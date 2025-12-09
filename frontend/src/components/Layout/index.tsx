import { Link, Outlet } from 'react-router-dom'
import css from './index.module.scss'
import { useSelector } from '../../services/store'
import { selectIsAuth } from '../../services/slices/authSlice'
import { ButtonLink } from '../Button'
import cn from 'classnames'
import { useEffect, useRef, useState } from 'react'

export const Layout = () => {
  const [open, setOpen] = useState(false)
  const isAuth = useSelector(selectIsAuth)
  const listRef = useRef<HTMLUListElement>(null)
  const buttonRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    const handleClickOutside = (e: Event) => {
      if (
        listRef.current &&
        buttonRef.current &&
        !listRef.current.contains(e.target as Node) &&
        !buttonRef.current.contains(e.target as Node)
      ) {
        setOpen(false)
      }
    }

    document.addEventListener('click', handleClickOutside)
    return () => document.removeEventListener('click', handleClickOutside)
  }, [])

  return (
    <div className={css.layout}>
      <header className={css.header}>
        <div className={css.container}>
          <div className={css.top}>
            <div className={css.logo}>
              <Link className={css.link} to="/">
                <h1>Тарифный калькулятор</h1>
              </Link>
            </div>
            {isAuth && (
              <nav className={css.menu}>
                <button className={css.button} onClick={() => setOpen(!open)} ref={buttonRef}>
                  <span></span>
                </button>
                <ul className={cn({ [css.list]: true, [css.active]: open })} ref={listRef}>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="/">
                      Главная
                    </Link>
                  </li>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="/tariff-calculation-rk">
                      Расчет тарифа по РК
                    </Link>
                  </li>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="/distance-calculation">
                      Расчет кратчайшего расстояния
                    </Link>
                  </li>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="/history">
                      История расчетов
                    </Link>
                  </li>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="/reference/stations">
                      Справочник
                    </Link>
                  </li>
                  <li className={css.item} onClick={() => setOpen(false)}>
                    <Link className={css.link} to="#">
                      Настройка
                    </Link>
                  </li>
                </ul>
              </nav>
            )}
            <div className={css.buttonLink}>
              {!isAuth ? (
                <ButtonLink to="/login" width="90px">
                  Войти
                </ButtonLink>
              ) : (
                <ButtonLink to="/logout">Выйти</ButtonLink>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className={css.content}>
        <div className={css.container}>
          <Outlet />
        </div>
      </main>

      <footer className={css.footer}>
        <div className={css.container}></div>
      </footer>
    </div>
  )
}
