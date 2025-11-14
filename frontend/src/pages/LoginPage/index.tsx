import css from './index.module.scss'
import eyeHide from '../../assets/images/eyeHide.svg'
import eyeShow from '../../assets/images/eyeShow.svg'
import checkmark from '../../assets/images/checkmark.svg'
import { useEffect, useState, type SyntheticEvent } from 'react'
import { useDispatch, useSelector } from '../../services/store'
import { selectGetError, selectIsLoading, login } from '../../services/slices/authSlice'
import { Button } from '../../components/Button'

export const LoginPage = () => {
  const dispatch = useDispatch()
  const error = useSelector(selectGetError)
  const isLoading = useSelector(selectIsLoading)

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [rememberMe, setRememberMe] = useState(false)

  useEffect(() => {
    const savedUsername = localStorage.getItem('login')
    if (savedUsername) {
      setUsername(savedUsername)
    }
  }, [])

  const onSubmit = (e: SyntheticEvent) => {
    e.preventDefault()
     if (rememberMe) {
      localStorage.setItem('login', username)
    } else {
      localStorage.removeItem('login')
    }
    dispatch(login({ username, password, rememberMe }))
  }

  return (
    <div className={css.container}>
      <div className={css.header}>
        <img src={checkmark} alt='' />
        <h2 className={css.title}>Вход в систему</h2>
        <p className={css.text}>Введите свои учетные данные</p>
      </div>

      <form className={css.form} name="auth" onSubmit={onSubmit}>
        <div className={css.inputWrapper}>
          <label className={css.label} htmlFor="username">
            Логин
          </label>
          <input
            className={css.input}
            type="text"
            id="username"
            name="username"
            placeholder="Введите логин"
            onChange={(e) => setUsername(e.target.value)}
            value={username}
            required
          />
        </div>
        <div className={css.inputWrapper}>
          <label className={css.label} htmlFor="password">
            Пароль
          </label>
          <input
            className={`${css.input}`}
            type={showPassword ? 'text' : 'password'}
            id="password"
            name="password"
            placeholder="Введите пароль"
            onChange={(e) => setPassword(e.target.value)}
            value={password}
            required
          />
          <div className={css.eye} onClick={() => setShowPassword((prew) => !prew)}>
            {showPassword ? <img src={eyeHide} /> : <img src={eyeShow} />}
          </div>
        </div>
        <div className={css.checkboxWrapper}>
          <input
            type="checkbox"
            id="rememberMe"
            checked={rememberMe}
            onChange={(e) => setRememberMe(e.target.checked)}
            disabled={isLoading}
          />
          <label htmlFor="rememberMe">
            Запомнить меня
          </label>
        </div>
        <Button loading={isLoading} width='100%'>Войти</Button>
        {error && <div className={css.error}>{error}</div>}
      </form>
    </div>
  )
}
