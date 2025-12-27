// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect, useState } from 'react'
import { Layout, Menu, Button, message } from 'antd'
import { useNavigate, useLocation } from 'react-router-dom'
import client from '../api/client'
import {
  DashboardOutlined,
  ProjectOutlined,
  LinkOutlined,
  SettingOutlined,
  KeyOutlined,
  FileTextOutlined,
  LogoutOutlined,
} from '@ant-design/icons'

const { Header, Sider, Content, Footer } = Layout

interface AppLayoutProps {
  children: React.ReactNode
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const [isAuthenticated, setIsAuthenticated] = useState(true)

  // Verify authentication on mount and periodically
  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('kkartifact_token')
      if (!token) {
        setIsAuthenticated(false)
        navigate('/login', { replace: true })
        return
      }

      try {
        // Verify token is still valid
        await client.get('/projects', { params: { limit: 1 } })
        setIsAuthenticated(true)
      } catch (error: any) {
        if (error.response?.status === 401) {
          localStorage.removeItem('kkartifact_token')
          setIsAuthenticated(false)
          navigate('/login', { replace: true })
          message.warning('Your session has expired. Please login again.')
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
    message.success('Logged out successfully')
  }

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
    },
    {
      key: '/projects',
      icon: <ProjectOutlined />,
      label: 'Projects',
    },
    {
      key: '/webhooks',
      icon: <LinkOutlined />,
      label: 'Webhooks',
    },
    {
      key: '/tokens',
      icon: <KeyOutlined />,
      label: 'Tokens',
    },
    {
      key: '/audit-logs',
      icon: <FileTextOutlined />,
      label: 'Audit Logs',
    },
    {
      key: '/config',
      icon: <SettingOutlined />,
      label: 'Configuration',
    },
  ]

  const selectedKey = menuItems.find((item) => location.pathname.startsWith(item.key))?.key || '/dashboard'

  if (!isAuthenticated) {
    return null // Will redirect to login
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible theme="dark">
        <div
          style={{ padding: '16px', color: 'white', fontSize: '18px', fontWeight: 'bold', cursor: 'pointer' }}
          onClick={() => navigate('/dashboard')}
        >
          kkArtifact
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
          <h1 style={{ margin: 0, fontSize: '20px' }}>Artifact Management</h1>
          <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>
            Logout
          </Button>
        </Header>
        <Content style={{ margin: '24px', padding: '24px', background: '#fff', minHeight: 280 }}>
          {children}
        </Content>
        <Footer style={{ textAlign: 'center', background: '#fff', borderTop: '1px solid #f0f0f0' }}>
          本系统由kk驱动
        </Footer>
      </Layout>
    </Layout>
  )
}

export default AppLayout

