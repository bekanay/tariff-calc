import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/Layout'
import './styles/global.scss'
import { MainPage } from './pages/MainPage'
import { ProtectedRoute } from './components/Protected-route'
import { useEffect } from 'react'
import { useDispatch } from './services/store'
import { initializeAuth } from './services/slices/authSlice'
import { LoginPage } from './pages/LoginPage'
import { LogoutPage } from './pages/LogoutPage'
import { useAutoRefreshToken } from './utils/hooks/useAutoRefreshToken'

function App() {
  const dispatch = useDispatch()

  useEffect(() => {
    dispatch(initializeAuth())
  }, [dispatch])

  useAutoRefreshToken()

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/logout"
          element={
            <ProtectedRoute>
              <LogoutPage />
            </ProtectedRoute>
          }
        />
        <Route element={<Layout />}>
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/login"
            element={
              <ProtectedRoute unAuthOnly>
                <LoginPage />
              </ProtectedRoute>
            }
          />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
