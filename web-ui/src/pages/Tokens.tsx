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
      console.error('Failed to load stored tokens:', e)
    }
  }, [])

  // Save token to localStorage when created
  const saveTokenToStorage = (tokenId: number, tokenValue: string) => {
    try {
      const updated = { ...storedTokens, [tokenId]: tokenValue }
      setStoredTokens(updated)
      localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(updated))
    } catch (e) {
      console.error('Failed to save token to storage:', e)
    }
  }

  const { data, isLoading, error } = useQuery({
    queryKey: ['tokens'],
    queryFn: () => tokensApi.list().then((res) => res.data),
  })

  // Log error for debugging
  if (error) {
    console.error('Failed to load tokens:', error)
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
        console.error('Failed to remove token from storage:', e)
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
        permissions: values.permissions || ['pull', 'push', 'promote'],
      }
      createMutation.mutate(data)
    })
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    message.success('Token copied to clipboard')
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
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Token',
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
          : 'Token not available'

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
                  title={isVisible ? 'Hide token' : 'Show token'}
                />
                <Button
                  icon={<CopyOutlined />}
                  onClick={() => {
                    if (storedToken) {
                      copyToClipboard(storedToken)
                    }
                  }}
                  title="Copy token"
                />
              </>
            )}
          </Space.Compact>
        )
      },
    },
    {
      title: 'Permissions',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <Space>
          {permissions.map((perm) => (
            <Tag key={perm}>{perm}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: 'Expires At',
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (expiresAt?: string) => expiresAt ? new Date(expiresAt).toLocaleString() : 'Never',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Popconfirm
            title="Are you sure to delete this token?"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
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
        <h2>Tokens</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Create Token
        </Button>
      </div>
      {error && (
        <div style={{ marginBottom: 16, padding: 16, background: '#fff1f0', border: '1px solid #ffccc7', borderRadius: 4 }}>
          Failed to load tokens. Please try again.
        </div>
      )}
      <Table
        columns={columns}
        dataSource={data || []}
        loading={isLoading}
        rowKey="id"
      />
      <Modal
        title="Create Token"
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalVisible(false)
          setCreatedToken(null)
          form.resetFields()
        }}
        okText="Create"
        cancelText="Cancel"
        confirmLoading={createMutation.isPending}
      >
        {createdToken ? (
          <div>
            <p><strong>Token created successfully!</strong></p>
            <p>Please copy this token now. You won't be able to see it again:</p>
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
                Copy
              </Button>
            </Input.Group>
            <Button
              style={{ marginTop: 16, width: '100%' }}
              onClick={() => {
                setCreatedToken(null)
                setIsModalVisible(false)
              }}
            >
              Done
            </Button>
          </div>
        ) : (
          <Form form={form} layout="vertical">
            <Form.Item
              name="name"
              label="Name"
              rules={[{ required: true, message: 'Please input token name!' }]}
            >
              <Input placeholder="e.g., web-ui-token" />
            </Form.Item>
            <Form.Item
              name="permissions"
              label="Permissions"
              rules={[{ required: true, message: 'Please select permissions!' }]}
              initialValue={['pull', 'push', 'promote']}
            >
              <Select mode="multiple" placeholder="Select permissions">
                <Option value="pull">Pull</Option>
                <Option value="push">Push</Option>
                <Option value="promote">Promote</Option>
                <Option value="admin">Admin</Option>
              </Select>
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  )
}

export default TokensPage

