// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Breadcrumb, message, Space, Popconfirm } from 'antd'
import { DeleteOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import { projectsApi, App } from '../api/projects'
import type { ColumnsType } from 'antd/es/table'

const AppsPage: React.FC = () => {
  const { project } = useParams<{ project: string }>()
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const pageSize = 50
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['apps', project, page],
    queryFn: () =>
      projectsApi.getApps(project!, pageSize, (page - 1) * pageSize).then((res) => res.data),
    enabled: !!project,
  })

  const deleteMutation = useMutation({
    mutationFn: (appName: string) => projectsApi.deleteApp(project!, appName),
    onSuccess: () => {
      message.success('应用删除成功')
      queryClient.invalidateQueries({ queryKey: ['apps', project] })
      refetch()
    },
    onError: (error: any) => {
      message.error(`删除应用失败：${error.response?.data?.error || error.message}`)
    },
  })

  if (error) {
    message.error('加载应用失败')
  }

  const columns: ColumnsType<App> = [
    {
      title: <span style={{ fontWeight: 600 }}>名称</span>,
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <span style={{ fontWeight: 500, color: '#1a1a1a' }}>{text}</span>
      ),
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
    },
    {
      title: <span style={{ fontWeight: 600 }}>操作</span>,
      key: 'actions',
      width: 200,
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            onClick={() => navigate(`/projects/${project}/apps/${record.name}/versions`)}
            style={{
              padding: '0 8px',
              height: '32px',
              fontWeight: 500,
            }}
          >
            查看版本
          </Button>
          <Popconfirm
            title="确定要删除此应用吗？"
            description="删除应用将同时删除该应用下的所有版本，此操作不可恢复！"
            onConfirm={() => deleteMutation.mutate(record.name)}
            okText="确定"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button 
              type="link" 
              danger 
              icon={<DeleteOutlined />}
              style={{
                padding: '0 8px',
                height: '32px',
              }}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Breadcrumb 
        style={{ marginBottom: '24px', fontSize: '14px' }}
        items={[
          {
            title: <a onClick={() => navigate('/projects')} style={{ color: '#1890ff' }}>项目</a>,
          },
          {
            title: <span style={{ color: '#8c8c8c' }}>{project}</span>,
          },
        ]}
      />
      <div style={{ marginBottom: '24px' }}>
        <h2 style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
          应用 - {project}
        </h2>
        <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
          管理项目下的所有应用
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
                return `共 ${total} 个应用`
              }
              return `第 ${range[0]}-${range[1]} 项，至少 ${total} 个应用`
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

export default AppsPage

