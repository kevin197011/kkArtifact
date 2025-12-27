// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Space, Popconfirm, message, Tag, Modal, Form, Input, Switch } from 'antd'
import { webhooksApi, Webhook, CreateWebhookRequest } from '../api/webhooks'
import type { ColumnsType } from 'antd/es/table'

const WebhooksPage: React.FC = () => {
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingWebhook, setEditingWebhook] = useState<Webhook | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

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
      message.success('Webhook created successfully')
    },
    onError: () => {
      message.error('Failed to create webhook')
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
      message.success('Webhook updated successfully')
    },
    onError: () => {
      message.error('Failed to update webhook')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => webhooksApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhooks'] })
      message.success('Webhook deleted successfully')
    },
    onError: () => {
      message.error('Failed to delete webhook')
    },
  })

  const handleCreate = () => {
    setEditingWebhook(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (webhook: Webhook) => {
    setEditingWebhook(webhook)
    form.setFieldsValue({
      name: webhook.name,
      url: webhook.url,
      event_types: webhook.event_types.join(','),
      enabled: webhook.enabled,
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
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
    {
      title: 'Event Types',
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
      title: 'Status',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'green' : 'red'}>{enabled ? 'Enabled' : 'Disabled'}</Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" size="small" onClick={() => handleEdit(record)}>
            Edit
          </Button>
          <Popconfirm
            title="Are you sure to delete this webhook?"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" size="small" danger>
              Delete
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
        title={editingWebhook ? 'Edit Webhook' : 'Create Webhook'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalVisible(false)
          setEditingWebhook(null)
          form.resetFields()
        }}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="url" label="URL" rules={[{ required: true, type: 'url' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="event_types" label="Event Types (comma separated)" rules={[{ required: true }]}>
            <Input placeholder="push,pull,promote" />
          </Form.Item>
          <Form.Item name="enabled" label="Enabled" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default WebhooksPage

