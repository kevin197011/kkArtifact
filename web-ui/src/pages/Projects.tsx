// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useMemo, useCallback, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Tree, Button, message, Empty, Space, Input, Spin } from 'antd'
import { ReloadOutlined, SearchOutlined, FolderOutlined, AppstoreOutlined, FileOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { projectsApi, Project, App, Version } from '../api/projects'
import { storageApi } from '../api/storage'
import type { DataNode } from 'antd/es/tree'

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

const ProjectsPage: React.FC = () => {
  const navigate = useNavigate()
  const [searchTerm, setSearchTerm] = useState('')
  const debouncedSearchTerm = useDebounce(searchTerm, 300) // 300ms debounce
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([])
  const queryClient = useQueryClient()

  // Fetch all projects
  const { data: allProjects, isLoading: projectsLoading, error, refetch } = useQuery({
    queryKey: ['projects', 'all'],
    queryFn: () => projectsApi.list(10000, 0).then((res) => res.data),
    retry: 1,
  })

  const syncMutation = useMutation({
    mutationFn: () => storageApi.syncStorage(),
    onSuccess: (response) => {
      message.success(
        `同步完成：${response.data.projects} 个项目，${response.data.apps} 个应用，${response.data.versions} 个版本`
      )
      // Refresh projects list
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      refetch()
    },
    onError: (error: any) => {
      message.error(`同步失败：${error.response?.data?.error || error.message}`)
    },
  })

  const handleSyncStorage = () => {
    syncMutation.mutate()
  }

  // Fetch apps for a specific project (lazy loading)
  const fetchAppsForProject = useCallback((projectName: string) => {
    return queryClient.fetchQuery({
      queryKey: ['apps', projectName, 'all'],
      queryFn: () => projectsApi.getApps(projectName, 1000, 0).then((res) => res.data),
      staleTime: 5 * 60 * 1000, // Cache for 5 minutes
    })
  }, [queryClient])

  // Fetch versions for a specific app (lazy loading)
  const fetchVersionsForApp = useCallback((projectName: string, appName: string) => {
    return queryClient.fetchQuery({
      queryKey: ['versions', projectName, appName, 'all'],
      queryFn: () => projectsApi.getVersions(projectName, appName, 1000, 0).then((res) => res.data),
      staleTime: 5 * 60 * 1000, // Cache for 5 minutes
    })
  }, [queryClient])

  // Get apps for a project from cache
  const getAppsForProject = useCallback((projectName: string) => {
    return queryClient.getQueryData<App[]>(['apps', projectName, 'all'])
  }, [queryClient])

  // Get versions for an app from cache
  const getVersionsForApp = useCallback((projectName: string, appName: string) => {
    return queryClient.getQueryData<Version[]>(['versions', projectName, appName, 'all'])
  }, [queryClient])

  // Handle tree expand/collapse
  const handleExpand = useCallback((expandedKeysValue: React.Key[]) => {
    setExpandedKeys(expandedKeysValue)

    // Fetch apps for newly expanded projects (fire and forget, React Query handles caching)
    const newlyExpanded = expandedKeysValue.filter(key => !expandedKeys.includes(key))
    newlyExpanded.forEach(key => {
      // Check if it's a project key
      const project = allProjects?.find(p => `project-${p.id}` === key)
      if (project) {
        // Only fetch if not already cached
        const cachedApps = getAppsForProject(project.name)
        if (!cachedApps) {
          // Fire and forget - React Query will handle caching and state
          fetchAppsForProject(project.name).catch(error => {
            console.error(`Failed to load apps for project ${project.name}:`, error)
          })
        }
        return
      }

      // Check if it's an app key
      if (typeof key === 'string' && key.startsWith('app-')) {
        const appId = parseInt(key.replace('app-', ''))
        // Find the app and its project
        for (const proj of allProjects || []) {
          const apps = getAppsForProject(proj.name)
          if (apps) {
            const app = apps.find(a => a.id === appId)
            if (app) {
              // Only fetch if not already cached
              const cachedVersions = getVersionsForApp(proj.name, app.name)
              if (!cachedVersions) {
                // Fire and forget - React Query will handle caching and state
                fetchVersionsForApp(proj.name, app.name).catch(error => {
                  console.error(`Failed to load versions for app ${proj.name}/${app.name}:`, error)
                })
              }
              break
            }
          }
        }
      }
    })
  }, [allProjects, expandedKeys, fetchAppsForProject, fetchVersionsForApp, getAppsForProject, getVersionsForApp])

  // Check if apps are loading for a project (memoized to avoid recreating on each render)
  const isAppsLoading = useCallback((projectName: string) => {
    const queryState = queryClient.getQueryState(['apps', projectName, 'all'])
    return queryState?.status === 'pending' || queryState?.fetchStatus === 'fetching'
  }, [queryClient])

  // Check if versions are loading for an app (memoized to avoid recreating on each render)
  const isVersionsLoading = useCallback((projectName: string, appName: string) => {
    const queryState = queryClient.getQueryState(['versions', projectName, appName, 'all'])
    return queryState?.status === 'pending' || queryState?.fetchStatus === 'fetching'
  }, [queryClient])

  // Load apps for all projects when searching (to enable app filtering)
  useEffect(() => {
    if (debouncedSearchTerm.trim() && allProjects) {
      // When searching, try to load apps for all projects to enable filtering
      // React Query will cache these, so repeated searches won't cause extra requests
      allProjects.forEach(project => {
        const cachedApps = getAppsForProject(project.name)
        if (!cachedApps) {
          // Only fetch if not already cached
          fetchAppsForProject(project.name)
            .then((apps) => {
              // Also try to load versions for all apps when searching (for version filtering)
              apps?.forEach(app => {
                const cachedVersions = getVersionsForApp(project.name, app.name)
                if (!cachedVersions) {
                  fetchVersionsForApp(project.name, app.name).catch(() => {
                    // Silently fail - will show app without versions
                  })
                }
              })
            })
            .catch(() => {
              // Silently fail - will show project without apps
            })
        } else {
          // Apps are cached, check if we need to load versions
          cachedApps.forEach(app => {
            const cachedVersions = getVersionsForApp(project.name, app.name)
            if (!cachedVersions) {
              fetchVersionsForApp(project.name, app.name).catch(() => {
                // Silently fail
              })
            }
          })
        }
      })
    }
  }, [debouncedSearchTerm, allProjects, getAppsForProject, fetchAppsForProject, getVersionsForApp, fetchVersionsForApp])

  // Auto-expand projects and apps with matching items when searching
  useEffect(() => {
    if (!debouncedSearchTerm.trim() || !allProjects) return

    const searchTermLower = debouncedSearchTerm.toLowerCase().trim()

    // Find projects with matching apps or versions
    const keysToExpand = new Set<string>()
    allProjects.forEach(project => {
      const apps = getAppsForProject(project.name)
      if (apps) {
        apps.forEach(app => {
          const hasMatchingApp = app.name.toLowerCase().includes(searchTermLower)
          if (hasMatchingApp) {
            keysToExpand.add(`project-${project.id}`)
            keysToExpand.add(`app-${app.id}`)
          } else {
            // Check versions
            const versions = getVersionsForApp(project.name, app.name)
            if (versions) {
              const hasMatchingVersion = versions.some(v => 
                v.version.toLowerCase().includes(searchTermLower)
              )
              if (hasMatchingVersion) {
                keysToExpand.add(`project-${project.id}`)
                keysToExpand.add(`app-${app.id}`)
              }
            }
          }
        })
      }
    })

    // Only update if there are new keys to expand
    const newKeys = Array.from(keysToExpand)
    const currentSet = new Set(expandedKeys.map(k => String(k)))
    const hasNewKeys = newKeys.some(key => !currentSet.has(key))
    
    if (hasNewKeys) {
      setExpandedKeys(prev => {
        const prevSet = new Set(prev.map(k => String(k)))
        newKeys.forEach(key => prevSet.add(key))
        return Array.from(prevSet)
      })
    }
  }, [debouncedSearchTerm, allProjects, getAppsForProject, getVersionsForApp, expandedKeys])

  // Build tree data structure
  const treeData = useMemo(() => {
    if (!allProjects) return []

    const searchTermLower = debouncedSearchTerm.toLowerCase().trim()
    const hasSearch = searchTermLower.length > 0

    // Filter projects
    const projectsToShow = allProjects.filter(project => {
      if (!hasSearch) return true
      return project.name.toLowerCase().includes(searchTermLower)
    })

    // Get apps for projects and versions for apps
    const appsByProject: Record<string, App[]> = {}
    const versionsByApp: Record<string, Version[]> = {} // Key: "projectName/appName"
    const projectsWithMatchingApps = new Set<string>()

    if (hasSearch) {
      // When searching, check all projects for matching apps and versions
      for (const project of allProjects) {
        const apps = getAppsForProject(project.name)
        if (apps) {
          const matchingApps: App[] = []
          apps.forEach(app => {
            const appMatches = app.name.toLowerCase().includes(searchTermLower)
            const versions = getVersionsForApp(project.name, app.name)
            const matchingVersions = versions?.filter(v => 
              v.version.toLowerCase().includes(searchTermLower)
            ) || []
            
            if (appMatches || matchingVersions.length > 0) {
              matchingApps.push(app)
              if (matchingVersions.length > 0) {
                versionsByApp[`${project.name}/${app.name}`] = matchingVersions
              }
            }
          })
          
          if (matchingApps.length > 0) {
            appsByProject[project.name] = matchingApps
            projectsWithMatchingApps.add(project.name)
          }
        }
      }
    } else {
      // When not searching, only get apps for expanded projects and versions for expanded apps
      for (const key of expandedKeys) {
        const project = allProjects.find(p => `project-${p.id}` === key)
        if (project) {
          const apps = getAppsForProject(project.name)
          if (apps) {
            appsByProject[project.name] = apps
            
            // Check if any apps are expanded to load their versions
            apps.forEach(app => {
              const appKey = `app-${app.id}`
              if (expandedKeys.includes(appKey)) {
                const versions = getVersionsForApp(project.name, app.name)
                if (versions) {
                  versionsByApp[`${project.name}/${app.name}`] = versions
                }
              }
            })
          }
        }
      }
    }

    // Combine projects that match search with projects that have matching apps
    const allProjectsToShow = new Set<string>()
    projectsToShow.forEach(p => allProjectsToShow.add(p.name))
    projectsWithMatchingApps.forEach(name => allProjectsToShow.add(name))

    // Build tree nodes
    const nodes: TreeDataNode[] = Array.from(allProjectsToShow)
      .map(projectName => {
        const project = allProjects.find(p => p.name === projectName)!
        const projectKey = `project-${project.id}`
        const apps = appsByProject[project.name] || []

        const appNodes: TreeDataNode[] = apps.map(app => {
          const appKey = `app-${app.id}`
          const isAppExpanded = expandedKeys.includes(appKey)
          
          // Get versions: from versionsByApp if available (filtered for search), or from cache if app is expanded
          let versions = versionsByApp[`${project.name}/${app.name}`]
          if (!versions && isAppExpanded) {
            // Try to get all versions from cache (for non-search mode when app is expanded)
            const cachedVersions = getVersionsForApp(project.name, app.name)
            if (cachedVersions) {
              if (hasSearch) {
                // Filter cached versions for search
                versions = cachedVersions.filter(v => 
                  v.version.toLowerCase().includes(searchTermLower)
                )
              } else {
                versions = cachedVersions
              }
            }
          }
          versions = versions || []
          // Check if versions are loading (only show loading if we don't have cached data)
          const versionsQueryState = queryClient.getQueryState(['versions', project.name, app.name, 'all'])
          const hasCachedVersions = versions.length > 0 || versionsQueryState?.dataUpdatedAt
          const isVersionsLoading_ = (versionsQueryState?.status === 'pending' || versionsQueryState?.fetchStatus === 'fetching') && !hasCachedVersions

          // Build version nodes
          const versionNodes: TreeDataNode[] = versions.map(version => ({
            key: `version-${version.id}`,
            title: (
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <span 
                  style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}
                  onClick={(e) => {
                    e.stopPropagation()
                    navigate(`/projects/${project.name}/apps/${app.name}/versions`)
                  }}
                >
                  <FileOutlined />
                  <span style={{ fontFamily: 'monospace', fontSize: '12px' }}>{version.version}</span>
                </span>
                <Button
                  type="link"
                  size="small"
                  icon={<EyeOutlined />}
                  onClick={(e) => {
                    e.stopPropagation()
                    navigate(`/projects/${project.name}/apps/${app.name}/versions`)
                  }}
                >
                  View
                </Button>
              </div>
            ),
            isLeaf: true,
            version,
            isVersion: true,
            app,
            project,
          }))

          // Add loading or empty state for versions
          if (isAppExpanded) {
            if (isVersionsLoading_) {
              versionNodes.push({
                key: `loading-versions-${app.id}`,
                title: <Spin size="small" />,
                isLeaf: true,
                disabled: true,
              })
            } else if (!hasSearch && versions.length === 0 && !isVersionsLoading_) {
              versionNodes.push({
                key: `empty-versions-${app.id}`,
                title: <span style={{ color: '#999', fontStyle: 'italic' }}>No versions</span>,
                isLeaf: true,
                disabled: true,
              })
            }
          }

          return {
            key: appKey,
            title: (
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <AppstoreOutlined />
                  {app.name}
                </span>
                {!isAppExpanded && (
                  <Button
                    type="link"
                    size="small"
                    onClick={(e) => {
                      e.stopPropagation()
                      navigate(`/projects/${project.name}/apps/${app.name}/versions`)
                    }}
                  >
                    View Versions
                  </Button>
                )}
              </div>
            ),
            children: versionNodes.length > 0 ? versionNodes : undefined,
            isLeaf: versionNodes.length === 0 && !isAppExpanded,
            app,
            isApp: true,
            project,
          }
        })

        // Check if apps are loading (only show loading if we don't have cached data)
        const isExpanded = expandedKeys.includes(projectKey)
        const appsQueryState = queryClient.getQueryState(['apps', project.name, 'all'])
        // Only show loading if we're expanded, don't have cached data, and query is actually fetching
        const hasCachedData = apps.length > 0 || appsQueryState?.dataUpdatedAt
        const isLoading = (appsQueryState?.status === 'pending' || appsQueryState?.fetchStatus === 'fetching') && !hasCachedData
        
        if (isLoading && isExpanded) {
          appNodes.push({
            key: `loading-${project.id}`,
            title: <Spin size="small" />,
            isLeaf: true,
            disabled: true,
          })
        } else if (!hasSearch && isExpanded && apps.length === 0 && !isLoading) {
          appNodes.push({
            key: `empty-${project.id}`,
            title: <span style={{ color: '#999', fontStyle: 'italic' }}>No apps</span>,
            isLeaf: true,
            disabled: true,
          })
        }

        return {
          key: projectKey,
          title: (
            <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
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
      .sort((a, b) => {
        // Sort projects by name
        const nameA = a.project?.name || ''
        const nameB = b.project?.name || ''
        return nameA.localeCompare(nameB)
      })

    return nodes
  }, [allProjects, debouncedSearchTerm, expandedKeys, getAppsForProject, getVersionsForApp, isAppsLoading, isVersionsLoading, navigate])

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value)
  }, [])

  if (error) {
    message.error('Failed to load projects')
  }

  const emptyText = debouncedSearchTerm.trim()
    ? `No projects, apps, or versions match "${debouncedSearchTerm}". Try a different term.`
    : 'No projects found. Sync storage to create projects.'

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>Projects</h2>
        <Space>
          <Button
            icon={<ReloadOutlined />}
            onClick={() => refetch()}
            loading={projectsLoading}
          >
            刷新列表
          </Button>
          <Button
            type="primary"
            icon={<ReloadOutlined />}
            onClick={handleSyncStorage}
            loading={syncMutation.isPending}
          >
            同步存储
          </Button>
        </Space>
      </div>

      <div style={{ marginBottom: 16 }}>
        <Input
          placeholder="Search projects, apps, and versions..."
          prefix={<SearchOutlined />}
          value={searchTerm}
          onChange={handleSearchChange}
          allowClear
          style={{ maxWidth: 400 }}
        />
      </div>

      {projectsLoading ? (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Spin size="large" />
        </div>
      ) : treeData.length === 0 ? (
        <Empty description={emptyText} />
      ) : (
        <div style={{ backgroundColor: '#fff', padding: '16px', borderRadius: '4px' }}>
          <Tree
            treeData={treeData}
            expandedKeys={expandedKeys}
            onExpand={handleExpand}
            showLine={false}
            showIcon={false}
            blockNode
          />
        </div>
      )}
    </div>
  )
}

export default ProjectsPage
