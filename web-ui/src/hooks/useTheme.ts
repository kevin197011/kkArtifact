// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'

export type Theme = 'light' | 'dark'
export type ThemeScope = 'frontend' | 'backend'

const THEME_STORAGE_KEY = {
  frontend: 'kkartifact_theme_frontend',
  backend: 'kkartifact_theme_backend',
}

export function useTheme(scope: ThemeScope = 'frontend') {
  const storageKey = THEME_STORAGE_KEY[scope]
  
  // Initialize theme from localStorage or default to light
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window === 'undefined') return 'light'
    const stored = localStorage.getItem(storageKey)
    return (stored === 'dark' || stored === 'light') ? stored : 'light'
  })

  // Apply theme to document on mount and when theme changes
  useEffect(() => {
    const root = document.documentElement
    if (theme === 'dark') {
      root.classList.add('theme-dark')
      root.classList.remove('theme-light')
    } else {
      root.classList.add('theme-light')
      root.classList.remove('theme-dark')
    }
    if (typeof window !== 'undefined') {
      localStorage.setItem(storageKey, theme)
    }
  }, [theme, storageKey])

  // Sync state with localStorage on mount (in case theme was set by inline script)
  useEffect(() => {
    if (typeof window === 'undefined') return
    const stored = localStorage.getItem(storageKey)
    const storedTheme = (stored === 'dark' || stored === 'light') ? stored : 'light'
    if (storedTheme !== theme) {
      setTheme(storedTheme)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []) // Only run once on mount

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'))
  }

  return { theme, setTheme, toggleTheme }
}
