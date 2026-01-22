// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { Row, Col, Card, Table, Typography, Spin } from 'antd'
import { ProjectOutlined, AppstoreOutlined, TagOutlined, ClockCircleOutlined, HistoryOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { projectsApi } from '../api/projects'
import { auditApi } from '../api/audit'
import type { AuditLog } from '../api/audit'
import type { ColumnType } from 'antd/es/table'
import styles from './Dashboard.module.css'

const { Title } = Typography

const Dashboard: React.FC = () => {
  const navigate = useNavigate()
  
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

  // Fetch versions for all apps to calculate total count
  const { data: allVersionsData } = useQuery({
    queryKey: ['all-versions', 'dashboard', allAppsData?.map((a) => `${a.project_id}/${a.name}`)],
    queryFn: async () => {
      if (!allAppsData || !projectsData || allAppsData.length === 0) return []
      // Create a map from project_id to project name
      const projectMap = new Map(projectsData.map((p) => [p.id, p.name]))
      const versionsPromises = allAppsData.map((app) => {
        const projectName = projectMap.get(app.project_id)
        if (!projectName) {
          return Promise.resolve([] as any[])
        }
        return projectsApi.getVersions(projectName, app.name, 1000, 0).then((res) => res.data)
      })
      const versionsArrays = await Promise.all(versionsPromises)
      return versionsArrays.flat()
    },
    enabled: !!allAppsData && !!projectsData && allAppsData.length > 0,
  })

  // Calculate statistics
  const projectsCount = projectsData?.length || 0
  const appsCount = allAppsData?.length || 0
  const versionsCount = allVersionsData?.length || 0
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
        const labels: Record<string, string> = {
          push: '推送',
          pull: '拉取',
          publish: '发布',
          unpublish: '取消发布',
          token_create: '创建令牌',
          token_delete: '删除令牌',
        }
        return (
          <span className={`${styles.operationBadge} ${styles[text] || ''}`}>
            {labels[text] || text}
          </span>
        )
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
    <div className={styles.dashboardContainer}>
      <div className={styles.header}>
        <Title level={2} className={styles.title}>
          仪表盘
        </Title>
        <div className={styles.subtitle}>
          系统概览和最近活动
        </div>
      </div>

      <Row gutter={[20, 20]} style={{ marginBottom: '32px' }}>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            loading={isLoading}
            className={`${styles.statCard} ${styles.projectCard}`}
            onClick={() => navigate('/projects')}
            bodyStyle={{ padding: 0 }}
          >
            <div className={styles.statCardBody}>
              <div className={`${styles.statIconWrapper} ${styles.projectIcon}`}>
                <ProjectOutlined className={styles.statIcon} style={{ color: 'var(--color-primary)' }} />
              </div>
              <div className={styles.statContent}>
                <span className={styles.statTitle}>项目数</span>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: '8px' }}>
                  <span className={styles.statPrefix}>
                    <ProjectOutlined style={{ color: 'var(--color-primary)', fontSize: '20px' }} />
                  </span>
                  <span className={styles.statValue}>{projectsCount}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            loading={isLoading}
            className={`${styles.statCard} ${styles.appCard}`}
            bodyStyle={{ padding: 0 }}
          >
            <div className={styles.statCardBody}>
              <div className={`${styles.statIconWrapper} ${styles.appIcon}`}>
                <AppstoreOutlined className={styles.statIcon} style={{ color: 'var(--color-success)' }} />
              </div>
              <div className={styles.statContent}>
                <span className={styles.statTitle}>应用总数</span>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: '8px' }}>
                  <span className={styles.statPrefix}>
                    <AppstoreOutlined style={{ color: 'var(--color-success)', fontSize: '20px' }} />
                  </span>
                  <span className={styles.statValue}>{appsCount}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card
            loading={isLoading}
            className={`${styles.statCard} ${styles.versionCard}`}
            bodyStyle={{ padding: 0 }}
          >
            <div className={styles.statCardBody}>
              <div className={`${styles.statIconWrapper} ${styles.versionIcon}`}>
                <TagOutlined className={styles.statIcon} style={{ color: 'var(--color-primary)' }} />
              </div>
              <div className={styles.statContent}>
                <span className={styles.statTitle}>版本总数</span>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: '8px' }}>
                  <span className={styles.statPrefix}>
                    <TagOutlined style={{ color: 'var(--color-primary)', fontSize: '20px' }} />
                  </span>
                  <span className={styles.statValue}>{versionsCount}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card
            className={`${styles.statCard} ${styles.activityCard}`}
            onClick={() => navigate('/audit-logs')}
            bodyStyle={{ padding: 0 }}
          >
            <div className={styles.statCardBody}>
              <div className={`${styles.statIconWrapper} ${styles.activityIcon}`}>
                <ClockCircleOutlined className={styles.statIcon} style={{ color: 'var(--color-warning)' }} />
              </div>
              <div className={styles.statContent}>
                <span className={styles.statTitle}>最近活动</span>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: '8px' }}>
                  <span className={styles.statPrefix}>
                    <ClockCircleOutlined style={{ color: 'var(--color-warning)', fontSize: '20px' }} />
                  </span>
                  <span className={styles.statValue}>{auditLogsData?.length || 0}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Card 
        className={styles.activityCard}
        bodyStyle={{ padding: 0 }}
      >
        <div className={styles.activityCardHeader}>
          <h3 className={styles.activityCardTitle}>
            <HistoryOutlined style={{ fontSize: '18px', color: 'var(--color-primary)' }} />
            最近活动
          </h3>
        </div>
        <div className={styles.activityCardBody}>
          {auditLogsLoading ? (
            <div className={styles.emptyState}>
              <Spin size="large" />
            </div>
          ) : auditLogsData && auditLogsData.length > 0 ? (
            <Table
              columns={auditColumns}
              dataSource={auditLogsData}
              rowKey="id"
              pagination={false}
              size="middle"
              style={{ background: 'transparent' }}
            />
          ) : (
            <div className={styles.emptyState}>
              <div className={styles.emptyStateIcon}>
                <HistoryOutlined />
              </div>
              <div className={styles.emptyStateText}>暂无活动记录</div>
            </div>
          )}
        </div>
      </Card>
    </div>
  )
}

export default Dashboard

