import { setCookie } from '../../utils/cookie'
import type { TLoginData } from '../../utils/types'

const URL = import.meta.env.VITE_URL

const checkResponse = <T>(res: Response): Promise<T> =>
  res.ok ? res.json() : res.json().then((err) => Promise.reject(err))

type TServerResponse<T> = {
  success: boolean
} & T

type TRefreshResponse = TServerResponse<{
  refresh_token: string
  access_token: string
  expires_in: number
}>

type TAuthResponse = TServerResponse<{
  refresh_token: string
  access_token: string
  expires_in: number
}>

export const refreshToken = async (): Promise<TRefreshResponse> => {
  const storage = localStorage.getItem('refreshToken') ? localStorage : sessionStorage

  const token = storage.getItem('refreshToken')

  if (!token) {
    throw new Error('No refresh token found')
  }

  const res = await fetch(`${URL}/refresh-token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
    },
    body: JSON.stringify({
      refresh_token: token,
    }),
  })

  const refreshData = await checkResponse<TRefreshResponse>(res)

  if (!refreshData?.access_token || !refreshData?.refresh_token) {
    throw new Error('Invalid refresh response')
  }

  storage.setItem('refreshToken', refreshData.refresh_token)
  setCookie('accessToken', refreshData.access_token, {
    SameSite: 'Strict',
    'max-age': refreshData.expires_in,
  })

  return refreshData
}

export const fetchWithRefresh = async <T>(url: RequestInfo, options: RequestInit) => {
  try {
    const res = await fetch(url, options)
    return await checkResponse<T>(res)
  } catch (err) {
    if ((err as { message: string }).message === 'jwt expired') {
      const refreshData = await refreshToken()
      if (options.headers) {
        ;(options.headers as { [key: string]: string }).authorization = refreshData.access_token
      }
      const res = await fetch(url, options)
      return await checkResponse<T>(res)
    } else {
      return Promise.reject(err)
    }
  }
}

export const loginApi = (data: TLoginData) =>
  fetch(`${URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
    },
    body: JSON.stringify(data),
    credentials: 'include',
  })
    .then((res) => checkResponse<TAuthResponse>(res))
    .then((data) => {
      if (data?.access_token && data?.refresh_token) return data
      return Promise.reject(data)
    })

export const logoutApi = () =>
  fetch(`${URL}/logout`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
    },
    body: JSON.stringify({
      refresh_token: localStorage.getItem('refreshToken') || sessionStorage.getItem('refreshToken'),
    }),
  }).then((res) => checkResponse<TServerResponse<Record<string, never>>>(res))
