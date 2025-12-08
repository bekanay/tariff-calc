import { ButtonLink } from '../../components/Button'
import css from './index.module.scss'

export const MainPage = () => {
  return (
    <div className={css.wrapper}>
      <div className={css.hero}>
        <div className={css.content}>
          <h2 className={css.title}>Рассчитайте оптимальный тариф на железнодорожную доставку</h2>
          <p className={css.subtitle}>
            Укажите станции отправления и прибытия — калькулятор подберет лучший маршрут и стоимость
          </p>
          <ButtonLink to="/calculator" width="260px">
            Перейти к калькулятору
          </ButtonLink>
        </div>
        {/* <div className={css.img_block}>
          <img src={routeIllustration} />
        </div> */}
      </div>
      {/* <div className={css.cards}>
        <div className={css.card}>
          <img src={route} />
          <p>Оптимальные маршруты</p>
          <p className={css.subtext}>Автоматический подбор наиболее выгодного пути</p>
        </div>
        <div className={css.card}>
          <img src={pricing} />
          <p>Прозрачные тарифы</p>
          <p className={css.subtext}>Детальная разбивка всех составляющих цены</p>
        </div>
        <div className={css.card}>
          <img src={accuracy} />
          <p>Точная логика расчетов</p>
          <p className={css.subtext}>Учет всех коэффициентов и дополнительных сборов</p>
        </div>
      </div> */}
    </div>
  )
}
