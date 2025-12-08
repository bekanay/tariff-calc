import { Button } from '../../components/Button'
import { useEffect, useMemo, useState, type ChangeEvent } from 'react'
import {
  addStation,
  clearStation,
  deleteStation,
  editStation,
  fetchStations,
  fetchWithKodStation,
  selectMetaData,
  selectStation,
  selectStations,
} from '../../services/slices/stationsSlice'
import { useDispatch, useSelector } from '../../services/store'
import cn from 'classnames'
import css from './index.module.scss'
import clear from '../../assets/images/close.svg'
import { Pagination } from '@mui/material'
import { debounce } from 'lodash'
import type { Sort } from '../../services/api/stationsApi'
import { Modal } from '../../components/Modal'
import { StationForm } from '../../components/StationForm'
import type { TStation } from '../../utils/types'
import { useSearchParams } from 'react-router-dom'

export const StationsPage = () => {
  const dispatch = useDispatch()

  const stations = useSelector(selectStations)
  const station = useSelector(selectStation)
  const metaData = useSelector(selectMetaData)

  const [isOpenAdd, setIsOpenAdd] = useState(false)
  const [isOpenEdit, setIsOpenEdit] = useState(false)
  const [isOpenDelete, setIsOpenDelete] = useState(false)
  const [selectedStation, setSelectedStation] = useState<TStation | null>(null)

  const [searchParams, setSearchParams] = useSearchParams()
  const initialSort = (searchParams.get('sort') ?? 'id') as Sort
  const initialSearch = searchParams.get('search') ?? ''
  const initialPage = Number(searchParams.get('page') ?? 1)
  const [sort, setSort] = useState<Sort>(initialSort)
  const [searchValue, setSearchValue] = useState(initialSearch)
  const [value, setValue] = useState(initialSearch)
  const [currentPage, setCurrentPage] = useState(initialPage)

  const updateURL = (params: { sort?: string; search?: string; page?: number }) => {
    const newParams = new URLSearchParams(searchParams)

    if (params.sort !== undefined) newParams.set('sort', params.sort)
    if (params.search !== undefined) newParams.set('search', params.search)
    if (params.page !== undefined) newParams.set('page', String(params.page))

    setSearchParams(newParams)
  }

  const debouncedSearch = useMemo(
    () =>
      debounce((text: string) => {
        setSearchValue(text)
        setCurrentPage(1)
      }, 400),
    []
  )

  useEffect(() => {
    if (!searchValue) {
      dispatch(fetchStations({ current_page: currentPage, sort }))
      return
    }

    if (/^\d{6}$/.test(searchValue)) {
      dispatch(fetchWithKodStation(searchValue))
      return
    }

    dispatch(
      fetchStations({
        current_page: currentPage,
        name: searchValue,
        sort,
      })
    )
  }, [dispatch, currentPage, searchValue, sort])

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value)

    setCurrentPage(1)
    updateURL({ search: e.target.value, page: 1 })

    if (station) dispatch(clearStation())
    debouncedSearch(e.target.value)
  }

  const onClick = () => {
    setValue('')
    setSearchValue('')
    updateURL({ search: '', page: 1 })
    dispatch(clearStation())
  }

  const toggleSort = (field: 'id' | 'stan_kod' | 'stan_name') => {
    setSort((prev) => {
      if (prev === field) {
        updateURL({ sort: `-${field}` })
        return `-${field}`
      }
      updateURL({ sort: field })
      return field
    })
  }

  const dataToRender = station ? [station] : stations

  return (
    <div className={css.wrapper}>
      <h2 className={css.title}>Справочник станции</h2>
      <div className={css.container}>
        <div className={css.top}>
          <div className={css.inputWrapper}>
            <input
              className={css.input}
              value={value}
              onChange={(e) => onChange(e)}
              placeholder="Поиск по коду или названию станции..."
            />
            {value && <img src={clear} onClick={onClick} />}
          </div>
          <Button onClick={() => setIsOpenAdd(true)}>&#x271A; Добавить станцию</Button>
        </div>
        <table>
          <thead>
            <tr>
              <th className={css.toggle} onClick={() => toggleSort('id')}>
                ID {sort === 'id' ? '▲' : sort === '-id' ? '▼' : ''}
              </th>
              <th className={css.toggle} onClick={() => toggleSort('stan_kod')}>
                Код станции {sort === 'stan_kod' ? '▲' : sort === '-stan_kod' ? '▼' : ''}
              </th>
              <th className={css.toggle} onClick={() => toggleSort('stan_name')}>
                Название {sort === 'stan_name' ? '▲' : sort === '-stan_name' ? '▼' : ''}
              </th>
              <th>Признак</th>
              <th>Параграф</th>
              <th>Дейсвия</th>
            </tr>
          </thead>
          <tbody>
            {dataToRender && dataToRender.length > 0 ? (
              dataToRender.map((s) => (
                <tr key={s.id}>
                  <td>{s.id}</td>
                  <td>{s.stan_kod}</td>
                  <td>{s.stan_name}</td>
                  <td>{s.stan_priznak}</td>
                  <td>{s.paragraph}</td>
                  <td>
                    <div className={css.actions}>
                      <button
                        className={css.button}
                        onClick={() => {
                          setSelectedStation(s)
                          setIsOpenEdit(true)
                        }}
                      >
                        Изменить
                      </button>
                      <button
                        className={cn(css.button, css.remove)}
                        onClick={() => {
                          setIsOpenDelete(true)
                          setSelectedStation(s)
                        }}
                      >
                        Удалить
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={6} style={{ textAlign: 'center', padding: 20 }}>
                  Нет данных для отображения
                </td>
              </tr>
            )}
          </tbody>
        </table>

        {metaData && (
          <>
            <div style={{ marginTop: '20px', fontSize: '16px' }}>Всего записей: {metaData.total_records}</div>
            <Pagination
              defaultPage={currentPage}
              count={metaData.last_page}
              color="primary"
              size="large"
              onChange={(_, page) => {
                setCurrentPage(page)
                updateURL({ page })
              }}
              sx={{ marginTop: 2 }}
            />
          </>
        )}

        <Modal title="Добавить станцию" isOpen={isOpenAdd} onClose={() => setIsOpenAdd(false)}>
          <StationForm
            onSubmit={(data) => {
              dispatch(addStation(data)).then(() => {
                setIsOpenAdd(false)
                dispatch(fetchStations({ current_page: currentPage, sort }))
              })
            }}
            submitText="Добавить"
          />
        </Modal>

        <Modal
          title="Изменить станцию"
          isOpen={isOpenEdit}
          onClose={() => {
            setIsOpenEdit(false)
            setSelectedStation(null)
          }}
        >
          {selectedStation && (
            <StationForm
              initialValues={{
                stan_kod: selectedStation.stan_kod,
                stan_name: selectedStation.stan_name,
                stan_priznak: selectedStation.stan_priznak,
                paragraph: selectedStation.paragraph,
              }}
              onSubmit={(data) => {
                dispatch(
                  editStation({
                    data,
                    kod: selectedStation.stan_kod,
                  })
                ).then(() => {
                  setIsOpenEdit(false)
                  setSelectedStation(null)
                  dispatch(fetchStations({ current_page: currentPage, sort }))
                })
              }}
              submitText="Сохранить"
            />
          )}
        </Modal>

        <Modal
          title="Удаление"
          isOpen={isOpenDelete}
          onClose={() => {
            setIsOpenDelete(false)
            setSelectedStation(null)
          }}
        >
          <div>
            <p>Удалить станцию {selectedStation?.stan_name}?</p>
            <div className={css.buttons}>
              <button
                className={css.button}
                onClick={() => {
                  setIsOpenDelete(false)
                  setSelectedStation(null)
                }}
              >
                Нет
              </button>
              <button
                className={`${css.button} ${css.yes}`}
                onClick={() => {
                  if (selectedStation) {
                    dispatch(deleteStation(selectedStation.stan_kod)).then(() => {
                      setIsOpenDelete(false)
                      setSelectedStation(null)
                      dispatch(fetchStations({ current_page: currentPage, sort }))
                    })
                  }
                }}
              >
                Да
              </button>
            </div>
          </div>
        </Modal>
      </div>
    </div>
  )
}
