// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { Button } from 'antd'
import { SunOutlined, MoonOutlined } from '@ant-design/icons'
import { useTheme, ThemeScope } from '../hooks/useTheme'

interface ThemeToggleProps {
  scope?: ThemeScope
  size?: 'small' | 'middle' | 'large'
  type?: 'text' | 'default' | 'primary'
}

const ThemeToggle: React.FC<ThemeToggleProps> = ({ 
  scope = 'frontend',
  size = 'middle',
  type = 'text'
}) => {
  const { theme, toggleTheme } = useTheme(scope)

  return (
    <Button
      type={type}
      icon={theme === 'light' ? <MoonOutlined /> : <SunOutlined />}
      onClick={toggleTheme}
      size={size}
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '6px',
        color: 'var(--color-text-secondary)',
      }}
    >
      {theme === 'light' ? '暗色' : '亮色'}
    </Button>
  )
}

export default ThemeToggle
