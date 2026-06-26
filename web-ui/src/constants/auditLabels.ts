// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

/** 审计日志操作类型中文标签 */
export const AUDIT_OPERATION_LABELS: Record<string, string> = {
  push: '推送',
  pull: '拉取',
  publish: '发布',
  unpublish: '取消发布',
  token_create: '创建令牌',
  token_delete: '删除令牌',
  version_delete: '删除',
}

export const AUDIT_OPERATION_OPTIONS = Object.entries(AUDIT_OPERATION_LABELS).map(
  ([value, label]) => ({ value, label })
)

export function formatAuditOperation(op: string): string {
  return AUDIT_OPERATION_LABELS[op] || op
}
