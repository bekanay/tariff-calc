import css from './index.module.scss'

export const HistoryPage = () => {
  return (
    <div className={css.wrapper}>
      <h2 className={css.title}>История расчетов</h2>
      <div className={css.tableCard}>
        <div className={css.container}>
          <table>
            <thead>
              <tr>
                <th>Дата</th>
                <th>Маршрут</th>
                <th>Груз</th>
                <th>Стоимость</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>20.11.2025г</td>
                <td>Алматы → Астана</td>
                <td>Уголь каменный</td>
                <td>2 526 156 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
              <tr>
                <td>19.11.2025</td>
                <td>Караганда → Шымкент</td>
                <td>Пшеница</td>
                <td>3 786 256 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
              <tr>
                <td>20.11.2025</td>
                <td>Алматы → Астана</td>
                <td>Уголь каменный</td>
                <td>2 526 156 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
              <tr>
                <td>19.11.2025</td>
                <td>Караганда → Шымкент</td>
                <td>Пшеница</td>
                <td>3 786 256 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
              <tr>
                <td>20.11.2025</td>
                <td>Алматы → Астана</td>
                <td>Уголь каменный</td>
                <td>2 526 156 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
              <tr>
                <td>19.11.2025</td>
                <td>Караганда → Шымкент</td>
                <td>Пшеница</td>
                <td>3 786 256 ₸</td>
                <td>
                  <button className={css.button}>Открыть</button>
                  <button className={css.button}>↻</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
