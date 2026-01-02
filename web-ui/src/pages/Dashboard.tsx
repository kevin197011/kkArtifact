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
        return date.toLocaleString('zh-CN')
      },
    },
    {
      title: '操作',
      dataIndex: 'operation',
      key: 'operation',
      width: 120,
      render: (text: string) => {
        const colors: Record<string, string> = {
          push: 'blue',
          pull: 'green',
          promote: 'orange',
          token_create: 'purple',
          token_delete: 'red',
        }
        const labels: Record<string, string> = {
          push: '推送',
          pull: '拉取',
          promote: '提升',
          token_create: '创建令牌',
          token_delete: '删除令牌',
        }
        return <span style={{ color: colors[text] || 'default' }}>{labels[text] || text}</span>
      },
    },
    {
      title: '代理ID',
      dataIndex: 'agent_id',
      key: 'agent_id',
      width: 200,
      ellipsis: true,
    },
    {
      title: '版本',
      dataIndex: 'version_hash',
      key: 'version_hash',
      width: 200,
      ellipsis: true,
      render: (hash: string, record: AuditLog) => {
        if (record.project_name && record.app_name) {
          return `${record.project_name}_${record.app_name}_${hash}`
        }
        return hash
      },
    },
  ]

  return (
    <div>
      <Title level={2}>仪表盘</Title>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="项目数"
              value={projectsCount}
              prefix={<ProjectOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="应用总数"
              value={appsCount}
              prefix={<AppstoreOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="最近活动"
              value={auditLogsData?.length || 0}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#fa8c16' }}
              suffix="最近10条"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="系统状态"
              value="运行中"
              prefix={<FileOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title="最近活动" style={{ marginBottom: 24 }}>
        {auditLogsLoading ? (
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <Spin size="large" />
          </div>
        ) : (
          <Table
            columns={auditColumns}
            dataSource={auditLogsData || []}
            rowKey="id"
            pagination={false}
            size="small"
          />
        )}
      </Card>
    </div>
  )
}

export default Dashboard

