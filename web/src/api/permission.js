import api from './index'

export const getPermissions = () => api.get('/permissions')
export const getPermissionTree = () => api.get('/permissions/tree')
export const getPermission = (id) => api.get(`/permissions/${id}`)
export const createPermission = (data) => api.post('/permissions', data)
export const updatePermission = (id, data) => api.put(`/permissions/${id}`, data)
export const deletePermission = (id) => api.delete(`/permissions/${id}`)
