// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useMemo, useCallback, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Tree, Input, Empty, Spin, Button, Typography, Card, Collapse, Tag, Space, message } from 'antd'
import { SearchOutlined, FolderOutlined, AppstoreOutlined, LoginOutlined, DownloadOutlined, CopyOutlined, CalendarOutlined, CheckCircleOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { publicProjectsApi, Project, App, Version } from '../api/projects'
import { downloadsApi } from '../api/downloads'
import ThemeToggle from '../components/ThemeToggle'
import type { DataNode } from 'antd/es/tree'
import styles from './InventoryPage.module.css'

const { Panel } = Collapse

const { Title } = Typography

type TreeDataNode = DataNode & {
  project?: Project
  app?: App
  version?: Version
  isProject?: boolean
  isApp?: boolean
  isVersion?: boolean
}

// Simple debounce hook
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value)

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])

  return debouncedValue
}

const InventoryPage: React.FC = () => {
  const navigate = useNavigate()
  const [searchTerm, setSearchTerm] = useState('')
  const debouncedSearchTerm = useDebounce(searchTerm, 300)
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([])
  
  // Get full script URL for install commands
  const getScriptUrl = (filename: string) => {
    const scriptPath = downloadsApi.downloadScript(filename)
    // If scriptPath is already a full URL, return it
    if (scriptPath.startsWith('http://') || scriptPath.startsWith('https://')) {
      return scriptPath
    }
    // Otherwise, prepend current origin
    return `${window.location.origin}${scriptPath}`
  }

  // Get current server URL (origin) for server_url environment variable
  const getCurrentServerUrl = () => {
    return window.location.origin
  }

  // Format date for display
  const formatDate = useCallback((dateString: string): string => {
    try {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      })
    } catch {
      return ''
    }
  }, [])

  // Copy version to clipboard
  const handleCopyVersion = useCallback((version: string, e: React.MouseEvent) => {
    e.stopPropagation() // Prevent tree node expansion
    navigator.clipboard.writeText(version).then(() => {
      message.success('版本号已复制到剪贴板')
    }).catch(() => {
      message.error('复制失败')
    })
  }, [])


  // Fetch all projects
  const { data: allProjects, isLoading: projectsLoading } = useQuery({
    queryKey: ['public-projects', 'all'],
    queryFn: () => publicProjectsApi.list(10000, 0).then((res) => res.data),
    retry: 1,
  })

  // Fetch apps for all projects
  const { data: allAppsData } = useQuery({
    queryKey: ['public-all-apps', allProjects?.map((p) => p.name)],
    queryFn: async () => {
      if (!allProjects || allProjects.length === 0) return {}
      const appsPromises = allProjects.map(async (project) => {
        try {
          const apps = await publicProjectsApi.getApps(project.name, 1000, 0).then((res) => res.data)
          return { projectName: project.name, apps }
          } catch (error) {
            return { projectName: project.name, apps: [] }
        }
      })
      const appsArrays = await Promise.all(appsPromises)
      const appsMap: Record<string, App[]> = {}
      appsArrays.forEach(({ projectName, apps }) => {
        appsMap[projectName] = apps
      })
      return appsMap
    },
    enabled: !!allProjects && allProjects.length > 0,
  })

  // Fetch versions for all apps
  const { data: allVersionsData } = useQuery({
    queryKey: ['public-all-versions', allAppsData],
    queryFn: async () => {
      if (!allAppsData) return {}
      const versionsMap: Record<string, Version[]> = {} // Key: "projectName/appName"
      for (const projectName of Object.keys(allAppsData)) {
        const apps = allAppsData[projectName]
        for (const app of apps) {
          try {
            const versions = await publicProjectsApi
              .getVersions(projectName, app.name, 1000, 0)
              .then((res) => res.data)
            versionsMap[`${projectName}/${app.name}`] = versions
          } catch (error) {
            versionsMap[`${projectName}/${app.name}`] = []
          }
        }
      }
      return versionsMap
    },
    enabled: !!allAppsData && Object.keys(allAppsData).length > 0,
  })

  // Build tree data structure with filtering
  const treeData = useMemo(() => {
    if (!allProjects) return []

    const searchTermLower = debouncedSearchTerm.toLowerCase().trim()
    const hasSearch = searchTermLower.length > 0

    // Filter projects
    let projectsToShow = allProjects
    if (hasSearch) {
      projectsToShow = allProjects.filter((project) =>
        project.name.toLowerCase().includes(searchTermLower)
      )
    }

    // Get apps and versions for projects
    const appsByProject: Record<string, App[]> = {}
    const versionsByApp: Record<string, Version[]> = {}
    const projectsWithMatchingItems = new Set<string>()

    if (hasSearch) {
      // When searching, check all projects for matching apps and versions
      for (const project of allProjects) {
        const apps = allAppsData?.[project.name] || []
        const matchingApps: App[] = []
        apps.forEach((app) => {
          const appMatches = app.name.toLowerCase().includes(searchTermLower)
          const versions = allVersionsData?.[`${project.name}/${app.name}`] || []
          const matchingVersions = versions.filter((v) =>
            v.version.toLowerCase().includes(searchTermLower)
          )

          if (appMatches || matchingVersions.length > 0) {
            matchingApps.push(app)
            if (matchingVersions.length > 0) {
              versionsByApp[`${project.name}/${app.name}`] = matchingVersions
            } else {
              versionsByApp[`${project.name}/${app.name}`] = versions
            }
            projectsWithMatchingItems.add(project.name)
          }
        })

        if (matchingApps.length > 0) {
          appsByProject[project.name] = matchingApps
        }
      }

      // Also include projects that match by name
      projectsToShow.forEach((p) => projectsWithMatchingItems.add(p.name))
    } else {
      // When not searching, show all data
      if (allAppsData) {
        Object.assign(appsByProject, allAppsData)
      }
      if (allVersionsData) {
        Object.assign(versionsByApp, allVersionsData)
      }
      projectsToShow.forEach((p) => projectsWithMatchingItems.add(p.name))
    }

    // Build tree nodes
    const nodes: TreeDataNode[] = Array.from(projectsWithMatchingItems)
      .map((projectName) => {
        const project = allProjects.find((p) => p.name === projectName)
        if (!project) return null as any // Will be filtered out

        const projectKey = `project-${project.id}`
        const apps = appsByProject[project.name] || []

        const appNodes: TreeDataNode[] = apps.map((app) => {
          const appKey = `app-${app.id}`
          const versions = versionsByApp[`${project.name}/${app.name}`] || []

          const versionNodes: TreeDataNode[] = versions.map((version) => ({
            key: `version-${version.id}`,
            title: (
              <div
                className={styles.versionNode}
                onClick={(e) => {
                  e.stopPropagation()
                  // Public inventory page - no navigation, just display information
                }}
              >
                <div className={styles.versionNodeContent}>
                  <span className={styles.versionNumber}>{version.version}</span>
                  {version.is_published && (
                    <Tag 
                      color="success" 
                      icon={<CheckCircleOutlined />}
                      className={styles.publishedTag}
                    >
                      已发布
                    </Tag>
                  )}
                  {version.created_at && (
                    <span className={styles.versionDate}>
                      <CalendarOutlined style={{ fontSize: '10px', marginRight: '3px' }} />
                      {formatDate(version.created_at)}
                    </span>
                  )}
                  <CopyOutlined
                    className={styles.copyIcon}
                    onClick={(e) => handleCopyVersion(version.version, e)}
                    title="复制版本号"
                  />
                </div>
              </div>
            ),
            isLeaf: true,
            version,
            isVersion: true,
            app,
            project,
          }))

          return {
            key: appKey,
            title: (
              <div
                className={styles.appNode}
                onClick={(e) => {
                  e.stopPropagation()
                  // Public inventory page - no navigation, just display information
                }}
              >
                <AppstoreOutlined className={styles.appIcon} />
                <span className={styles.appName}>{app.name}</span>
                {versions.length > 0 && (
                  <span className={styles.versionCount}>{versions.length} 个版本</span>
                )}
              </div>
            ),
            children: versionNodes.length > 0 ? versionNodes : undefined,
            isLeaf: versionNodes.length === 0,
            app,
            isApp: true,
            project,
          }
        })

        return {
          key: projectKey,
            title: (
              <div
                className={styles.projectNode}
                onClick={(e) => {
                  e.stopPropagation()
                  // Public inventory page - no navigation, just display information
                }}
              >
                <FolderOutlined className={styles.projectIcon} />
                <span className={styles.projectName}>{project.name}</span>
                {apps.length > 0 && (
                  <span className={styles.appCount}>{apps.length} 个应用</span>
                )}
              </div>
            ),
          children: appNodes.length > 0 ? appNodes : undefined,
          isLeaf: false,
          project,
          isProject: true,
        }
      })
      .filter((node) => node !== null)
      .sort((a, b) => {
        const nameA = a?.project?.name || ''
        const nameB = b?.project?.name || ''
        return nameA.localeCompare(nameB)
      }) as TreeDataNode[]

    return nodes
  }, [allProjects, allAppsData, allVersionsData, debouncedSearchTerm, formatDate])

  // Auto-expand when searching
  useEffect(() => {
    if (debouncedSearchTerm.trim() && treeData.length > 0) {
      const keysToExpand: React.Key[] = []
      const collectKeys = (nodes: TreeDataNode[]) => {
        nodes.forEach((node) => {
          keysToExpand.push(node.key)
          if (node.children) {
            collectKeys(node.children as TreeDataNode[])
          }
        })
      }
      collectKeys(treeData)
      setExpandedKeys(keysToExpand)
    } else if (!debouncedSearchTerm.trim()) {
      setExpandedKeys([])
    }
  }, [debouncedSearchTerm, treeData])

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value)
  }, [])

  const isLoading = projectsLoading

  const emptyText = debouncedSearchTerm.trim()
    ? `没有匹配 "${debouncedSearchTerm}" 的项目、应用或版本。请尝试其他关键词。`
    : '暂无项目'

  // Fetch agent version info
  const { data: agentVersionInfo } = useQuery({
    queryKey: ['agent-version'],
    queryFn: () => downloadsApi.getAgentVersionInfo().then((res) => res.data),
    retry: 1,
  })

  return (
    <div className={styles.inventoryContainer}>
      {/* Content wrapper */}
      <div className={styles.contentWrapper}>
        {/* Elegant header */}
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <img 
                src="/logo-icon.svg" 
                alt="kkArtifact" 
                style={{ 
                  width: '32px', 
                  height: '32px',
                  filter: 'drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1))',
                }} 
              />
              <div>
                <Title level={2} className={styles.title} style={{ margin: 0 }}>
                  制品清单
                </Title>
                <div style={{ fontSize: '13px', color: 'var(--color-text-secondary)', marginTop: '4px' }}>
                  浏览所有项目和版本
                </div>
              </div>
            </div>
          </div>
          <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
            <ThemeToggle scope="frontend" type="text" />
            <Button
              type="primary"
              icon={<LoginOutlined />}
              onClick={() => navigate('/login')}
              style={{
                height: '40px',
                padding: '0 20px',
                fontWeight: 500,
                fontSize: '14px',
                borderRadius: '8px',
                boxShadow: '0 2px 8px rgba(22, 93, 255, 0.2)',
              }}
            >
              登录后台
            </Button>
          </div>
        </div>

      {/* Agent Download Section */}
      {agentVersionInfo && agentVersionInfo.binaries.length > 0 && (
        <Card 
          style={{ 
            marginBottom: 32,
            borderRadius: 'var(--radius-lg)',
            border: '1px solid var(--color-border-light)',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.04)',
          }}
        >
          <Collapse defaultActiveKey={[]} ghost>
            <Panel 
              header={
                <Space>
                  <DownloadOutlined style={{ color: 'var(--color-primary)', fontSize: '16px' }} />
                  <span style={{ fontWeight: 600, fontSize: '15px', color: 'var(--color-text-primary)' }}>
                    安装 agent 客户端工具
                  </span>
                </Space>
              } 
              key="download"
            >
              <div className={styles.downloadContent}>
                <div className={styles.downloadDescription}>
                  <p>kkartifact-agent 是一个命令行工具，用于推送和拉取制品。支持并发传输、断点续传等特性。</p>
                  <div style={{ marginTop: '20px', marginBottom: '20px' }}>
                    <h4 style={{ marginBottom: '12px' }}>一键安装（推荐）：</h4>
                    <div style={{ marginBottom: '16px' }}>
                      <div style={{ marginBottom: '12px' }}>
                        <strong style={{ fontSize: '14px', color: 'var(--color-text-primary)' }}>Unix/Linux/macOS：</strong>
                        <div style={{ 
                          background: 'var(--color-bg-secondary)', 
                          padding: '12px', 
                          borderRadius: '6px', 
                          fontSize: '13px',
                          fontFamily: 'monospace',
                          marginTop: '8px',
                          position: 'relative',
                        }}>
                          <code style={{ color: 'var(--color-text-primary)' }}>
                            {`curl -fsSL ${getScriptUrl('install-agent.sh')} | server_url="${getCurrentServerUrl()}" bash`}
                          </code>
                          <CopyOutlined
                            style={{ 
                              position: 'absolute',
                              right: '12px',
                              top: '12px',
                              fontSize: '14px', 
                              color: 'var(--color-text-tertiary)', 
                              cursor: 'pointer',
                              padding: '4px',
                              borderRadius: '4px',
                              transition: 'all 0.2s ease',
                            }}
                            onClick={() => {
                              const scriptUrl = getScriptUrl('install-agent.sh')
                              const serverUrl = getCurrentServerUrl()
                              const cmd = `curl -fsSL ${scriptUrl} | server_url="${serverUrl}" bash`
                              navigator.clipboard.writeText(cmd)
                              message.success('命令已复制到剪贴板')
                            }}
                            title="复制命令"
                            onMouseEnter={(e) => {
                              e.currentTarget.style.color = 'var(--color-primary)'
                              e.currentTarget.style.backgroundColor = 'var(--color-primary-light)'
                            }}
                            onMouseLeave={(e) => {
                              e.currentTarget.style.color = 'var(--color-text-tertiary)'
                              e.currentTarget.style.backgroundColor = 'transparent'
                            }}
                          />
                        </div>
                      </div>
                      <div>
                        <strong style={{ fontSize: '14px', color: 'var(--color-text-primary)' }}>Windows (PowerShell)：</strong>
                        <div style={{ 
                          background: 'var(--color-bg-secondary)', 
                          padding: '12px', 
                          borderRadius: '6px', 
                          fontSize: '13px',
                          fontFamily: 'monospace',
                          marginTop: '8px',
                          position: 'relative',
                        }}>
                          <code style={{ color: 'var(--color-text-primary)' }}>
                            {`$env:server_url="${getCurrentServerUrl()}"; irm ${getScriptUrl('install-agent.ps1')} | iex`}
                          </code>
                          <CopyOutlined
                            style={{ 
                              position: 'absolute',
                              right: '12px',
                              top: '12px',
                              fontSize: '14px', 
                              color: 'var(--color-text-tertiary)', 
                              cursor: 'pointer',
                              padding: '4px',
                              borderRadius: '4px',
                              transition: 'all 0.2s ease',
                            }}
                            onClick={() => {
                              const scriptUrl = getScriptUrl('install-agent.ps1')
                              const serverUrl = getCurrentServerUrl()
                              const cmd = `$env:server_url="${serverUrl}"; irm ${scriptUrl} | iex`
                              navigator.clipboard.writeText(cmd)
                              message.success('命令已复制到剪贴板')
                            }}
                            title="复制命令"
                            onMouseEnter={(e) => {
                              e.currentTarget.style.color = 'var(--color-primary)'
                              e.currentTarget.style.backgroundColor = 'var(--color-primary-light)'
                            }}
                            onMouseLeave={(e) => {
                              e.currentTarget.style.color = 'var(--color-text-tertiary)'
                              e.currentTarget.style.backgroundColor = 'transparent'
                            }}
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </Panel>
          </Collapse>
        </Card>
      )}

      {/* Main content */}
      <div className={styles.contentCard}>
        <div className={styles.searchContainer}>
          <Input
            placeholder="搜索项目、应用和版本..."
            prefix={<SearchOutlined style={{ fontSize: '16px' }} />}
            value={searchTerm}
            onChange={handleSearchChange}
            allowClear
            size="large"
            className={styles.searchInput}
            style={{
              borderRadius: '10px',
              border: '2px solid var(--color-border-light)',
              transition: 'all 0.2s ease',
              background: 'var(--color-bg-primary)',
              color: 'var(--color-text-primary)',
            }}
            onFocus={(e) => {
              e.target.style.borderColor = 'var(--color-primary)'
              e.target.style.boxShadow = '0 0 0 3px var(--color-primary-light)'
              e.target.style.background = 'var(--color-bg-primary)'
            }}
            onBlur={(e) => {
              e.target.style.borderColor = 'var(--color-border-light)'
              e.target.style.boxShadow = 'none'
              e.target.style.background = 'var(--color-bg-primary)'
            }}
          />
        </div>

        {isLoading ? (
          <div className={styles.loadingContainer}>
            <Spin size="large" />
          </div>
        ) : treeData.length === 0 ? (
          <div className={styles.emptyContainer}>
            <Empty description={emptyText} />
          </div>
        ) : (
          <div className={styles.treeContainer}>
            <Tree
              treeData={treeData}
              expandedKeys={expandedKeys}
              onExpand={setExpandedKeys}
              showLine={{ showLeafIcon: false }}
              showIcon={false}
              blockNode
              style={{
                backgroundColor: 'transparent',
              }}
            />
          </div>
        )}
      </div>
      </div>

      {/* Footer */}
      <div className={styles.footer}>
        系统运行部驱动
      </div>
    </div>
  )
}

export default InventoryPage

