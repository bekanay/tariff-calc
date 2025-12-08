import { combineReducers, configureStore } from '@reduxjs/toolkit'
import type { TypedUseSelectorHook } from 'react-redux'
import { useDispatch as dispatchHook, useSelector as selectorHook } from 'react-redux'
import { authSlice } from './slices/authSlice'
import { stationsSlice } from './slices/stationsSlice'

const rootReducer = combineReducers({
  auth: authSlice.reducer,
  stations: stationsSlice.reducer,
})

const store = configureStore({
  reducer: rootReducer,
})

export type RootState = ReturnType<typeof rootReducer>

export type AppDispatch = typeof store.dispatch

export const useDispatch: () => AppDispatch = () => dispatchHook()
export const useSelector: TypedUseSelectorHook<RootState> = selectorHook

export default store
