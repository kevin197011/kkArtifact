// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Table, Tag, message } from 'antd'
import { auditApi, AuditLog } from '../api/audit'
import type { ColumnsType } from 'antd/es/table'

const AuditLogsPage: React.FC = () => {
  const [page, setPage] = useState(1)
  const pageSize = 50

  const { data, isLoading, error } = useQuery({
    queryKey: ['audit-logs', page],
    queryFn: () => auditApi.list(pageSize, (page - 1) * pageSize).then((res) => res.data),
  })

  if (error) {
    message.error('加载审计日志失败')
  }

  const columns: ColumnsType<AuditLog> = [
    {
      title: '操作',
      dataIndex: 'operation',
      key: 'operation',
      render: (op: string) => {
        const labels: Record<string, string> = {
          push: '推送',
          pull: '拉取',
          promote: '提升',
          token_create: '创建令牌',
          token_delete: '删除令牌',
        }
        return <Tag color="blue">{labels[op] || op}</Tag>
      },
    },
    {
      title: '版本哈希',
      dataIndex: 'version_hash',
      key: 'version_hash',
      render: (hash: string | undefined, record: AuditLog) => {
        if (!hash) return '-'
        // Display as project_app_version format if project_name and app_name are available
        if (record.project_name && record.app_name) {
          const displayText = `${record.project_name}_${record.app_name}_${hash}`
          return <Tag style={{ fontFamily: 'monospace' }}>{displayText}</Tag>
        }
        // Fallback to just the hash
        return <Tag style={{ fontFamily: 'monospace' }}>{hash}</Tag>
      },
    },
    {
      title: '代理ID',
      dataIndex: 'agent_id',
      key: 'agent_id',
      render: (id?: string) => id || '-',
    },
    {
      title: '元数据',
      dataIndex: 'metadata',
      key: 'metadata',
      render: (meta?: string) => {
        if (!meta) return '-'
        try {
          const parsed = JSON.parse(meta)
          return <pre style={{ margin: 0, fontSize: '12px' }}>{JSON.stringify(parsed, null, 2)}</pre>
        } catch {
          return meta
        }
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
      width: 200,
    },
  ]

  return (
    <div>
      <h2>审计日志</h2>
      <Table
        columns={columns}
        dataSource={data}
        loading={isLoading}
        rowKey="id"
        pagination={{
          current: page,
          pageSize,
          total: data?.length || 0,
          onChange: setPage,
        }}
      />
    </div>
  )
}

export default AuditLogsPage

