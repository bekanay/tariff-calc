import { Button } from '../../components/Button'
import css from './index.module.scss'

export const MainPage = () => {
  return (
    <div className={css.container}>
      <h2 className={css.title}>Добро пожаловать</h2>
      <p className={css.text}>Выберите тариф и рассчитайте стоимость</p>
      <Button type="button">Начать</Button>
    </div>
  )
}
