import { MyForm } from '../../components/Form'
import { ResultCulc } from '../../components/ResultCulc'

import css from './index.module.scss'

export const CalculatorPage = () => {
  return (
    <div className={css.wrapper}>
      <h2 className={css.title}>Тарифный калькулятор</h2>
      <p className={css.subtitle}>Введите данные для расчета стоимости перевозки</p>
      <div className={css.card}>
        <MyForm />
        <ResultCulc />
      </div>
    </div>
  )
}
