import { Button } from '../Button'
import css from './index.module.scss'

export const ResultCulc = () => {
  return (
    <div>
      <h3 className={css.title}>3. Результат расчета</h3>
      <div className={css.card}>
        <div className={`${css.total} ${css.row}`}>
          <div>
            <p className={css.accent}>Итоговая стоимость перевозки</p>
            <p className={css.accent}>1 247 850 ₸</p>
          </div>
          <div>
            <p>Стоимость на 1 вагон</p>
            <p>247 850 ₸</p>
          </div>
          <div>
            <p>Стоимость на 1 тонну</p>
            <p>24 785 ₸</p>
          </div>
        </div>
        <div className={css.divider}></div>
        <div className={`${css.breakdown} ${css.row}`}>
          <div>
            <p>Базовый тариф</p>
            <p>45 785 ₸</p>
          </div>
          <div>
            <p>Коэффициенты</p>
            <p>85 785 ₸</p>
          </div>
          <div>
            <p>Сборы</p>
            <p>35 785 ₸</p>
          </div>
          <div>
            <p>НДС (12%)</p>
            <p>19 785 ₸</p>
          </div>
        </div>
        <div className={css.buttons}>
          <Button className="resultCulcBtn">📄 Скачать PDF</Button>
          <Button className="resultCulcBtn">✉️ Отправить расчёт</Button>
          <Button className="resultCulcBtn">💾 Сохранить в историю</Button>
        </div>
      </div>
    </div>
  )
}
