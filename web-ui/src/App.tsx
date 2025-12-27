// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import AppLayout from './components/Layout'
import ProjectsPage from './pages/Projects'
import AppsPage from './pages/Apps'
import VersionsPage from './pages/Versions'
import WebhooksPage from './pages/Webhooks'
import ConfigPage from './pages/Config'
import TokensPage from './pages/Tokens'
import AuditLogsPage from './pages/AuditLogs'
import LoginPage from './pages/Login'
import ProtectedRoute from './components/ProtectedRoute'

function App() {
  return (
    <BrowserRouter
      future={{
        v7_relativeSplatPath: true,
        v7_startTransition: true,
      }}
    >
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <AppLayout>
                <Navigate to="/projects" replace />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects"
          element={
            <ProtectedRoute>
              <AppLayout>
                <ProjectsPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:project/apps"
          element={
            <ProtectedRoute>
              <AppLayout>
                <AppsPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:project/apps/:app/versions"
          element={
            <ProtectedRoute>
              <AppLayout>
                <VersionsPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/webhooks"
          element={
            <ProtectedRoute>
              <AppLayout>
                <WebhooksPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/tokens"
          element={
            <ProtectedRoute>
              <AppLayout>
                <TokensPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/audit-logs"
          element={
            <ProtectedRoute>
              <AppLayout>
                <AuditLogsPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/config"
          element={
            <ProtectedRoute>
              <AppLayout>
                <ConfigPage />
              </AppLayout>
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
