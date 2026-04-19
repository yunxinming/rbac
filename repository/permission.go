package repository

import (
	"hello/model"
)

type PermissionRepository struct{}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{}
}

func (r *PermissionRepository) Create(permission *model.Permission) error {
	return DB.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *model.Permission) error {
	return DB.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	return DB.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) GetAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetByCodes(codes []string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Where("code IN ?", codes).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := DB.Model(&model.Permission{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// GetAllWithChildren 获取所有权限（含子权限）- 优化：一次查询+内存构建
func (r *PermissionRepository) GetAllWithChildren() ([]model.Permission, error) {
	var permissions []model.Permission
	// 一次性查询所有权限
	err := DB.Order("sort asc").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	// 在内存中构建树
	return r.buildTree(permissions, 0), nil
}

// buildTree 在内存中构建权限树
func (r *PermissionRepository) buildTree(permissions []model.Permission, parentID uint) []model.Permission {
	var result []model.Permission
	for i := range permissions {
		if permissions[i].ParentID == parentID {
			permissions[i].Children = r.buildTree(permissions, permissions[i].ID)
			result = append(result, permissions[i])
		}
	}
	return result
}

// GetUserMenuIDs 获取用户的菜单ID列表（只有read权限才显示菜单）
func (r *PermissionRepository) GetUserMenuIDs(userID uint) ([]uint, error) {
	// 1. 获取用户直接拥有的菜单/目录权限
	var directMenuIDs []uint
	err := DB.Table("permissions").
		Select("DISTINCT permissions.id").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.type IN ?", userID, []model.PermissionType{model.PermissionTypeDirectory, model.PermissionTypeMenu}).
		Pluck("permissions.id", &directMenuIDs).Error
	if err != nil {
		return nil, err
	}

	// 2. 获取用户拥有的 read 权限的父级菜单ID（只有read权限才显示菜单）
	var readParentIDs []uint
	err = DB.Table("permissions AS op").
		Select("DISTINCT op.parent_id").
		Joins("JOIN role_permissions ON op.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND op.type = ? AND op.code LIKE ?", userID, model.PermissionTypeOperation, "%:read").
		Pluck("op.parent_id", &readParentIDs).Error
	if err != nil {
		return nil, err
	}

	// 合并菜单ID
	menuIDs := append(directMenuIDs, readParentIDs...)

	// 3. 预加载所有权限的 parent_id 关系，避免 N+1 查询
	var allPerms []struct {
		ID       uint
		ParentID uint
	}
	err = DB.Table("permissions").Select("id, parent_id").Find(&allPerms).Error
	if err != nil {
		return nil, err
	}

	// 构建 id -> parent_id 映射
	permMap := make(map[uint]uint)
	for _, p := range allPerms {
		permMap[p.ID] = p.ParentID
	}

	// 4. 使用内存查找递归获取所有父级目录
	allMenuIDs := make(map[uint]bool)
	for _, id := range menuIDs {
		allMenuIDs[id] = true
		r.loadParentMenuIDsInMemory(id, permMap, allMenuIDs)
	}

	// 转换为切片
	result := make([]uint, 0, len(allMenuIDs))
	for id := range allMenuIDs {
		result = append(result, id)
	}
	return result, nil
}

// loadParentMenuIDsInMemory 使用内存映射递归加载父级菜单ID
func (r *PermissionRepository) loadParentMenuIDsInMemory(menuID uint, permMap map[uint]uint, menuIDs map[uint]bool) {
	parentID, exists := permMap[menuID]
	if !exists || parentID == 0 {
		return
	}
	if !menuIDs[parentID] {
		menuIDs[parentID] = true
		r.loadParentMenuIDsInMemory(parentID, permMap, menuIDs)
	}
}

// GetUserMenus 获取用户有权限的菜单树
func (r *PermissionRepository) GetUserMenus(userID uint) ([]model.Permission, error) {
	// 获取用户所有菜单ID
	menuIDs, err := r.GetUserMenuIDs(userID)
	if err != nil {
		return nil, err
	}

	if len(menuIDs) == 0 {
		return []model.Permission{}, nil
	}

	// 获取所有目录和菜单
	var allMenus []model.Permission
	err = DB.Where("type IN ?", []model.PermissionType{model.PermissionTypeDirectory, model.PermissionTypeMenu}).
		Order("sort asc").Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// 构建菜单树并过滤权限
	return r.buildUserMenuTree(allMenus, menuIDs, 0), nil
}

// buildUserMenuTree 构建用户菜单树
func (r *PermissionRepository) buildUserMenuTree(allMenus []model.Permission, menuIDs []uint, parentID uint) []model.Permission {
	var result []model.Permission
	menuIDSet := make(map[uint]bool)
	for _, id := range menuIDs {
		menuIDSet[id] = true
	}

	for i := range allMenus {
		if allMenus[i].ParentID == parentID {
			// 检查是否有权限（或子菜单有权限）
			if r.hasMenuAccess(allMenus[i], allMenus, menuIDSet) {
				allMenus[i].Children = r.buildUserMenuTree(allMenus, menuIDs, allMenus[i].ID)
				result = append(result, allMenus[i])
			}
		}
	}
	return result
}

// hasMenuAccess 检查用户是否有菜单访问权限
func (r *PermissionRepository) hasMenuAccess(menu model.Permission, allMenus []model.Permission, menuIDSet map[uint]bool) bool {
	// 直接有权限
	if menuIDSet[menu.ID] {
		return true
	}
	// 检查子菜单是否有权限
	for _, m := range allMenus {
		if m.ParentID == menu.ID {
			if r.hasMenuAccess(m, allMenus, menuIDSet) {
				return true
			}
		}
	}
	return false
}

// GetByType 按类型查询
func (r *PermissionRepository) GetByType(permType model.PermissionType) ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Where("type = ?", permType).Order("sort asc").Find(&permissions).Error
	return permissions, err
}
