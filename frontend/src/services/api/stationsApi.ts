import type { TStationForm } from '../../components/StationForm'
import { getCookie } from '../../utils/cookie'
import type { TMetaData, TStation } from '../../utils/types'
import { fetchWithRefresh } from './authApi'

const URL = import.meta.env.VITE_URL

type TServerResponse<T> = {
  success: boolean
} & T

export type TStationsResponse = TServerResponse<{
  stations: TStation[]
  metadata: TMetaData
}>

export type TStationResponse = TServerResponse<{
  station: TStation
}>

export type Sort = 'id' | '-id' | 'stan_kod' | '-stan_kod' | 'stan_name' | '-stan_name'

export type TStationsQuery = {
  current_page?: number
  page_size?: number
  sort?: Sort
  name?: string
}

export const getStationsApi = (params: TStationsQuery = {}) => {
  const searchParams = new URLSearchParams()
  if (params.current_page) searchParams.append('page', params.current_page.toString())
  if (params.page_size) searchParams.append('page_size', params.page_size.toString())
  if (params.sort) searchParams.append('sort', params.sort)
  if (params.name) searchParams.append('name', params.name)

  const queryString = searchParams.toString()
  const url = `${URL}/stations${queryString ? `?${queryString}` : ''}`

  return fetchWithRefresh<TStationsResponse>(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
      authorization: getCookie('accessToken'),
    } as HeadersInit,
  }).then((data) => {
    if (data) return data
    return Promise.reject(data)
  })
}

export const getStationWithKodApi = (kod: string) =>
  fetchWithRefresh<TStation>(`${URL}/stations/${kod}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
      authorization: getCookie('accessToken'),
    } as HeadersInit,
  }).then((data) => {
    console.log(data)
    if (data) return data
    return Promise.reject(data)
  })

export const addStationApi = (data: TStationForm) =>
  fetchWithRefresh<TStationResponse>(`${URL}/stations`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
      authorization: getCookie('accessToken'),
    } as HeadersInit,
    body: JSON.stringify(data),
  }).then((data) => {
    if (data) return data
    return Promise.reject(data)
  })

export const editStationApi = (data: TStationForm, kod: string) =>
  fetchWithRefresh<TStationForm>(`${URL}/stations/${kod}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
      authorization: getCookie('accessToken'),
    } as HeadersInit,
    body: JSON.stringify(data),
  })

export const deleteStationApi = (kod: string) =>
  fetchWithRefresh<void>(`${URL}/stations/${kod}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
      authorization: getCookie('accessToken'),
    } as HeadersInit,
  })
