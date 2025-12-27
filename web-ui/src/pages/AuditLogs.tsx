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
    message.error('Failed to load audit logs')
  }

  const columns: ColumnsType<AuditLog> = [
    {
      title: 'Operation',
      dataIndex: 'operation',
      key: 'operation',
      render: (op: string) => <Tag color="blue">{op}</Tag>,
    },
    {
      title: 'Version Hash',
      dataIndex: 'version_hash',
      key: 'version_hash',
      render: (hash?: string) =>
        hash ? <Tag style={{ fontFamily: 'monospace' }}>{hash.substring(0, 12)}...</Tag> : '-',
    },
    {
      title: 'Agent ID',
      dataIndex: 'agent_id',
      key: 'agent_id',
      render: (id?: string) => id || '-',
    },
    {
      title: 'Metadata',
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
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
      width: 200,
    },
  ]

  return (
    <div>
      <h2>Audit Logs</h2>
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

