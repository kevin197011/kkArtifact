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
      <div style={{ marginBottom: '24px' }}>
        <h2 style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
          全局配置
        </h2>
        <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
          管理系统全局设置
        </div>
      </div>
      <Card 
        style={{ 
          maxWidth: 700,
        }}
        bodyStyle={{ padding: '32px' }}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={data}
        >
          <Form.Item
            name="version_retention_limit"
            label={<span style={{ fontWeight: 600, fontSize: '15px', color: 'var(--color-text-primary)' }}>版本保留数量</span>}
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
              size="large"
            />
          </Form.Item>
          <Form.Item
            name="audit_log_retention_days"
            label={<span style={{ fontWeight: 600, fontSize: '15px', color: 'var(--color-text-primary)' }}>审计日志保留天数</span>}
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
              size="large"
            />
          </Form.Item>
          <div style={{ 
            marginTop: '24px', 
            padding: '20px', 
            background: 'var(--color-bg-secondary)', 
            borderRadius: 'var(--radius-md)',
            border: '1px solid var(--color-border-light)',
          }}>
            <div style={{ marginBottom: '12px', fontWeight: 600, fontSize: '14px', color: 'var(--color-text-primary)' }}>
              定时清理任务说明
            </div>
            <ul style={{ margin: 0, paddingLeft: '24px', color: 'var(--color-text-secondary)', lineHeight: '1.8', fontSize: '13px' }}>
              <li>清理任务每天凌晨 3:00 自动运行</li>
              <li><strong style={{ color: 'var(--color-text-primary)' }}>版本清理</strong>：保留每个应用最新的 N 个版本，删除更旧的版本</li>
              <li><strong style={{ color: 'var(--color-text-primary)' }}>审计日志清理</strong>：删除超过保留天数的审计日志记录</li>
              <li>清理范围：所有项目和应用</li>
              <li>清理内容：从存储和数据库同时删除旧数据</li>
            </ul>
          </div>
          <Form.Item style={{ marginTop: '32px', marginBottom: 0 }}>
            <Space>
              <Button 
                type="primary" 
                htmlType="submit" 
                loading={updateMutation.isPending}
                style={{
                  borderRadius: '6px',
                  height: '40px',
                  padding: '0 24px',
                  fontWeight: 500,
                }}
              >
                保存
              </Button>
              <Button 
                onClick={() => form.resetFields()}
                style={{
                  borderRadius: '6px',
                  height: '40px',
                  padding: '0 24px',
                  color: 'var(--color-text-primary)',
                  borderColor: 'var(--color-border)',
                }}
              >
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ConfigPage
