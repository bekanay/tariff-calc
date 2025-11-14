import { Navigate, useLocation } from 'react-router-dom'
import { useSelector } from '../../services/store'
import { selectIsAuth, selectIsInit } from '../../services/slices/authSlice'
import { Loader } from '../Loader'

type ProtectedRouteProps = {
  children: React.ReactElement
  unAuthOnly?: boolean
}

export const ProtectedRoute = ({ children, unAuthOnly }: ProtectedRouteProps) => {
  const isAuth = useSelector(selectIsAuth)
  const isInit = useSelector(selectIsInit)
  const location = useLocation()

  if (!isInit) {
    return (
      <Loader type='page'/>
    )
  }

  if (!unAuthOnly && !isAuth) {
    const safeFrom = location.pathname === '/logout' || location.pathname === '/login' ? '/' : location.pathname

    return <Navigate to="/login" state={{ from: safeFrom }} />
  }

  if (unAuthOnly && isAuth) {
    const from = (location.state as { from?: string })?.from || '/'

    const forbidden = ['/logout', '/login']
    const redirectTo = forbidden.includes(from) ? '/' : from

    return <Navigate to={redirectTo} replace />
  }

  return children
}

// export const ProtectedRoute = ({
//   children,
//   unAuthOnly
// }: ProtectedRouteProps) => {
//   const isAuth = useSelector(selectisAuth);
//   const location = useLocation();

//   if (!unAuthOnly && !isAuth) {
//     return <Navigate replace to='/login' state={{ from: location }} />;
//   }

//   if (unAuthOnly && isAuth) {
//     const from = location.state?.from || { pathname: '/' };
//     return <Navigate replace to={from} />;
//   }

//   return children;
// };
