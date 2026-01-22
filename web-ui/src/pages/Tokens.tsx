// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Space, Popconfirm, message, Tag, Modal, Form, Input, Select } from 'antd'
import { tokensApi, Token, CreateTokenRequest } from '../api/tokens'
import type { ColumnsType } from 'antd/es/table'
import { CopyOutlined, DeleteOutlined, PlusOutlined, EyeOutlined, EyeInvisibleOutlined } from '@ant-design/icons'

const { Option } = Select

// Local storage key for storing token values (mapping token ID to token value)
const TOKEN_STORAGE_KEY = 'kkartifact_token_values'

const TokensPage: React.FC = () => {
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [createdToken, setCreatedToken] = useState<string | null>(null)
  const [visibleTokens, setVisibleTokens] = useState<Set<number>>(new Set())
  const [storedTokens, setStoredTokens] = useState<Record<number, string>>({})
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  // Load stored tokens from localStorage on mount
  useEffect(() => {
    try {
      const stored = localStorage.getItem(TOKEN_STORAGE_KEY)
      if (stored) {
        setStoredTokens(JSON.parse(stored))
      }
    } catch (e) {
      if (import.meta.env.DEV) {
        console.error('Failed to load stored tokens:', e)
      }
    }
  }, [])

  // Save token to localStorage when created
  const saveTokenToStorage = (tokenId: number, tokenValue: string) => {
    try {
      const updated = { ...storedTokens, [tokenId]: tokenValue }
      setStoredTokens(updated)
      localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(updated))
    } catch (e) {
      if (import.meta.env.DEV) {
        console.error('Failed to save token to storage:', e)
      }
    }
  }

  const { data, isLoading, error } = useQuery({
    queryKey: ['tokens'],
    queryFn: () => tokensApi.list().then((res) => res.data),
  })

  // Log error for debugging
  if (error) {
    if (import.meta.env.DEV) {
      console.error('Failed to load tokens:', error)
    }
  }

  const createMutation = useMutation({
    mutationFn: (data: CreateTokenRequest) => tokensApi.create(data),
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['tokens'] })
      setIsModalVisible(false)
      form.resetFields()
      if (response.data.token && response.data.id) {
        // Save token to localStorage for later viewing
        saveTokenToStorage(response.data.id, response.data.token)
        setCreatedToken(response.data.token)
        message.success('Token created successfully!')
      }
    },
    onError: () => {
      message.error('Failed to create token')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => tokensApi.delete(id),
    onSuccess: (_, id) => {
      // Remove token from localStorage when deleted
      try {
        const updated = { ...storedTokens }
        delete updated[id]
        setStoredTokens(updated)
        localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(updated))
      } catch (e) {
        if (import.meta.env.DEV) {
          console.error('Failed to remove token from storage:', e)
        }
      }
      queryClient.invalidateQueries({ queryKey: ['tokens'] })
      message.success('Token deleted successfully')
    },
    onError: () => {
      message.error('Failed to delete token')
    },
  })

  const handleCreate = () => {
    setCreatedToken(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleSubmit = () => {
    form.validateFields().then((values) => {
      const data: CreateTokenRequest = {
        name: values.name,
        permissions: values.permissions || ['pull', 'push', 'publish'],
      }
      createMutation.mutate(data)
    })
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    message.success('令牌已复制到剪贴板')
  }

  const toggleTokenVisibility = (id: number) => {
    const newVisible = new Set(visibleTokens)
    if (newVisible.has(id)) {
      newVisible.delete(id)
    } else {
      newVisible.add(id)
    }
    setVisibleTokens(newVisible)
  }

  const maskToken = (token: string) => {
    if (!token || token.length <= 8) {
      return '****'
    }
    return `${token.substring(0, 4)}${'*'.repeat(Math.max(0, token.length - 8))}${token.substring(token.length - 4)}`
  }

  const columns: ColumnsType<Token> = [
    {
      title: <span style={{ fontWeight: 600 }}>名称</span>,
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <span style={{ fontWeight: 500, color: '#1a1a1a' }}>{text}</span>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>令牌</span>,
      key: 'token',
      width: 400,
      render: (_, record) => {
        const storedToken = storedTokens[record.id]
        const hasToken = !!storedToken
        const isVisible = visibleTokens.has(record.id)
        const displayToken = hasToken
          ? isVisible
            ? storedToken
            : maskToken(storedToken)
          : '令牌不可用'

        return (
          <Space.Compact style={{ width: '100%' }}>
            <Input
              style={{ width: 'calc(100% - 120px)', fontFamily: 'monospace', fontSize: '12px' }}
              value={displayToken}
              readOnly
            />
            {hasToken && (
              <>
                <Button
                  icon={isVisible ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                  onClick={() => toggleTokenVisibility(record.id)}
                  title={isVisible ? '隐藏令牌' : '显示令牌'}
                />
                <Button
                  icon={<CopyOutlined />}
                  onClick={() => {
                    if (storedToken) {
                      copyToClipboard(storedToken)
                    }
                  }}
                  title="复制令牌"
                />
              </>
            )}
          </Space.Compact>
        )
      },
    },
    {
      title: <span style={{ fontWeight: 600 }}>权限</span>,
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <Space>
          {permissions.map((perm) => {
            const labels: Record<string, string> = {
              pull: '拉取',
              push: '推送',
              publish: '发布',
            }
            return <Tag key={perm}>{labels[perm] || perm}</Tag>
          })}
        </Space>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>过期时间</span>,
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (expiresAt?: string) => (
        <span style={{ color: '#8c8c8c', fontSize: '14px' }}>
          {expiresAt ? new Date(expiresAt).toLocaleString('zh-CN') : '永不过期'}
        </span>
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
      width: 120,
      render: (_, record) => (
        <Space>
          <Popconfirm
            title="确定要删除此令牌吗？"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button 
              type="link" 
              size="small" 
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
              令牌
            </h2>
            <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
              管理 API 访问令牌
            </div>
          </div>
          <Button 
            type="primary" 
            icon={<PlusOutlined />} 
            onClick={handleCreate}
            style={{
              borderRadius: '6px',
              height: '36px',
              padding: '0 16px',
              fontWeight: 500,
            }}
          >
            创建令牌
          </Button>
        </div>
      </div>
      {error && (
        <div style={{ 
          marginBottom: '16px', 
          padding: '16px', 
          background: '#fff1f0', 
          border: '1px solid #ffccc7', 
          borderRadius: '6px',
          color: '#cf1322',
        }}>
          加载令牌失败，请重试。
        </div>
      )}
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
          dataSource={data || []}
          loading={isLoading}
          rowKey="id"
        />
      </div>
      <Modal
        title="创建令牌"
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalVisible(false)
          setCreatedToken(null)
          form.resetFields()
        }}
        okText="创建"
        cancelText="取消"
        confirmLoading={createMutation.isPending}
      >
        {createdToken ? (
          <div>
            <p><strong>令牌创建成功！</strong></p>
            <p>请立即复制此令牌。您将无法再次查看它：</p>
            <Input.Group compact>
              <Input
                style={{ width: 'calc(100% - 100px)' }}
                value={createdToken}
                readOnly
              />
              <Button
                icon={<CopyOutlined />}
                onClick={() => copyToClipboard(createdToken)}
              >
                复制
              </Button>
            </Input.Group>
            <Button
              style={{ marginTop: 16, width: '100%' }}
              onClick={() => {
                setCreatedToken(null)
                setIsModalVisible(false)
              }}
            >
              完成
            </Button>
          </div>
        ) : (
          <Form form={form} layout="vertical">
            <Form.Item
              name="name"
              label="名称"
              rules={[{ required: true, message: '请输入令牌名称！' }]}
            >
              <Input placeholder="例如：web-ui-token" />
            </Form.Item>
            <Form.Item
              name="permissions"
              label="权限"
              rules={[{ required: true, message: '请选择权限！' }]}
              initialValue={['pull', 'push', 'publish']}
            >
              <Select mode="multiple" placeholder="选择权限">
                <Option value="pull">拉取</Option>
                <Option value="push">推送</Option>
                <Option value="publish">发布</Option>
                <Option value="admin">管理员</Option>
              </Select>
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  )
}

export default TokensPage

