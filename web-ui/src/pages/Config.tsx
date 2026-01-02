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
    mutationFn: (data: { version_retention_limit?: number; audit_log_retention_days?: number }) => configApi.update(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['config'] })
      message.success('配置更新成功')
    },
    onError: () => {
      message.error('配置更新失败')
    },
  })

  const handleSubmit = (values: { version_retention_limit?: number; audit_log_retention_days?: number }) => {
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
            label="版本保留数量"
            tooltip="每个应用保留的最新版本数量。超过此数量的旧版本将在每天凌晨3点自动清理。"
            rules={[
              { required: true, message: '请输入版本保留数量' },
              { type: 'number', min: 1, message: '必须至少为 1' },
            ]}
          >
            <InputNumber 
              min={1} 
              style={{ width: '100%' }} 
              addonAfter="个版本"
              placeholder="例如：30"
            />
          </Form.Item>
          <Form.Item
            name="audit_log_retention_days"
            label="审计日志保留天数"
            tooltip="审计日志保留的天数。超过此天数的审计日志将在每天凌晨3点自动清理。"
            rules={[
              { required: true, message: '请输入审计日志保留天数' },
              { type: 'number', min: 1, message: '必须至少为 1' },
            ]}
          >
            <InputNumber 
              min={1} 
              style={{ width: '100%' }} 
              addonAfter="天"
              placeholder="例如：90"
            />
          </Form.Item>
          <div style={{ marginTop: 16, padding: 12, background: '#f0f2f5', borderRadius: 4 }}>
            <div style={{ marginBottom: 8, fontWeight: 500 }}>定时清理任务说明：</div>
            <ul style={{ margin: 0, paddingLeft: 20, color: '#666' }}>
              <li>清理任务每天凌晨 3:00 自动运行</li>
              <li><strong>版本清理</strong>：保留每个应用最新的 N 个版本，删除更旧的版本</li>
              <li><strong>审计日志清理</strong>：删除超过保留天数的审计日志记录</li>
              <li>清理范围：所有项目和应用</li>
              <li>清理内容：从存储和数据库同时删除旧数据</li>
            </ul>
          </div>
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
