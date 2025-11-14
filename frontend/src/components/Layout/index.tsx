import { Link, Outlet } from 'react-router-dom'
import css from './index.module.scss'
import logo from '../../assets/images/Logo.png'
import { useSelector } from '../../services/store'
import { selectIsAuth } from '../../services/slices/authSlice'
import { ButtonLink } from '../Button'

export const Layout = () => {
  const isAuth = useSelector(selectIsAuth)

  return (
    <div className={css.layout}>
      <div className={css.navigation}>
        <div className={css.logo}>
          <img src={logo} alt="Logo" width={100} height={100} />
          <h1>Тарифный калькулятор</h1>
        </div>
        <ul className={css.menu}>
          <li className={css.item}>
            <Link className={css.link} to="/">
              Главная
            </Link>
          </li>
          <li className={css.item}>
            {!isAuth ? <ButtonLink to="/login">Войти</ButtonLink> : <ButtonLink to="/logout">Выйти</ButtonLink>}
          </li>
        </ul>
      </div>
      <div className={css.content}>
        <Outlet />
      </div>
      <div className={css.footer}>© Филиал АО "НК "ҚТЖ" - "ЦДАЦ" 2025г.</div>
    </div>
  )
}
