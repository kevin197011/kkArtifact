// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { Row, Col, Card, Statistic, Table, Typography, Spin } from 'antd'
import { ProjectOutlined, AppstoreOutlined, FileOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { projectsApi } from '../api/projects'
import { auditApi } from '../api/audit'
import type { AuditLog } from '../api/audit'
import type { ColumnType } from 'antd/es/table'

const { Title } = Typography

const Dashboard: React.FC = () => {
  // Fetch projects count
  const { data: projectsData, isLoading: projectsLoading } = useQuery({
    queryKey: ['projects', 'dashboard'],
    queryFn: () => projectsApi.list(1000, 0).then((res) => res.data),
  })

  // Fetch audit logs for recent activity
  const { data: auditLogsData, isLoading: auditLogsLoading } = useQuery({
    queryKey: ['audit-logs', 'dashboard'],
    queryFn: () => auditApi.list(10, 0).then((res) => res.data),
  })

  // Fetch apps for all projects to calculate total count
  const { data: allAppsData } = useQuery({
    queryKey: ['all-apps', 'dashboard', projectsData?.map((p) => p.name)],
    queryFn: async () => {
      if (!projectsData || projectsData.length === 0) return []
      const appsPromises = projectsData.map((project) =>
        projectsApi.getApps(project.name, 1000, 0).then((res) => res.data)
      )
      const appsArrays = await Promise.all(appsPromises)
      return appsArrays.flat()
    },
    enabled: !!projectsData && projectsData.length > 0,
  })

  // Calculate statistics
  const projectsCount = projectsData?.length || 0
  const appsCount = allAppsData?.length || 0
  const isLoading = projectsLoading || auditLogsLoading

  // Audit logs columns
  const auditColumns: ColumnType<AuditLog>[] = [
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => {
        const date = new Date(text)
        return <span style={{ color: 'var(--color-text-secondary)', fontSize: '13px' }}>{date.toLocaleString('zh-CN')}</span>
      },
    },
    {
      title: '操作',
      dataIndex: 'operation',
      key: 'operation',
      width: 120,
      render: (text: string) => {
        const colorMap: Record<string, string> = {
          push: 'var(--color-primary)',
          pull: 'var(--color-success)',
          publish: 'var(--color-warning)',
          unpublish: 'var(--color-warning)',
          token_create: '#722ed1',
          token_delete: 'var(--color-error)',
        }
        const labels: Record<string, string> = {
          push: '推送',
          pull: '拉取',
          publish: '发布',
          unpublish: '取消发布',
          token_create: '创建令牌',
          token_delete: '删除令牌',
        }
        return <span style={{ color: colorMap[text] || 'var(--color-text-primary)', fontWeight: 500, fontSize: '13px' }}>{labels[text] || text}</span>
      },
    },
    {
      title: '代理ID',
      dataIndex: 'agent_id',
      key: 'agent_id',
      width: 200,
      ellipsis: true,
      render: (text: string) => <span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontFamily: 'monospace' }}>{text || '-'}</span>,
    },
    {
      title: '版本',
      dataIndex: 'version_hash',
      key: 'version_hash',
      width: 200,
      ellipsis: true,
      render: (hash: string, record: AuditLog) => {
        const displayText = record.project_name && record.app_name
          ? `${record.project_name}_${record.app_name}_${hash}`
          : hash
        return <span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontFamily: 'monospace' }}>{displayText}</span>
      },
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: '32px' }}>
        <Title level={2} style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
          仪表盘
        </Title>
        <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
          系统概览和最近活动
        </div>
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: '32px' }}>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            loading={isLoading}
            hoverable
            onClick={() => window.location.href = '/projects'}
            style={{ cursor: 'pointer' }}
            bodyStyle={{ padding: '20px' }}
          >
            <Statistic
              title={<span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontWeight: 500 }}>项目数</span>}
              value={projectsCount}
              prefix={<ProjectOutlined style={{ color: 'var(--color-primary)', fontSize: '18px', marginRight: '8px' }} />}
              valueStyle={{ color: 'var(--color-text-primary)', fontSize: '24px', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            loading={isLoading}
            hoverable
            bodyStyle={{ padding: '20px' }}
          >
            <Statistic
              title={<span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontWeight: 500 }}>应用总数</span>}
              value={appsCount}
              prefix={<AppstoreOutlined style={{ color: 'var(--color-success)', fontSize: '18px', marginRight: '8px' }} />}
              valueStyle={{ color: 'var(--color-text-primary)', fontSize: '24px', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card
            hoverable
            onClick={() => window.location.href = '/audit-logs'}
            style={{ cursor: 'pointer' }}
            bodyStyle={{ padding: '20px' }}
          >
            <Statistic
              title={<span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontWeight: 500 }}>最近活动</span>}
              value={auditLogsData?.length || 0}
              prefix={<ClockCircleOutlined style={{ color: 'var(--color-warning)', fontSize: '18px', marginRight: '8px' }} />}
              valueStyle={{ color: 'var(--color-text-primary)', fontSize: '24px', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card
            bodyStyle={{ padding: '20px' }}
          >
            <Statistic
              title={<span style={{ color: 'var(--color-text-secondary)', fontSize: '13px', fontWeight: 500 }}>系统状态</span>}
              value="运行中"
              prefix={<FileOutlined style={{ color: 'var(--color-success)', fontSize: '18px', marginRight: '8px' }} />}
              valueStyle={{ color: 'var(--color-success)', fontSize: '16px', fontWeight: 600 }}
            />
          </Card>
        </Col>
      </Row>

      <Card 
        title={
          <span style={{ fontSize: '16px', fontWeight: 600, color: 'var(--color-text-primary)' }}>
            最近活动
          </span>
        }
        bodyStyle={{ padding: '20px' }}
      >
        {auditLogsLoading ? (
          <div style={{ textAlign: 'center', padding: '60px' }}>
            <Spin size="large" />
          </div>
        ) : (
          <Table
            columns={auditColumns}
            dataSource={auditLogsData || []}
            rowKey="id"
            pagination={false}
            size="middle"
          />
        )}
      </Card>
    </div>
  )
}

export default Dashboard

