import { useEffect } from 'react'
import { useDispatch } from '../../services/store'
import { refreshToken } from '../../services/api/authApi'
import { logout } from '../../services/slices/authSlice'

export const useAutoRefreshToken = (intervalMs = 4 * 60 * 1000) => {
  const dispatch = useDispatch()

  useEffect(() => {
    const refresh = async () => {
      try {
        await refreshToken()
      } catch (err) {
        console.error(err)
        dispatch(logout())
      }
    }

    window.addEventListener('focus', refresh)
    const interval = setInterval(refresh, intervalMs)

    return () => {
      window.removeEventListener('focus', refresh)
      clearInterval(interval)
    }
  }, [dispatch, intervalMs])
}
