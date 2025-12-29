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
      message.success('配置更新成功')
    },
    onError: () => {
      message.error('配置更新失败')
    },
  })

  const handleSubmit = (values: { version_retention_limit: number }) => {
    updateMutation.mutate(values)
  }

  return (
    <div>
      <h2>全局配置</h2>
      <Card style={{ maxWidth: 600 }}>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={data}
        >
          <Form.Item
            name="version_retention_limit"
            label="版本保留限制"
            rules={[
              { required: true, message: '请输入版本保留限制' },
              { type: 'number', min: 1, message: '必须至少为 1' },
            ]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={updateMutation.isPending}>
                保存
              </Button>
              <Button onClick={() => form.resetFields()}>重置</Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ConfigPage
