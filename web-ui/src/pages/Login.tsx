// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Form, Input, Button, Card, message, Typography } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { authApi } from '../api/auth'
import client from '../api/client'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    // Check if already logged in
    const token = localStorage.getItem('kkartifact_token')
    if (token) {
      // Verify token is valid by making a test request
      client
        .get('/projects', { params: { limit: 1, offset: 0 } })
        .then(() => {
          navigate('/projects')
        })
        .catch(() => {
          // Token invalid, clear it
          localStorage.removeItem('kkartifact_token')
        })
    }
  }, [navigate])

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const response = await authApi.login({
        username: values.username,
        password: values.password,
      })

      // Store token
      localStorage.setItem('kkartifact_token', response.data.token)
      message.success(`Welcome, ${response.data.name}!`)

      // Redirect to projects page
      navigate('/projects')
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error || 'Login failed. Please check your credentials.'
      message.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        background: '#f0f2f5',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          flex: 1,
        }}
      >
        <Card style={{ width: 400, boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}>
          <div style={{ textAlign: 'center', marginBottom: 32 }}>
            <Title level={2}>kkArtifact</Title>
            <Text type="secondary">Sign in to your account</Text>
          </div>

          <Form
            name="login"
            onFinish={onFinish}
            autoComplete="off"
            size="large"
            layout="vertical"
          >
            <Form.Item
              label="Username"
              name="username"
              rules={[{ required: true, message: 'Please input your username!' }]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="Username"
                autoComplete="username"
              />
            </Form.Item>

            <Form.Item
              label="Password"
              name="password"
              rules={[{ required: true, message: 'Please input your password!' }]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="Password"
                autoComplete="current-password"
              />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" block loading={loading}>
                Sign In
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
      <div
        style={{
          textAlign: 'center',
          padding: '16px',
          background: '#fff',
          borderTop: '1px solid #f0f0f0',
        }}
      >
        本系统由kk驱动
      </div>
    </div>
  )
}

export default LoginPage
