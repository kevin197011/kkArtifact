// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Space, Popconfirm, message, Tag, Modal, Form, Input, Switch, Select } from 'antd'
import { webhooksApi, Webhook, CreateWebhookRequest } from '../api/webhooks'
import { projectsApi } from '../api/projects'
import type { ColumnsType } from 'antd/es/table'

const WebhooksPage: React.FC = () => {
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingWebhook, setEditingWebhook] = useState<Webhook | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [selectedProject, setSelectedProject] = useState<string | undefined>(undefined)

  // Fetch projects
  const { data: projects } = useQuery({
    queryKey: ['projects'],
    queryFn: () => projectsApi.list(1000, 0).then((res) => res.data),
  })

  // Fetch apps for selected project
  const { data: apps } = useQuery({
    queryKey: ['apps', selectedProject],
    queryFn: () => projectsApi.getApps(selectedProject!, 1000, 0).then((res) => res.data),
    enabled: !!selectedProject,
  })

  const { data, isLoading } = useQuery({
    queryKey: ['webhooks'],
    queryFn: () => webhooksApi.list().then((res) => res.data),
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateWebhookRequest) => webhooksApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhooks'] })
      setIsModalVisible(false)
      form.resetFields()
      message.success('Webhook 创建成功')
    },
    onError: () => {
      message.error('Webhook 创建失败')
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<CreateWebhookRequest> }) =>
      webhooksApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhooks'] })
      setIsModalVisible(false)
      setEditingWebhook(null)
      form.resetFields()
      message.success('Webhook 更新成功')
    },
    onError: () => {
      message.error('Webhook 更新失败')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => webhooksApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhooks'] })
      message.success('Webhook 删除成功')
    },
    onError: () => {
      message.error('Webhook 删除失败')
    },
  })

  const handleProjectChange = (projectId: number | undefined) => {
    if (projectId && projects) {
      const project = projects.find((p) => p.id === projectId)
      setSelectedProject(project?.name)
    } else {
      setSelectedProject(undefined)
    }
    form.setFieldsValue({ app_id: undefined })
  }

  const handleCreate = () => {
    setEditingWebhook(null)
    setSelectedProject(undefined)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (webhook: Webhook) => {
    setEditingWebhook(webhook)
    if (webhook.project_id && projects) {
      const project = projects.find((p) => p.id === webhook.project_id)
      setSelectedProject(project?.name)
    } else {
      setSelectedProject(undefined)
    }
    form.setFieldsValue({
      name: webhook.name,
      url: webhook.url,
      event_types: webhook.event_types.join(','),
      enabled: webhook.enabled,
      project_id: webhook.project_id,
      app_id: webhook.app_id,
    })
    setIsModalVisible(true)
  }

  const handleSubmit = () => {
    form.validateFields().then((values) => {
      const data: CreateWebhookRequest = {
        name: values.name,
        url: values.url,
        event_types: values.event_types.split(',').map((s: string) => s.trim()),
        enabled: values.enabled ?? true,
        project_id: values.project_id,
        app_id: values.app_id,
      }

      if (editingWebhook) {
        updateMutation.mutate({ id: editingWebhook.id, data })
      } else {
        createMutation.mutate(data)
      }
    })
  }

  const columns: ColumnsType<Webhook> = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
    {
      title: '项目 / 应用',
      key: 'project_app',
      render: (_: any, record: Webhook) => {
        if (record.project_name) {
          if (record.app_name) {
            return <Tag>{record.project_name} / {record.app_name}</Tag>
          }
          return <Tag>{record.project_name}</Tag>
        }
        return <Tag color="default">全局</Tag>
      },
    },
    {
      title: '事件类型',
      dataIndex: 'event_types',
      key: 'event_types',
      render: (types: string[]) => (
        <Space>
          {types.map((type) => (
            <Tag key={type}>{type}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'green' : 'red'}>{enabled ? '启用' : '禁用'}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" size="small" onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此 Webhook 吗？"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" size="small" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>Webhooks</h2>
        <Button type="primary" onClick={handleCreate}>
          Create Webhook
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={data}
        loading={isLoading}
        rowKey="id"
      />
      <Modal
        title={editingWebhook ? '编辑 Webhook' : '创建 Webhook'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalVisible(false)
          setEditingWebhook(null)
          setSelectedProject(undefined)
          form.resetFields()
        }}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="url" label="URL" rules={[{ required: true, type: 'url' }]}>
            <Input />
          </Form.Item>
          <Form.Item
            name="project_id"
            label="项目（可选）"
            tooltip="选择一个项目以限制 Webhook 到特定项目。留空则为全局 Webhook。"
          >
            <Select
              placeholder="所有项目（全局）"
              allowClear
              onChange={handleProjectChange}
              showSearch
              filterOption={(input, option) =>
                String(option?.children || '').toLowerCase().includes(input.toLowerCase())
              }
            >
              {projects?.map((project) => (
                <Select.Option key={project.id} value={project.id}>
                  {project.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="app_id"
            label="应用（可选）"
            tooltip="选择一个应用以限制 Webhook 到特定应用。需要先选择一个项目。"
          >
            <Select
              placeholder="所有应用"
              allowClear
              disabled={!selectedProject}
              showSearch
              filterOption={(input, option) =>
                String(option?.children || '').toLowerCase().includes(input.toLowerCase())
              }
            >
              {apps?.map((app) => (
                <Select.Option key={app.id} value={app.id}>
                  {app.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="event_types" label="事件类型（逗号分隔）" rules={[{ required: true }]}>
            <Input placeholder="push,pull,publish" />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default WebhooksPage

