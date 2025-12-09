import { MyForm } from '../../components/Form'
import { ResultCulc } from '../../components/ResultCulc'

import css from './index.module.scss'

export const TariffCalculationPage = () => {
  return (
    <div className={css.wrapper}>
      <h2 className={css.title}>Расчет тарифа по РК</h2>
      <p className={css.subtitle}>Введите данные для расчета стоимости перевозки</p>
      <div className={css.container}>
        <MyForm />
        <ResultCulc />
      </div>
    </div>
  )
}
