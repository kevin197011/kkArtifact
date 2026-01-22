// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table,
  Breadcrumb,
  Tag,
  message,
  Button,
  Space,
  Modal,
  Descriptions,
  Popconfirm,
  Tooltip,
  Typography,
} from 'antd'
import { useParams, useNavigate } from 'react-router-dom'
import { projectsApi, Version } from '../api/projects'
import { versionsApi } from '../api/versions'
import type { ColumnsType } from 'antd/es/table'
import { EyeOutlined, StarOutlined, StarFilled, DownloadOutlined, DeleteOutlined } from '@ant-design/icons'

const { Text } = Typography

const VersionsPage: React.FC = () => {
  const { project, app } = useParams<{ project: string; app: string }>()
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [selectedVersion, setSelectedVersion] = useState<string | null>(null)
  const [isManifestVisible, setIsManifestVisible] = useState(false)
  const [manifestPageSize, setManifestPageSize] = useState(10) // Default 10 items per page
  const [manifestPage, setManifestPage] = useState(1) // Current page for manifest files
  const pageSize = 50
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['versions', project, app, page],
    queryFn: () =>
      projectsApi.getVersions(project!, app!, pageSize, (page - 1) * pageSize).then((res) => res.data),
    enabled: !!project && !!app,
  })

  const { data: manifest, isLoading: isManifestLoading } = useQuery({
    queryKey: ['manifest', project, app, selectedVersion],
    queryFn: () => versionsApi.getManifest(project!, app!, selectedVersion!).then((res) => res.data),
    enabled: !!selectedVersion && !!project && !!app,
  })

  const publishMutation = useMutation({
    mutationFn: (version: string) =>
      versionsApi.publish({ project: project!, app: app!, version }),
    onSuccess: () => {
      message.success('版本发布成功')
      queryClient.invalidateQueries({ queryKey: ['versions', project, app] })
    },
    onError: () => {
      message.error('发布版本失败')
    },
  })

  const unpublishMutation = useMutation({
    mutationFn: (version: string) =>
      versionsApi.unpublish({ project: project!, app: app!, version }),
    onSuccess: () => {
      message.success('版本已取消发布')
      queryClient.invalidateQueries({ queryKey: ['versions', project, app] })
    },
    onError: () => {
      message.error('取消发布失败')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (version: string) => projectsApi.deleteVersion(project!, app!, version),
    onSuccess: () => {
      message.success('版本删除成功')
      queryClient.invalidateQueries({ queryKey: ['versions', project, app] })
      refetch()
    },
    onError: (error: any) => {
      message.error(`删除版本失败：${error.response?.data?.error || error.message}`)
    },
  })

  const handleViewManifest = (version: string) => {
    setSelectedVersion(version)
    setIsManifestVisible(true)
    setManifestPage(1) // Reset to first page when opening manifest
  }

  const handlePublish = (version: string) => {
    publishMutation.mutate(version)
  }

  const handleUnpublish = (version: string) => {
    unpublishMutation.mutate(version)
  }

  const handleTogglePublish = (version: string, isPublished: boolean) => {
    if (isPublished) {
      handleUnpublish(version)
    } else {
      handlePublish(version)
    }
  }

  const handleDownloadFile = (version: string, filePath: string) => {
    versionsApi
      .downloadFile(project!, app!, version, filePath)
      .then((response) => {
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', filePath.split('/').pop() || 'file')
        document.body.appendChild(link)
        link.click()
        link.remove()
        message.success('File download started')
      })
      .catch(() => {
        message.error('Failed to download file')
      })
  }

  if (error) {
    message.error('加载版本失败')
  }

  const columns: ColumnsType<Version> = [
    {
      title: <span style={{ fontWeight: 600 }}>版本</span>,
      dataIndex: 'version',
      key: 'version',
      width: 300,
      render: (version: string) => (
        <Text style={{ fontFamily: 'monospace', fontSize: '13px', color: '#1a1a1a' }} copyable>
          {version}
        </Text>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>创建时间</span>,
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => (
        <span style={{ color: '#8c8c8c', fontSize: '14px' }}>
          {new Date(date).toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
          })}
        </span>
      ),
    },
    {
      title: <span style={{ fontWeight: 600 }}>操作</span>,
      key: 'actions',
      width: 280,
      render: (_, record) => (
        <Space>
          <Tooltip title="查看清单">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewManifest(record.version)}
              style={{
                padding: '0 8px',
                height: '32px',
                fontWeight: 500,
              }}
            >
              清单
            </Button>
          </Tooltip>
          <Tooltip title={record.is_published ? "点击取消发布" : "发布版本"}>
            <Popconfirm
              title={record.is_published ? "确定要取消发布此版本吗？" : "确定要发布此版本吗？"}
              description={record.is_published ? "取消发布后，将无法通过 pull latest 获取此版本" : "发布此版本后，其他已发布版本将被取消发布"}
              onConfirm={() => handleTogglePublish(record.version, record.is_published)}
            >
              <Button 
                type="link" 
                size="small" 
                icon={record.is_published ? <StarFilled style={{ color: '#faad14' }} /> : <StarOutlined />}
                style={{
                  padding: '0 8px',
                  height: '32px',
                  fontWeight: 500,
                }}
              >
                {record.is_published ? '已发布' : '发布'}
              </Button>
            </Popconfirm>
          </Tooltip>
          <Tooltip title="删除版本">
            <Popconfirm
              title="确定要删除此版本吗？"
              description="删除版本将永久删除该版本的所有文件，此操作不可恢复！"
              onConfirm={() => deleteMutation.mutate(record.version)}
              okText="确定"
              cancelText="取消"
              okButtonProps={{ danger: true }}
            >
              <Button 
                type="link" 
                size="small" 
                danger 
                icon={<DeleteOutlined />}
                style={{
                  padding: '0 8px',
                  height: '32px',
                }}
              >
                删除
              </Button>
            </Popconfirm>
          </Tooltip>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Breadcrumb 
        style={{ marginBottom: '24px', fontSize: '14px' }}
        items={[
          {
            title: <a onClick={() => navigate('/projects')} style={{ color: '#1890ff' }}>项目</a>,
          },
          {
            title: <a onClick={() => navigate(`/projects/${project}/apps`)} style={{ color: '#1890ff' }}>{project}</a>,
          },
          {
            title: <span style={{ color: '#8c8c8c' }}>{app}</span>,
          },
        ]}
      />
      <div style={{ marginBottom: '24px' }}>
        <h2 style={{ margin: 0, fontSize: '24px', fontWeight: 600, color: 'var(--color-text-primary)', letterSpacing: '-0.3px' }}>
          版本 - {project}/{app}
        </h2>
        <div style={{ marginTop: '6px', color: 'var(--color-text-secondary)', fontSize: '13px' }}>
          管理应用的所有版本
        </div>
      </div>
      <div
        style={{
          background: 'var(--color-bg-primary)',
          borderRadius: 'var(--radius-md)',
          border: '1px solid var(--color-border-light)',
          overflow: 'hidden',
        }}
      >
        <Table
          columns={columns}
          dataSource={data}
          loading={isLoading}
          rowKey="id"
          pagination={{
            current: page,
            pageSize,
            total: data && data.length < pageSize 
              ? (page - 1) * pageSize + data.length 
              : data 
              ? page * pageSize + 1 
              : 0,
            onChange: setPage,
            showSizeChanger: true,
            showTotal: (total, range) => {
              if (data && data.length < pageSize) {
                return `共 ${total} 个版本`
              }
              return `第 ${range[0]}-${range[1]} 项，至少 ${total} 个版本`
            },
            style: {
              padding: '16px 24px',
            },
          }}
          scroll={{ x: 'max-content' }}
        />
      </div>
      <Modal
        title={
          <span>
            清单详情 - {selectedVersion && (
              <Text style={{ fontFamily: 'monospace' }} copyable>{selectedVersion}</Text>
            )}
          </span>
        }
        open={isManifestVisible}
        onCancel={() => {
          setIsManifestVisible(false)
          setSelectedVersion(null)
        }}
        footer={null}
        width={1000}
      >
        {isManifestLoading ? (
          <div>加载中...</div>
        ) : manifest ? (
          <div>
            <Descriptions bordered column={1} size="small">
              <Descriptions.Item label="Project">
                <Tag>{manifest.project}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="App">
                <Tag>{manifest.app}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Version">
                <Text style={{ fontFamily: 'monospace' }} copyable>
                  {manifest.version}
                </Text>
              </Descriptions.Item>
              {manifest.git_commit && (
                <Descriptions.Item label="Git Commit">
                  <Text style={{ fontFamily: 'monospace' }} copyable>
                    {manifest.git_commit}
                  </Text>
                </Descriptions.Item>
              )}
              <Descriptions.Item label="Builder">
                <Tag>{manifest.builder}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Build Time">
                {new Date(manifest.build_time).toLocaleString('zh-CN', {
                  year: 'numeric',
                  month: '2-digit',
                  day: '2-digit',
                  hour: '2-digit',
                  minute: '2-digit',
                  second: '2-digit',
                })}
              </Descriptions.Item>
              <Descriptions.Item label="文件数量">
                <Tag color="blue">{manifest.files?.length || 0}</Tag>
              </Descriptions.Item>
            </Descriptions>
            <div style={{ marginTop: 16 }}>
              <h4>文件列表：</h4>
              <Table
                dataSource={manifest.files}
                rowKey="path"
                size="small"
                pagination={
                  manifest.files && manifest.files.length > 0
                    ? {
                        current: manifestPage,
                        pageSize: manifestPageSize,
                        total: manifest.files.length,
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 个文件`,
                        pageSizeOptions: ['10', '20', '50', '100'],
                        showQuickJumper: true,
                        simple: false,
                        onChange: (page) => {
                          setManifestPage(page)
                        },
                        onShowSizeChange: (_current, size) => {
                          setManifestPageSize(size)
                          setManifestPage(1) // Reset to first page when page size changes
                        },
                      }
                    : false
                }
                scroll={{ x: 'max-content' }}
                columns={[
                  {
                    title: '文件路径',
                    dataIndex: 'path',
                    key: 'path',
                    width: 300,
                    ellipsis: true,
                  },
                  {
                    title: '大小',
                    dataIndex: 'size',
                    key: 'size',
                    width: 120,
                    render: (size: number) => {
                      if (size < 1024) {
                        return `${size} B`
                      } else if (size < 1024 * 1024) {
                        return `${(size / 1024).toFixed(2)} KB`
                      } else {
                        return `${(size / (1024 * 1024)).toFixed(2)} MB`
                      }
                    },
                  },
                  {
                    title: 'SHA256 哈希',
                    dataIndex: 'hash',
                    key: 'hash',
                    width: 300,
                    render: (hash: string) => (
                      <Text style={{ fontFamily: 'monospace', fontSize: '11px' }} copyable>
                        {hash}
                      </Text>
                    ),
                  },
                  {
                    title: '操作',
                    key: 'actions',
                    render: (_, record) => (
                      <Button
                        type="link"
                        size="small"
                        icon={<DownloadOutlined />}
                        onClick={() => handleDownloadFile(selectedVersion!, record.path)}
                      >
                        下载
                      </Button>
                    ),
                  },
                ]}
              />
            </div>
          </div>
        ) : (
          <div>无清单数据</div>
        )}
      </Modal>
    </div>
  )
}

export default VersionsPage
