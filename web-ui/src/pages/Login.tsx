// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect, useRef } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Form, Input, Button, Card, message, Typography } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { authApi } from '../api/auth'
import client from '../api/client'
import styles from './Login.module.css'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  const particlesRef = useRef<HTMLDivElement>(null)
  const connectionsRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Check if already logged in
    const token = localStorage.getItem('kkartifact_token')
    if (token) {
      // Verify token is valid by making a test request
      client
        .get('/projects', { params: { limit: 1, offset: 0 } })
        .then(() => {
          navigate('/dashboard')
        })
        .catch(() => {
          // Token invalid, clear it
          localStorage.removeItem('kkartifact_token')
        })
    }

    // Create particles
    if (particlesRef.current) {
      const particles = particlesRef.current
      particles.innerHTML = ''
      
      for (let i = 0; i < 20; i++) {
        const particle = document.createElement('div')
        particle.className = styles.particle
        particle.style.left = `${Math.random() * 100}%`
        particle.style.width = particle.style.height = `${Math.random() * 4 + 2}px`
        particle.style.animationDelay = `${Math.random() * 15}s`
        particle.style.animationDuration = `${Math.random() * 10 + 10}s`
        particles.appendChild(particle)
      }
    }

    // Create connection lines
    if (connectionsRef.current) {
      const connectionsContainer = connectionsRef.current
      connectionsContainer.innerHTML = ''
      
      for (let i = 0; i < 10; i++) {
        const connection = document.createElement('div')
        connection.className = styles.connection
        connection.style.top = `${Math.random() * 100}%`
        connection.style.left = `${Math.random() * 100}%`
        connection.style.width = `${Math.random() * 200 + 100}px`
        connection.style.animationDelay = `${Math.random() * 3}s`
        connection.style.transform = `rotate(${Math.random() * 360}deg)`
        connectionsContainer.appendChild(connection)
      }
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

      // Redirect to the specified page or dashboard
      const redirect = searchParams.get('redirect')
      navigate(redirect || '/dashboard', { replace: true })
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error || 'Login failed. Please check your credentials.'
      message.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.loginContainer}>
      <div className={styles.gridBackground}></div>
      <div className={styles.particles} ref={particlesRef}></div>
      <div className={styles.connections} ref={connectionsRef}></div>
      
      <div className={styles.contentWrapper}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            flex: 1,
          }}
        >
          <Card 
            className={styles.loginCard}
            style={{ 
              width: 400, 
              boxShadow: '0 8px 32px rgba(0,0,0,0.2)',
              background: 'rgba(255, 255, 255, 0.95)',
              backdropFilter: 'blur(10px)',
              border: '1px solid rgba(255, 255, 255, 0.3)',
            }}
          >
            <div style={{ textAlign: 'center', marginBottom: 32 }}>
              <div style={{ marginBottom: 16 }}>
                <img src="/logo-icon.svg" alt="kkArtifact" style={{ width: '64px', height: '64px' }} />
              </div>
              <Title level={2} style={{ marginBottom: 8 }}>kkArtifact</Title>
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
          className={styles.footer}
          style={{
            textAlign: 'center',
            padding: '16px',
            background: 'rgba(255, 255, 255, 0.9)',
            backdropFilter: 'blur(10px)',
            borderTop: '1px solid rgba(255, 255, 255, 0.3)',
          }}
        >
          本系统由kk驱动
        </div>
      </div>
    </div>
  )
}

export default LoginPage
