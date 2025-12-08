import { useEffect, useRef, useState } from 'react'
import cn from 'classnames'
import css from './index.module.scss'

export type TInput = {
  type?: string
  label: string
  options?: string[]
  name: string
  onChange: (name: string, value: string | number) => void
  width?: string
  placeholder?: string
  min?: number
  max?: number
  step?: number
  errors?: Record<string, string>
  touched?: Record<string, boolean>
  initialValues?: string | number
}

export const Input = ({
  type = 'text',
  name,
  label,
  options,
  onChange,
  width,
  placeholder = 'Начните вводить название...',
  min = 0,
  max,
  step = 1,
  errors,
  touched,
  initialValues,
}: TInput) => {
  const [open, setOpen] = useState(false)
  const [value, setValue] = useState<string | number>(type === 'number' && initialValues ? initialValues : '')
  const [filtered, setFiltered] = useState(options)

  const wrapperRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleInputChange = (searchValue: string) => {
    if (type === 'number') {
      const numValue = searchValue === '' ? 0 : Number(searchValue)
      setValue(numValue)
      onChange(name, numValue)
    } else {
      setValue(searchValue)

      const filteredOptions = options?.filter((option) => option.toLowerCase().includes(searchValue.toLowerCase()))
      setFiltered(filteredOptions)
      setOpen(true)
      onChange(name, searchValue)
    }
  }

  const handleSelect = (option: string) => {
    setValue(option)
    onChange(name, option)
    setOpen(false)
  }

  const hasError = touched?.[name] && errors?.[name]

  return (
    <div className={css.wrapper} ref={wrapperRef} style={{ width }}>
      <p className={css.label}>{label}</p>
      <input
        type={type}
        className={cn({
          [css.input]: true,
          [css.errors]: hasError,
        })}
        value={value}
        onChange={(e) => handleInputChange(e.target.value)}
        onFocus={() => type !== 'number' && setOpen(true)}
        placeholder={type === 'number' ? '0' : placeholder}
        min={type === 'number' ? min : undefined}
        max={type === 'number' ? max : undefined}
        step={type === 'number' ? step : undefined}
      />
      {hasError && <div className={css.error}>{errors[name]}</div>}
      {type !== 'number' && open && (
        <ul className={`${css.dropdown} ${open ? css.open : ''}`}>
          {filtered?.map((option) => (
            <li key={option} className={css.item} onClick={() => handleSelect(option)}>
              {option}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
