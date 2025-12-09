import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import type { TMetaData, TStation } from '../../utils/types'
import {
  addStationApi,
  deleteStationApi,
  editStationApi,
  getStationsApi,
  getStationWithKodApi,
  type TStationsQuery,
  type TStationsResponse,
} from '../api/stationsApi'
import type { TStationForm } from '../../components/StationForm'

interface StationsInitState {
  stations: TStation[]
  station: TStation | null
  metaData: TMetaData | null
  isLoading: boolean
  errorText: string | null
}

export const fetchStations = createAsyncThunk<TStationsResponse, TStationsQuery>('stations/fetch', async (params) =>
  getStationsApi(params)
)

export const fetchWithKodStation = createAsyncThunk('station/fetchWithCodeStation', async (kod: string) =>
  getStationWithKodApi(kod)
)

export const addStation = createAsyncThunk('station/addStation', async (data: TStationForm) => addStationApi(data))

export const editStation = createAsyncThunk(
  'station/editStation',
  async ({ data, kod }: { data: TStationForm; kod: string }) => editStationApi(data, kod)
)

export const deleteStation = createAsyncThunk('station/deleteStation', async (kod: string) => deleteStationApi(kod))

const initialState: StationsInitState = {
  stations: [],
  station: null,
  metaData: null,
  isLoading: false,
  errorText: null,
}

export const stationsSlice = createSlice({
  name: 'station',
  initialState,
  reducers: {
    clearStation(state) {
      state.station = null
    },
  },

  extraReducers(builder) {
    builder
      .addCase(fetchStations.pending, (state) => {
        state.isLoading = true
      })
      .addCase(fetchStations.fulfilled, (state, action) => {
        state.stations = action.payload.stations
        state.metaData = action.payload.metadata
        state.isLoading = false
      })
      .addCase(fetchStations.rejected, (state, action) => {
        state.errorText = action.error.message || 'Ошибка'
        state.isLoading = false
      })

    builder
      .addCase(fetchWithKodStation.pending, (state) => {
        state.isLoading = true
      })
      .addCase(fetchWithKodStation.fulfilled, (state, action) => {
        state.station = action.payload
        state.isLoading = false
      })
      .addCase(fetchWithKodStation.rejected, (state, action) => {
        state.errorText = action.error.message || 'Ошибка'
        state.isLoading = false
      })

    builder
      .addCase(addStation.pending, (state) => {
        state.isLoading = true
      })
      .addCase(addStation.fulfilled, (state) => {
        state.isLoading = false
        state.errorText = null
      })
      .addCase(addStation.rejected, (state, action) => {
        state.isLoading = false
        state.errorText = action.error.message || 'Не удалось добавить'
      })

    builder
      .addCase(editStation.pending, (state) => {
        state.isLoading = true
      })
      .addCase(editStation.fulfilled, (state) => {
        state.isLoading = false
      })
      .addCase(editStation.rejected, (state, action) => {
        state.isLoading = false
        state.errorText = action.error.message || 'Не удалось изменить'
      })
  },
})

export const selectStations = (state: { stations: StationsInitState }) => state.stations.stations
export const selectMetaData = (state: { stations: StationsInitState }) => state.stations.metaData
export const selectStation = (state: { stations: StationsInitState }) => state.stations.station
export const selectIsLoading = (state: { stations: StationsInitState }) => state.stations.isLoading
export const selectErrorText = (state: { stations: StationsInitState }) => state.stations.errorText

export const { clearStation } = stationsSlice.actions
