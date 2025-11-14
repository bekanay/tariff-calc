import { useEffect } from 'react'
import { useDispatch, useSelector } from '../../services/store'
import { refreshToken } from '../../services/api/authApi'
import { logout, selectIsAuth } from '../../services/slices/authSlice'

export const useAutoRefreshToken = (intervalMs = 14 * 60 * 1000) => {
  const isAuth = useSelector(selectIsAuth)
  const dispatch = useDispatch()

  useEffect(() => {
    if (!isAuth) return

    const refresh = async () => {
      try {
        await refreshToken()
      } catch (err) {
        console.error(err)
        dispatch(logout())
      }
    }

    const interval = setInterval(refresh, intervalMs)

    return () => {
      clearInterval(interval)
    }
  }, [dispatch, intervalMs, isAuth])
}
