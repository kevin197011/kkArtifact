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
      title: <span style={{ fontWeight: 600 }}>操作</span>,
      dataIndex: 'operation',
      key: 'operation',
      render: (op: string) => {
        const labels: Record<string, string> = {
          push: '推送',
          pull: '拉取',
          publish: '发布',
          unpublish: '取消发布',
          token_create: '创建令牌',
          token_delete: '删除令牌',
        }
        return <Tag color="blue">{labels[op] || op}</Tag>
      },
    },
    {
      title: <span style={{ fontWeight: 600 }}>版本哈希</span>,
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
      title: <span style={{ fontWeight: 600 }}>代理ID</span>,
      dataIndex: 'agent_id',
      key: 'agent_id',
      render: (id?: string) => (
        <span style={{ color: '#8c8c8c', fontSize: '14px' }}>{id || '-'}</span>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>元数据</span>,
      dataIndex: 'metadata',
      key: 'metadata',
      render: (meta?: string) => {
        if (!meta) return <span style={{ color: '#8c8c8c' }}>-</span>
        try {
          const parsed = JSON.parse(meta)
          return <pre style={{ margin: 0, fontSize: '12px', fontFamily: 'monospace', color: '#1a1a1a' }}>{JSON.stringify(parsed, null, 2)}</pre>
        } catch {
          return <span style={{ color: '#8c8c8c', fontSize: '14px' }}>{meta}</span>
        }
      },
    },
    {
      title: <span style={{ fontWeight: 600 }}>创建时间</span>,
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => (
        <span style={{ color: '#8c8c8c', fontSize: '14px' }}>
          {new Date(date).toLocaleString('zh-CN')}
        </span>
      ),
      width: 200,
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: '24px' }}>
        <h2 style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
          审计日志
        </h2>
        <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
          查看系统操作记录
        </div>
      </div>
      <div
        style={{
          background: 'var(--color-bg-primary)',
          borderRadius: 'var(--radius-md)',
          border: '1px solid var(--color-border-light)',
          overflow: 'hidden',
        }}
      >
        <Table
          columns={columns}
          dataSource={data}
          loading={isLoading}
          rowKey="id"
          pagination={{
            current: page,
            pageSize,
            total: data && data.length < pageSize 
              ? (page - 1) * pageSize + data.length 
              : data 
              ? page * pageSize + 1 
              : 0,
            onChange: setPage,
            showSizeChanger: false,
            showTotal: (total, range) => {
              if (data && data.length < pageSize) {
                return `共 ${total} 条审计日志`
              }
              return `第 ${range[0]}-${range[1]} 项，至少 ${total} 条审计日志`
            },
            style: {
              padding: '16px 24px',
            },
          }}
        />
      </div>
    </div>
  )
}

export default AuditLogsPage

