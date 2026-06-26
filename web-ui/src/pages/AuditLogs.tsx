// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Table, Tag, message, Select, Space, Button } from 'antd'
import { auditApi, AuditLog } from '../api/audit'
import { projectsApi } from '../api/projects'
import { AUDIT_OPERATION_OPTIONS, formatAuditOperation } from '../constants/auditLabels'
import type { ColumnsType } from 'antd/es/table'

const AuditLogsPage: React.FC = () => {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(50)
  const [projectId, setProjectId] = useState<number | undefined>()
  const [appId, setAppId] = useState<number | undefined>()
  const [operation, setOperation] = useState<string | undefined>()

  const { data: projects } = useQuery({
    queryKey: ['projects', 'audit-filter'],
    queryFn: () => projectsApi.list(1000, 0).then((res) => res.data),
  })

  const selectedProjectName = projects?.find((p) => p.id === projectId)?.name

  const { data: apps } = useQuery({
    queryKey: ['apps', 'audit-filter', selectedProjectName],
    queryFn: () => projectsApi.getApps(selectedProjectName!, 1000, 0).then((res) => res.data),
    enabled: !!selectedProjectName,
  })

  const filters = {
    project_id: projectId,
    app_id: appId,
    operation,
  }

  const { data, isLoading, error } = useQuery({
    queryKey: ['audit-logs', page, pageSize, filters],
    queryFn: async () => {
      const response = await auditApi.list(pageSize, (page - 1) * pageSize, filters)
      return response.data
    },
  })

  useEffect(() => {
    if (error) {
      message.error('加载审计日志失败')
    }
  }, [error])

  const handleResetFilters = () => {
    setProjectId(undefined)
    setAppId(undefined)
    setOperation(undefined)
    setPage(1)
  }

  const columns: ColumnsType<AuditLog> = [
    {
      title: <span style={{ fontWeight: 600 }}>操作</span>,
      dataIndex: 'operation',
      key: 'operation',
      width: 120,
      render: (op: string) => <Tag color="blue">{formatAuditOperation(op)}</Tag>,
    },
    {
      title: <span style={{ fontWeight: 600 }}>项目 / 应用</span>,
      key: 'scope',
      width: 220,
      render: (_: unknown, record: AuditLog) => {
        if (record.project_name || record.app_name) {
          return (
            <span style={{ fontSize: '13px', color: 'var(--color-text-secondary)' }}>
              {[record.project_name, record.app_name].filter(Boolean).join(' / ') || '-'}
            </span>
          )
        }
        return '-'
      },
    },
    {
      title: <span style={{ fontWeight: 600 }}>版本哈希</span>,
      dataIndex: 'version_hash',
      key: 'version_hash',
      render: (hash: string | undefined, record: AuditLog) => {
        if (!hash) return '-'
        if (record.project_name && record.app_name) {
          const displayText = `${record.project_name}_${record.app_name}_${hash}`
          return <Tag style={{ fontFamily: 'monospace' }}>{displayText}</Tag>
        }
        return <Tag style={{ fontFamily: 'monospace' }}>{hash}</Tag>
      },
    },
    {
      title: <span style={{ fontWeight: 600 }}>代理ID</span>,
      dataIndex: 'agent_id',
      key: 'agent_id',
      render: (id?: string) => (
        <span style={{ color: 'var(--color-text-secondary)', fontSize: '14px' }}>{id || '-'}</span>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>创建时间</span>,
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => (
        <span style={{ color: 'var(--color-text-secondary)', fontSize: '14px' }}>
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
          查看系统操作记录，支持按项目、应用、操作类型筛选
        </div>
      </div>

      <div
        style={{
          marginBottom: '16px',
          padding: '16px',
          background: 'var(--color-bg-primary)',
          borderRadius: 'var(--radius-md)',
          border: '1px solid var(--color-border-light)',
        }}
      >
        <Space wrap>
          <Select
            allowClear
            placeholder="筛选项目"
            style={{ minWidth: 160 }}
            value={projectId}
            onChange={(value) => {
              setProjectId(value)
              setAppId(undefined)
              setPage(1)
            }}
            options={projects?.map((p) => ({ value: p.id, label: p.name }))}
          />
          <Select
            allowClear
            placeholder="筛选应用"
            style={{ minWidth: 200 }}
            value={appId}
            disabled={!projectId}
            onChange={(value) => {
              setAppId(value)
              setPage(1)
            }}
            options={apps?.map((a) => ({ value: a.id, label: a.name }))}
          />
          <Select
            allowClear
            placeholder="筛选操作"
            style={{ minWidth: 140 }}
            value={operation}
            onChange={(value) => {
              setOperation(value)
              setPage(1)
            }}
            options={AUDIT_OPERATION_OPTIONS}
          />
          <Button onClick={handleResetFilters}>重置</Button>
        </Space>
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
          dataSource={data?.data || []}
          loading={isLoading}
          rowKey="id"
          pagination={{
            current: page,
            pageSize,
            total: data?.total || 0,
            onChange: (newPage, newPageSize) => {
              setPage(newPage)
              if (newPageSize !== pageSize) {
                setPageSize(newPageSize)
                setPage(1)
              }
            },
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100'],
            showTotal: (total, range) => `共 ${total} 条审计日志，显示第 ${range[0]}-${range[1]} 项`,
            style: { padding: '16px 24px' },
          }}
        />
      </div>
    </div>
  )
}

export default AuditLogsPage
