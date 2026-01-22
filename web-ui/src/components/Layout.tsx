// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect, useState } from 'react'
import { Layout, Menu, Button, message } from 'antd'
import { useNavigate, useLocation } from 'react-router-dom'
import client, { publicClient } from '../api/client'
import {
  DashboardOutlined,
  ProjectOutlined,
  LinkOutlined,
  SettingOutlined,
  KeyOutlined,
  FileTextOutlined,
  LogoutOutlined,
  HomeOutlined,
} from '@ant-design/icons'

const { Header, Sider, Content, Footer } = Layout

interface AppLayoutProps {
  children: React.ReactNode
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  // Initialize with token check - if token exists, assume authenticated initially to avoid flash
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(() => {
    const token = localStorage.getItem('kkartifact_token')
    return token ? true : null // If token exists, start as authenticated to avoid flash
  })

  // Verify authentication on mount and periodically
  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('kkartifact_token')
      if (!token) {
        // No token, use public endpoint to check if service is available
        try {
          await publicClient.get('/public/projects', { params: { limit: 1 } })
          // Service is available but user is not authenticated
          setIsAuthenticated(false)
          navigate('/login', { replace: true })
        } catch (error: any) {
          // Service might be down, but still redirect to login
          setIsAuthenticated(false)
          navigate('/login', { replace: true })
        }
        return
      }

      try {
        // Verify token is still valid using authenticated endpoint
        await client.get('/projects', { params: { limit: 1 } })
        setIsAuthenticated(true)
      } catch (error: any) {
        if (error.response?.status === 401) {
          // Token is invalid, try using public endpoint instead to avoid 401
          try {
            await publicClient.get('/public/projects', { params: { limit: 1 } })
            // Public endpoint works, but user is not authenticated
            localStorage.removeItem('kkartifact_token')
            setIsAuthenticated(false)
            navigate('/login', { replace: true })
            message.warning('您的会话已过期，请重新登录')
          } catch (publicError: any) {
            // Even public endpoint failed, remove token and redirect
            localStorage.removeItem('kkartifact_token')
            setIsAuthenticated(false)
            navigate('/login', { replace: true })
            message.warning('您的会话已过期，请重新登录')
          }
        } else {
          // For other errors (network, etc.), still assume authenticated if we have token
          // This prevents flash on network issues
          setIsAuthenticated(true)
        }
      }
    }

    checkAuth()
    // Check auth every 5 minutes
    const interval = setInterval(checkAuth, 5 * 60 * 1000)
    return () => clearInterval(interval)
  }, [navigate])

  const handleLogout = () => {
    localStorage.removeItem('kkartifact_token')
    navigate('/login', { replace: true })
      message.success('退出登录成功')
  }

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表盘',
    },
    {
      key: '/projects',
      icon: <ProjectOutlined />,
      label: '项目',
    },
    {
      key: '/webhooks',
      icon: <LinkOutlined />,
      label: 'Webhooks',
    },
    {
      key: '/tokens',
      icon: <KeyOutlined />,
      label: '令牌',
    },
    {
      key: '/audit-logs',
      icon: <FileTextOutlined />,
      label: '审计日志',
    },
    {
      key: '/config',
      icon: <SettingOutlined />,
      label: '配置',
    },
  ]

  const selectedKey = menuItems.find((item) => location.pathname.startsWith(item.key))?.key || '/dashboard'

  // Show loading state while checking authentication
  if (isAuthenticated === null) {
    return null // Or return a loading spinner
  }

  if (!isAuthenticated) {
    return null // Will redirect to login
  }

  return (
    <Layout style={{ minHeight: '100vh', background: 'var(--color-bg-tertiary)' }}>
      <Sider 
        collapsible 
        theme="light"
        width={240}
        style={{
          background: 'var(--color-bg-primary)',
          borderRight: '1px solid var(--color-border-light)',
        }}
      >
        <div
          style={{ 
            padding: '24px 20px', 
            display: 'flex', 
            alignItems: 'center', 
            gap: '12px', 
            cursor: 'pointer',
            color: 'var(--color-text-primary)',
            transition: 'background-color 0.2s',
            borderRadius: 'var(--radius-sm)',
            margin: '8px',
          }}
          onClick={() => navigate('/dashboard')}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = 'var(--color-bg-secondary)'
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = 'transparent'
          }}
        >
          <img src="/logo-icon.svg" alt="kkArtifact" style={{ width: '28px', height: '28px' }} />
          <span style={{ fontSize: '16px', fontWeight: 600, letterSpacing: '-0.2px' }}>kkArtifact</span>
        </div>
        <Menu
          theme="light"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{
            borderRight: 'none',
            background: 'transparent',
            padding: '8px',
          }}
        />
      </Sider>
      <Layout>
        <Header 
          style={{ 
            background: 'var(--color-bg-primary)', 
            padding: '0 32px', 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            borderBottom: '1px solid var(--color-border-light)',
            height: '64px',
            position: 'sticky',
            top: 0,
            zIndex: 100,
          }}
        >
          <h1 style={{ margin: 0, fontSize: '18px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.2px' }}>
            制品管理
          </h1>
          <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
            <Button 
              type="text" 
              icon={<HomeOutlined />} 
              onClick={() => navigate('/')}
              style={{
                height: '36px',
                padding: '0 12px',
                color: 'var(--color-text-secondary)',
              }}
            >
              版本清单
            </Button>
            <Button 
              type="text" 
              icon={<LogoutOutlined />} 
              onClick={handleLogout}
              style={{
                height: '36px',
                padding: '0 12px',
                color: 'var(--color-text-secondary)',
              }}
            >
              退出
            </Button>
          </div>
        </Header>
        <Content 
          style={{ 
            margin: '24px', 
            padding: 0,
            minHeight: 280,
          }}
        >
          <div
            style={{
              background: 'var(--color-bg-primary)',
              borderRadius: 'var(--radius-md)',
              padding: '32px',
              border: '1px solid var(--color-border-light)',
            }}
          >
            {children}
          </div>
        </Content>
        <Footer 
          style={{ 
            textAlign: 'center', 
            background: 'var(--color-bg-primary)', 
            borderTop: '1px solid var(--color-border-light)',
            padding: '20px 24px',
            color: 'var(--color-text-secondary)',
            fontSize: '13px',
          }}
        >
          系统运行部驱动
        </Footer>
      </Layout>
    </Layout>
  )
}

export default AppLayout

