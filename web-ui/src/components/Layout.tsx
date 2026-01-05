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
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible theme="dark">
        <div
          style={{ 
            padding: '16px', 
            display: 'flex', 
            alignItems: 'center', 
            gap: '12px', 
            cursor: 'pointer',
            color: 'white'
          }}
          onClick={() => navigate('/dashboard')}
        >
          <img src="/logo-icon.svg" alt="kkArtifact" style={{ width: '32px', height: '32px' }} />
          <span style={{ fontSize: '18px', fontWeight: 'bold' }}>kkArtifact</span>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h1 style={{ margin: 0, fontSize: '20px' }}>制品管理</h1>
          <div style={{ display: 'flex', gap: '16px', alignItems: 'center' }}>
            <Button type="default" icon={<HomeOutlined />} onClick={() => navigate('/')}>
              查看版本清单
            </Button>
            <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>
              退出登录
            </Button>
          </div>
        </Header>
        <Content style={{ margin: '24px', padding: '24px', background: '#fff', minHeight: 280 }}>
          {children}
        </Content>
        <Footer style={{ textAlign: 'center', background: '#fff', borderTop: '1px solid #f0f0f0' }}>
          本系统由系统部驱动
        </Footer>
      </Layout>
    </Layout>
  )
}

export default AppLayout

