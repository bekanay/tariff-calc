import { useEffect } from 'react'
import { Loader } from '../../components/Loader'
import { logout } from '../../services/slices/authSlice'
import { useDispatch } from '../../services/store'
import { useNavigate } from 'react-router-dom'

export const LogoutPage = () => {
  const dispatch = useDispatch()
  const navigate = useNavigate()

  useEffect(() => {
    dispatch(logout()).then(() => navigate('/login', { replace: true }))
  }, [dispatch, navigate])

  return <Loader type="page" />
}
