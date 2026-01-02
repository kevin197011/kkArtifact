// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useMemo, useCallback, useEffect, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Tree, Input, Empty, Spin, Button, Typography, Card, Collapse, Tag, Space } from 'antd'
import { SearchOutlined, FolderOutlined, AppstoreOutlined, FileOutlined, LoginOutlined, DownloadOutlined, CodeOutlined, DesktopOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { publicProjectsApi, Project, App, Version } from '../api/projects'
import { downloadsApi } from '../api/downloads'
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
  const particlesRef = useRef<HTMLDivElement>(null)

  // Create subtle particle effects
  useEffect(() => {
    if (particlesRef.current) {
      const particles = particlesRef.current
      particles.innerHTML = ''

      // Create fewer particles with slower animation for subtle effect
      for (let i = 0; i < 12; i++) {
        const particle = document.createElement('div')
        particle.className = styles.particle
        particle.style.left = `${Math.random() * 100}%`
        particle.style.width = particle.style.height = `${Math.random() * 3 + 2}px`
        particle.style.animationDelay = `${Math.random() * 20}s`
        particle.style.animationDuration = `${Math.random() * 15 + 20}s`
        particles.appendChild(particle)
      }
    }
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
          console.error(`Failed to load apps for project ${project.name}:`, error)
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
            console.error(`Failed to load versions for ${projectName}/${app.name}:`, error)
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
              <span
                className={styles.treeNode}
                onClick={(e) => {
                  e.stopPropagation()
                  // Public inventory page - no navigation, just display information
                }}
              >
                <FileOutlined className={styles.treeNodeIcon} />
                <span className={styles.treeNodeTextVersion}>{version.version}</span>
                {version.created_at && (
                  <span style={{ marginLeft: '8px', fontSize: '12px', color: '#999', fontWeight: 'normal' }}>
                    ({formatDate(version.created_at)})
                  </span>
                )}
              </span>
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
              <span
                className={styles.treeNode}
                onClick={(e) => {
                  e.stopPropagation()
                  // Public inventory page - no navigation, just display information
                }}
              >
                <AppstoreOutlined className={styles.treeNodeIcon} />
                <span className={styles.treeNodeText}>{app.name}</span>
              </span>
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
            <span
              className={styles.treeNode}
              onClick={(e) => {
                e.stopPropagation()
                // Public inventory page - no navigation, just display information
              }}
            >
              <FolderOutlined className={styles.treeNodeIcon} />
              <span className={styles.treeNodeText}>{project.name}</span>
            </span>
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
  }, [allProjects, allAppsData, allVersionsData, debouncedSearchTerm, navigate])

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
  const { data: agentVersionInfo, error: agentError } = useQuery({
    queryKey: ['agent-version'],
    queryFn: () => downloadsApi.getAgentVersionInfo().then((res) => res.data),
    retry: 1,
  })
  
  // Debug: log agent version info
  useEffect(() => {
    if (agentVersionInfo) {
      console.log('Agent version info:', agentVersionInfo)
    }
    if (agentError) {
      console.error('Agent version error:', agentError)
    }
  }, [agentVersionInfo, agentError])

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const getPlatformIcon = (platform: string) => {
    if (platform.includes('windows')) return <DesktopOutlined />
    if (platform.includes('darwin')) return <DesktopOutlined />
    if (platform.includes('linux')) return <CodeOutlined />
    return <CodeOutlined />
  }

  // Format date for display
  const formatDate = (dateString: string): string => {
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
  }

  return (
    <div className={styles.inventoryContainer}>
      {/* DevOps style background effects */}
      <div className={styles.gridBackground}></div>
      <div className={styles.particles} ref={particlesRef}></div>

      {/* Content wrapper */}
      <div className={styles.contentWrapper}>
        {/* Enhanced header */}
        <div className={styles.header}>
        <div className={styles.headerContent}>
          <Title level={2} className={styles.title}>
            制品清单
          </Title>
        </div>
        <Button
          type="primary"
          icon={<LoginOutlined />}
          onClick={() => navigate('/login')}
          size="large"
          style={{
            background: 'rgba(255, 255, 255, 0.25)',
            borderColor: 'rgba(255, 255, 255, 0.4)',
            color: '#ffffff',
            fontWeight: 500,
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.background = 'rgba(255, 255, 255, 0.35)'
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = 'rgba(255, 255, 255, 0.25)'
          }}
        >
          登录后台
        </Button>
      </div>

      {/* Agent Download Section */}
      {agentVersionInfo && agentVersionInfo.binaries.length > 0 && (
        <Card className={styles.downloadCard} style={{ marginBottom: 24 }}>
          <Collapse defaultActiveKey={[]} ghost>
            <Panel 
              header={
                <Space>
                  <DownloadOutlined />
                  <span style={{ fontWeight: 500 }}>下载 Agent 客户端工具</span>
                  {agentVersionInfo.version && agentVersionInfo.version !== 'unknown' && (
                    <Tag color="blue">
                      {agentVersionInfo.version.startsWith('v') 
                        ? agentVersionInfo.version 
                        : `v${agentVersionInfo.version}`}
                    </Tag>
                  )}
                </Space>
              } 
              key="download"
            >
              <div className={styles.downloadContent}>
                <div className={styles.downloadDescription}>
                  <p>kkartifact-agent 是一个命令行工具，用于推送和拉取制品。支持并发传输、断点续传等特性。</p>
                  <div className={styles.downloadSteps}>
                    <h4>使用步骤：</h4>
                    <ol>
                      <li>下载对应平台的二进制文件</li>
                      <li>添加执行权限（Linux/macOS）：<code>chmod +x kkartifact-agent-*</code></li>
                      <li>移动到系统路径（可选）：<code>mv kkartifact-agent-* /usr/local/bin/kkartifact-agent</code></li>
                      <li>创建配置文件 <code>.kkartifact.yml</code>：</li>
                    </ol>
                    <pre className={styles.codeBlock}>
{`server_url: http://localhost:3000  # 服务器地址
token: YOUR_TOKEN_HERE             # API Token（从管理后台获取）`}
                    </pre>
                    <div className={styles.downloadUsage}>
                      <h4>常用命令：</h4>
                      <ul>
                        <li><code>kkartifact-agent push &lt;project&gt; &lt;app&gt; &lt;version&gt; [目录]</code> - 推送制品</li>
                        <li><code>kkartifact-agent pull &lt;project&gt; &lt;app&gt; &lt;version&gt; [目录]</code> - 拉取制品</li>
                      </ul>
                    </div>
                  </div>
                </div>
                <div className={styles.binariesList}>
                  <h4>可用的二进制文件：</h4>
                  <Space direction="vertical" style={{ width: '100%' }} size="middle">
                    {agentVersionInfo.binaries.map((binary) => (
                      <Card key={binary.filename} size="small" className={styles.binaryCard}>
                        <Space style={{ width: '100%', justifyContent: 'space-between' }}>
                          <Space>
                            {getPlatformIcon(binary.platform)}
                            <div>
                              <div style={{ fontWeight: 500 }}>{binary.platform}</div>
                              <div style={{ fontSize: '12px', color: '#666' }}>{binary.filename}</div>
                            </div>
                          </Space>
                          <Space>
                            <Tag>{formatFileSize(binary.size)}</Tag>
                            <Button
                              type="primary"
                              icon={<DownloadOutlined />}
                              href={downloadsApi.downloadAgent(binary.filename)}
                              download={binary.filename}
                              size="small"
                            >
                              下载
                            </Button>
                          </Space>
                        </Space>
                      </Card>
                    ))}
                  </Space>
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
            prefix={<SearchOutlined />}
            value={searchTerm}
            onChange={handleSearchChange}
            allowClear
            size="large"
            className={styles.searchInput}
            style={{
              borderRadius: '8px',
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
              showLine={false}
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
        本系统由系统部驱动
      </div>
    </div>
  )
}

export default InventoryPage

