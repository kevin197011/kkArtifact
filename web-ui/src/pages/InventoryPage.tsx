// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useMemo, useCallback, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Tree, Input, Empty, Spin, Button, Typography } from 'antd'
import { SearchOutlined, FolderOutlined, AppstoreOutlined, FileOutlined, LoginOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { publicProjectsApi, Project, App, Version } from '../api/projects'
import type { DataNode } from 'antd/es/tree'

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
                style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}
                onClick={(e) => {
                  e.stopPropagation()
                  const token = localStorage.getItem('kkartifact_token')
                  if (token) {
                    navigate(`/projects/${project.name}/apps/${app.name}/versions`)
                  } else {
                    navigate(`/login?redirect=/projects/${project.name}/apps/${app.name}/versions`)
                  }
                }}
              >
                <FileOutlined />
                <span style={{ fontFamily: 'monospace', fontSize: '12px' }}>{version.version}</span>
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
                style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}
                onClick={(e) => {
                  e.stopPropagation()
                  const token = localStorage.getItem('kkartifact_token')
                  if (token) {
                    navigate(`/projects/${project.name}/apps/${app.name}/versions`)
                  } else {
                    navigate(`/login?redirect=/projects/${project.name}/apps/${app.name}/versions`)
                  }
                }}
              >
                <AppstoreOutlined />
                {app.name}
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
              style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}
              onClick={(e) => {
                e.stopPropagation()
                const token = localStorage.getItem('kkartifact_token')
                if (token) {
                  navigate(`/projects/${project.name}/apps`)
                } else {
                  navigate(`/login?redirect=/projects/${project.name}/apps`)
                }
              }}
            >
              <FolderOutlined />
              {project.name}
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
    ? `No projects, apps, or versions match "${debouncedSearchTerm}". Try a different term.`
    : 'No projects found.'

  const hasToken = !!localStorage.getItem('kkartifact_token')

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f0f2f5', padding: '24px' }}>
      {/* Simple header */}
      <div
        style={{
          backgroundColor: '#fff',
          padding: '16px 24px',
          marginBottom: '24px',
          borderRadius: '8px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <Title level={3} style={{ margin: 0 }}>
            Artifact Inventory
          </Title>
        </div>
        {!hasToken && (
          <Button type="primary" icon={<LoginOutlined />} onClick={() => navigate('/login')}>
            Login
          </Button>
        )}
      </div>

      {/* Main content */}
      <div style={{ backgroundColor: '#fff', padding: '24px', borderRadius: '8px' }}>
        <div style={{ marginBottom: 16 }}>
          <Input
            placeholder="Search projects, apps, and versions..."
            prefix={<SearchOutlined />}
            value={searchTerm}
            onChange={handleSearchChange}
            allowClear
            size="large"
            style={{ maxWidth: 400 }}
          />
        </div>

        {isLoading ? (
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <Spin size="large" />
          </div>
        ) : treeData.length === 0 ? (
          <Empty description={emptyText} />
        ) : (
          <Tree
            treeData={treeData}
            expandedKeys={expandedKeys}
            onExpand={setExpandedKeys}
            showLine={false}
            showIcon={false}
            blockNode
          />
        )}
      </div>
    </div>
  )
}

export default InventoryPage

