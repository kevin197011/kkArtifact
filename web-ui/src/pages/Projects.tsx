// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, message, Empty, Space, Popconfirm } from 'antd'
import { ReloadOutlined, DeleteOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { projectsApi, Project } from '../api/projects'
import { storageApi } from '../api/storage'
import type { ColumnsType } from 'antd/es/table'

const ProjectsPage: React.FC = () => {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const pageSize = 50
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['projects', page],
    queryFn: () => projectsApi.list(pageSize, (page - 1) * pageSize).then((res) => res.data),
    retry: 1,
  })

  const syncMutation = useMutation({
    mutationFn: () => storageApi.syncStorage(),
    onSuccess: (response) => {
      message.success(
        `同步完成：${response.data.projects} 个项目，${response.data.apps} 个应用，${response.data.versions} 个版本`
      )
      // Refresh projects list
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      refetch()
    },
    onError: (error: any) => {
      message.error(`同步失败：${error.response?.data?.error || error.message}`)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (projectName: string) => projectsApi.deleteProject(projectName),
    onSuccess: () => {
      message.success('项目删除成功')
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      refetch()
    },
    onError: (error: any) => {
      message.error(`删除项目失败：${error.response?.data?.error || error.message}`)
    },
  })

  const handleSyncStorage = () => {
    syncMutation.mutate()
  }

  if (error) {
    message.error('加载项目失败')
  }

  const columns: ColumnsType<Project> = [
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
            onClick={() => navigate(`/projects/${record.name}/apps`)}
            style={{
              padding: '0 8px',
              height: '32px',
              fontWeight: 500,
            }}
          >
            查看应用
          </Button>
          <Popconfirm
            title="确定要删除此项目吗？"
            description="删除项目将同时删除该项目下的所有应用和版本，此操作不可恢复！"
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
      <div style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '8px' }}>
          <div>
            <h2 style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
              项目
            </h2>
            <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
              管理所有项目及其应用和版本
            </div>
          </div>
          <Space>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => refetch()}
              loading={isLoading}
              style={{
                borderRadius: '6px',
                height: '36px',
                padding: '0 16px',
              }}
            >
              刷新列表
            </Button>
            <Button
              type="primary"
              icon={<ReloadOutlined />}
              onClick={handleSyncStorage}
              loading={syncMutation.isPending}
              style={{
                borderRadius: '6px',
                height: '36px',
                padding: '0 16px',
                fontWeight: 500,
              }}
            >
              同步存储
            </Button>
          </Space>
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
                return `共 ${total} 个项目`
              }
              return `第 ${range[0]}-${range[1]} 项，至少 ${total} 个项目`
            },
            style: {
              padding: '16px 24px',
            },
          }}
          locale={{
            emptyText: <Empty description="暂无项目" />,
          }}
          style={{
            borderRadius: '8px',
          }}
        />
      </div>
    </div>
  )
}

export default ProjectsPage
