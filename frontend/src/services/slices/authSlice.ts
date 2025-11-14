import { createAsyncThunk, createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { TAuthState, TUserData } from '../../utils/types'
import { loginApi, logoutApi } from '../api/authApi'
import { deleteCookie, getCookie, setCookie } from '../../utils/cookie'

export const login = createAsyncThunk(
  'auth/login',
  async ({ username, password, rememberMe }: TUserData & { rememberMe: boolean }) => {
    const response = await loginApi({ username, password })
    return { ...response, rememberMe }
  }
)

export const logout = createAsyncThunk('auth/logout', async () => logoutApi())

export const initialState: TAuthState = {
  isAuth: false,
  errorText: null,
  isLoading: false,
  isInit: false,
  rememberMe: false
}

export const AuthSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    initializeAuth(state) {
      const accessToken = getCookie('accessToken')
      const refreshToken = localStorage.getItem('refreshToken') || sessionStorage.getItem('refreshToken')
      state.isAuth = !!(accessToken && refreshToken)
      state.isInit = true
    },
    setRememberMe(state, action: PayloadAction<boolean>) {
      state.rememberMe = action.payload
    }
  },
  extraReducers(builder) {
    builder
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
          Secure: true,
          SameSite: 'Strict',
          'max-age': 900,
        })
      })
      .addCase(login.rejected, (state, action) => {
        state.isLoading = false
        state.errorText = action.error.message || 'Неверный логин или пароль'
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

export const selectIsAuth = (state: { auth: TAuthState }) => state.auth.isAuth
export const selectGetError = (state: { auth: TAuthState }) => state.auth.errorText
export const selectIsLoading = (state: { auth: TAuthState }) => state.auth.isLoading
export const selectIsInit = (state: { auth: TAuthState }) => state.auth.isInit
export const selectRememberMe = (state: { auth: TAuthState }) => state.auth.rememberMe

export const { initializeAuth } = AuthSlice.actions
