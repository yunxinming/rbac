import { ElMessageBox, ElMessage } from 'element-plus'

// 计算表格操作列宽度
export function calculateActionWidth(buttons) {
  if (buttons.length === 0) return '0'
  const totalWidth = buttons.reduce((sum, text) => {
    const charWidth = 14  // 中文字符宽度
    const padding = 12    // 按钮左右padding
    const textWidth = text.length * charWidth + padding
    return sum + textWidth + 8  // 8px 按钮间距
  }, 0)
  return (totalWidth + 10) + 'px'  // 额外留10px边距
}

// 删除确认对话框
export async function confirmDelete(message, onConfirm) {
  await ElMessageBox.confirm(message, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await onConfirm()
  ElMessage.success('删除成功')
}
