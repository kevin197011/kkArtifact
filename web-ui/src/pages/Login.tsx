// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect, useRef } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Form, Input, Button, Card, message, Typography, Spin } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { authApi } from '../api/auth'
import client from '../api/client'
import styles from './Login.module.css'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  const [checkingToken, setCheckingToken] = useState(true)
  const particlesRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Check if already logged in
    const token = localStorage.getItem('kkartifact_token')
    if (token) {
      // Verify token is valid by making a test request
      client
        .get('/projects', { params: { limit: 1, offset: 0 } })
        .then(() => {
          navigate('/dashboard', { replace: true })
        })
        .catch(() => {
          // Token invalid, clear it
          localStorage.removeItem('kkartifact_token')
          setCheckingToken(false)
        })
    } else {
      setCheckingToken(false)
    }
  }, [navigate])

  // Create floating particles for dynamic effect
  useEffect(() => {
    if (!particlesRef.current) return

    const container = particlesRef.current
    container.innerHTML = ''

    // Create multiple floating particles
    for (let i = 0; i < 15; i++) {
      const particle = document.createElement('div')
      particle.className = styles.particle
      particle.style.left = `${Math.random() * 100}%`
      particle.style.top = `${Math.random() * 100}%`
      particle.style.width = particle.style.height = `${Math.random() * 6 + 4}px`
      particle.style.animationDelay = `${Math.random() * 20}s`
      particle.style.animationDuration = `${Math.random() * 15 + 20}s`
      container.appendChild(particle)
    }

    // Create additional floating orbs
    for (let i = 0; i < 3; i++) {
      const orb = document.createElement('div')
      orb.className = styles.floatingOrb
      orb.style.left = `${Math.random() * 100}%`
      orb.style.top = `${Math.random() * 100}%`
      orb.style.width = orb.style.height = `${Math.random() * 200 + 150}px`
      orb.style.animationDelay = `${Math.random() * 10}s`
      orb.style.animationDuration = `${Math.random() * 20 + 25}s`
      container.appendChild(orb)
    }
  }, [])

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const response = await authApi.login({
        username: values.username,
        password: values.password,
      })

      // Store token
      localStorage.setItem('kkartifact_token', response.data.token)
      message.success(`欢迎，${response.data.name}！`)

      // Redirect to the specified page or dashboard
      const redirect = searchParams.get('redirect')
      navigate(redirect || '/dashboard', { replace: true })
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error || '登录失败，请检查您的凭据。'
      message.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // Show loading state while checking token
  if (checkingToken) {
    return (
      <div className={styles.loginContainer}>
        <div className={styles.particles} ref={particlesRef}></div>
        <div className={styles.contentWrapper}>
          <div
            style={{
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              flex: 1,
              minHeight: '100vh',
            }}
          >
            <Card 
              className={styles.loginCard}
              style={{ 
                width: '100%',
                maxWidth: 360,
              }}
              bodyStyle={{ padding: '32px' }}
            >
              <div style={{ textAlign: 'center', padding: '50px' }}>
                <Spin size="large" />
              </div>
            </Card>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={styles.loginContainer}>
      <div className={styles.particles} ref={particlesRef}></div>
      <div className={styles.contentWrapper}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            flex: 1,
            height: '100vh',
            padding: '24px',
          }}
        >
            <Card 
              className={styles.loginCard}
              style={{ 
                width: '100%',
                maxWidth: 360,
              }}
              bodyStyle={{ padding: '32px' }}
            >
              <div style={{ textAlign: 'center', marginBottom: '32px' }}>
                <div 
                  style={{ 
                    marginBottom: '16px',
                    display: 'inline-flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '56px',
                    height: '56px',
                    borderRadius: '12px',
                    background: 'linear-gradient(135deg, #165dff 0%, #4080ff 100%)',
                    boxShadow: '0 4px 12px rgba(22, 93, 255, 0.2)',
                    transition: 'transform 0.3s ease',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.transform = 'scale(1.05) rotate(5deg)'
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.transform = 'scale(1) rotate(0deg)'
                  }}
                >
                  <img 
                    src="/logo-icon.svg" 
                    alt="kkArtifact" 
                    style={{ 
                      width: '36px', 
                      height: '36px',
                      filter: 'brightness(0) invert(1)',
                    }} 
                  />
                </div>
                <Title 
                  level={3} 
                  style={{ 
                    margin: 0, 
                    marginBottom: '6px', 
                    fontSize: '22px', 
                    fontWeight: 600, 
                    color: '#1d2129',
                    letterSpacing: '-0.3px',
                  }}
                >
                  kkArtifact
                </Title>
                <Text 
                  type="secondary" 
                  style={{ 
                    fontSize: '13px', 
                    color: '#86909c',
                    fontWeight: 400,
                  }}
                >
                  登录您的账户
                </Text>
              </div>

            <Form
              name="login"
              onFinish={onFinish}
              autoComplete="off"
              size="large"
              layout="vertical"
            >
              <Form.Item
                label={<span style={{ fontSize: '13px', fontWeight: 500, color: '#1d2129' }}>用户名</span>}
                name="username"
                rules={[{ required: true, message: '请输入用户名！' }]}
                style={{ marginBottom: '20px' }}
              >
                <Input
                  prefix={<UserOutlined style={{ color: '#c9cdd4' }} />}
                  placeholder="请输入用户名"
                  autoComplete="username"
                  style={{
                    height: '40px',
                    borderRadius: '8px',
                  }}
                />
              </Form.Item>

              <Form.Item
                label={<span style={{ fontSize: '13px', fontWeight: 500, color: '#1d2129' }}>密码</span>}
                name="password"
                rules={[{ required: true, message: '请输入密码！' }]}
                style={{ marginBottom: '24px' }}
              >
                <Input.Password
                  prefix={<LockOutlined style={{ color: '#c9cdd4' }} />}
                  placeholder="请输入密码"
                  autoComplete="current-password"
                  style={{
                    height: '40px',
                    borderRadius: '8px',
                  }}
                />
              </Form.Item>

              <Form.Item style={{ marginBottom: 0 }}>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  block 
                  loading={loading} 
                  style={{ 
                    height: '40px', 
                    fontWeight: 500,
                    fontSize: '14px',
                    borderRadius: '8px',
                    boxShadow: '0 2px 8px rgba(22, 93, 255, 0.2)',
                    transition: 'all 0.2s ease',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.boxShadow = '0 4px 12px rgba(22, 93, 255, 0.3)'
                    e.currentTarget.style.transform = 'translateY(-1px)'
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.boxShadow = '0 2px 8px rgba(22, 93, 255, 0.2)'
                    e.currentTarget.style.transform = 'translateY(0)'
                  }}
                >
                  登录
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </div>
        <div
          className={styles.footer}
          style={{
            textAlign: 'center',
            padding: '20px',
            color: '#86909c',
            fontSize: '13px',
            background: 'transparent',
          }}
        >
          系统运行部驱动
        </div>
      </div>
    </div>
  )
}

export default LoginPage
