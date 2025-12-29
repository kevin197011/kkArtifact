// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useMemo, useCallback, useEffect, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Tree, Input, Empty, Spin, Button, Typography } from 'antd'
import { SearchOutlined, FolderOutlined, AppstoreOutlined, FileOutlined, LoginOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { publicProjectsApi, Project, App, Version } from '../api/projects'
import type { DataNode } from 'antd/es/tree'
import styles from './InventoryPage.module.css'

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

