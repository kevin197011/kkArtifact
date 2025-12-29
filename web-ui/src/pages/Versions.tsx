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
import { EyeOutlined, StarOutlined, DownloadOutlined } from '@ant-design/icons'

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

  const { data, isLoading, error } = useQuery({
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

  const promoteMutation = useMutation({
    mutationFn: (version: string) =>
      versionsApi.promote({ project: project!, app: app!, version }),
    onSuccess: () => {
      message.success('Version promoted successfully')
      queryClient.invalidateQueries({ queryKey: ['versions', project, app] })
    },
    onError: () => {
      message.error('Failed to promote version')
    },
  })

  const handleViewManifest = (version: string) => {
    setSelectedVersion(version)
    setIsManifestVisible(true)
    setManifestPage(1) // Reset to first page when opening manifest
  }

  const handlePromote = (version: string) => {
    promoteMutation.mutate(version)
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
    message.error('Failed to load versions')
  }

  const columns: ColumnsType<Version> = [
    {
      title: 'Version',
      dataIndex: 'version',
      key: 'version',
      width: 300,
      render: (version: string) => (
        <Text style={{ fontFamily: 'monospace', fontSize: '12px' }} copyable>
          {version}
        </Text>
      ),
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      }),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Tooltip title="View Manifest">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewManifest(record.version)}
            >
              Manifest
            </Button>
          </Tooltip>
          <Tooltip title="Promote Version">
            <Popconfirm
              title="Are you sure to promote this version?"
              onConfirm={() => handlePromote(record.version)}
            >
              <Button type="link" size="small" icon={<StarOutlined />}>
                Promote
              </Button>
            </Popconfirm>
          </Tooltip>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={() => navigate('/projects')}>Projects</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>
          <a onClick={() => navigate(`/projects/${project}/apps`)}>{project}</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>{app}</Breadcrumb.Item>
      </Breadcrumb>
      <h2>
        Versions - {project}/{app}
      </h2>
      <Table
        columns={columns}
        dataSource={data}
        loading={isLoading}
        rowKey="id"
        pagination={{
          current: page,
          pageSize,
          total: data?.length || 0,
          onChange: setPage,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 个版本`,
        }}
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={
          <span>
            Manifest Details - {selectedVersion && (
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
          <div>Loading...</div>
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
              <Descriptions.Item label="Files Count">
                <Tag color="blue">{manifest.files?.length || 0}</Tag>
              </Descriptions.Item>
            </Descriptions>
            <div style={{ marginTop: 16 }}>
              <h4>Files:</h4>
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
                    title: 'File Path',
                    dataIndex: 'path',
                    key: 'path',
                    width: 300,
                    ellipsis: true,
                  },
                  {
                    title: 'Size',
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
                    title: 'SHA256 Hash',
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
                    title: 'Actions',
                    key: 'actions',
                    render: (_, record) => (
                      <Button
                        type="link"
                        size="small"
                        icon={<DownloadOutlined />}
                        onClick={() => handleDownloadFile(selectedVersion!, record.path)}
                      >
                        Download
                      </Button>
                    ),
                  },
                ]}
              />
            </div>
          </div>
        ) : (
          <div>No manifest data</div>
        )}
      </Modal>
    </div>
  )
}

export default VersionsPage
