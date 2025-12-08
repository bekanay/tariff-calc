import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import type { TLoginData } from '../../utils/types'
import { loginApi, logoutApi, refreshToken } from '../api/authApi'
import { deleteCookie, getCookie, setCookie } from '../../utils/cookie'

interface AuthInitState {
  isAuth: boolean
  errorText: string | null
  isLoading: boolean
  isInit: boolean
  rememberMe: boolean
}

export const initializeAuth = createAsyncThunk('auth/initializeAuth', async () => {
  const access = getCookie('accessToken')
  const refresh = localStorage.getItem('refreshToken') || sessionStorage.getItem('refreshToken')

  if (!refresh) return false

  if (access) return true

  try {
    await refreshToken()
    return true
  } catch {
    return false
  }
})

export const login = createAsyncThunk(
  'auth/login',
  async ({ username, password, rememberMe }: TLoginData & { rememberMe: boolean }) => {
    const response = await loginApi({ username, password })
    return { ...response, rememberMe }
  }
)

export const logout = createAsyncThunk('auth/logout', async () => logoutApi())

export const initialState: AuthInitState = {
  isAuth: false,
  errorText: null,
  isLoading: false,
  isInit: false,
  rememberMe: false,
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {},
  extraReducers(builder) {
    builder
      .addCase(initializeAuth.pending, (state) => {
        state.isInit = false
      })
      .addCase(initializeAuth.fulfilled, (state, action) => {
        state.isAuth = action.payload
        state.isInit = true
      })
      .addCase(initializeAuth.rejected, (state) => {
        state.isAuth = false
        state.isInit = true
      })
      .addCase(login.pending, (state) => {
        state.isLoading = true
      })
      .addCase(login.fulfilled, (state, action) => {
        state.isLoading = false
        state.isAuth = true
        state.errorText = null
        if (action.payload.rememberMe) {
          localStorage.setItem('refreshToken', action.payload.refresh_token)
        } else {
          sessionStorage.setItem('refreshToken', action.payload.refresh_token)
        }
        setCookie('accessToken', action.payload.access_token, {
          SameSite: 'Strict',
          'max-age': action.payload.expires_in,
        })
      })
      .addCase(login.rejected, (state, action) => {
        state.isLoading = false
        state.errorText = action.error.message || 'Не удалось войти'
      })
      .addCase(logout.pending, (state) => {
        state.isLoading = true
      })
      .addCase(logout.fulfilled, (state) => {
        deleteCookie('accessToken')
        localStorage.removeItem('refreshToken')
        sessionStorage.removeItem('refreshToken')
        state.isLoading = false
        state.isAuth = false
      })
      .addCase(logout.rejected, (state, action) => {
        state.isLoading = false
        state.errorText = action.error.message || 'Не удалось выйти'
      })
  },
})

export const selectIsAuth = (state: { auth: AuthInitState }) => state.auth.isAuth
export const selectGetError = (state: { auth: AuthInitState }) => state.auth.errorText
export const selectIsLoading = (state: { auth: AuthInitState }) => state.auth.isLoading
export const selectIsInit = (state: { auth: AuthInitState }) => state.auth.isInit
export const selectRememberMe = (state: { auth: AuthInitState }) => state.auth.rememberMe
