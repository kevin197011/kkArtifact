// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Table, Button, Breadcrumb, message } from 'antd'
import { useParams, useNavigate } from 'react-router-dom'
import { projectsApi, App } from '../api/projects'
import type { ColumnsType } from 'antd/es/table'

const AppsPage: React.FC = () => {
  const { project } = useParams<{ project: string }>()
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const pageSize = 50

  const { data, isLoading, error } = useQuery({
    queryKey: ['apps', project, page],
    queryFn: () =>
      projectsApi.getApps(project!, pageSize, (page - 1) * pageSize).then((res) => res.data),
    enabled: !!project,
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
        <Button
          type="link"
          onClick={() => navigate(`/projects/${project}/apps/${record.name}/versions`)}
        >
          查看版本
        </Button>
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

