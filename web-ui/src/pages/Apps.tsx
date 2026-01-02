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
          <Button
            type="link"
            onClick={() => navigate(`/projects/${project}/apps/${record.name}/versions`)}
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
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={() => navigate('/projects')}>项目</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>{project}</Breadcrumb.Item>
      </Breadcrumb>
      <h2>应用 - {project}</h2>
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

export default AppsPage

