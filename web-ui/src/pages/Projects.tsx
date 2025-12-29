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
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" onClick={() => navigate(`/projects/${record.name}/apps`)}>
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
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>项目</h2>
        <Space>
          <Button
            icon={<ReloadOutlined />}
            onClick={() => refetch()}
            loading={isLoading}
          >
            刷新列表
          </Button>
          <Button
            type="primary"
            icon={<ReloadOutlined />}
            onClick={handleSyncStorage}
            loading={syncMutation.isPending}
          >
            同步存储
          </Button>
        </Space>
      </div>
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
        locale={{
          emptyText: <Empty description="暂无项目" />,
        }}
      />
    </div>
  )
}

export default ProjectsPage
