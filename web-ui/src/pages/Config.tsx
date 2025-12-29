// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Form, InputNumber, Button, Card, message, Space } from 'antd'
import { configApi } from '../api/config'

const ConfigPage: React.FC = () => {
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['config'],
    queryFn: () => configApi.get().then((res) => res.data),
  })

  useEffect(() => {
    if (data) {
      form.setFieldsValue(data)
    }
  }, [data, form])

  const updateMutation = useMutation({
    mutationFn: (data: { version_retention_limit: number }) => configApi.update(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['config'] })
      message.success('Configuration updated successfully')
    },
    onError: () => {
      message.error('Failed to update configuration')
    },
  })

  const handleSubmit = (values: { version_retention_limit: number }) => {
    updateMutation.mutate(values)
  }

  return (
    <div>
      <h2>Global Configuration</h2>
      <Card style={{ maxWidth: 600 }}>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={data}
        >
          <Form.Item
            name="version_retention_limit"
            label="Version Retention Limit"
            rules={[
              { required: true, message: 'Please input version retention limit' },
              { type: 'number', min: 1, message: 'Must be at least 1' },
            ]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={updateMutation.isPending}>
                Save
              </Button>
              <Button onClick={() => form.resetFields()}>Reset</Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ConfigPage
