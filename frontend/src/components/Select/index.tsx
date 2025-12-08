import { useEffect, useRef, useState } from 'react'
import cn from 'classnames'
import css from './index.module.scss'

export type TSelect = {
  label: string
  options: string[]
  name: string
  onChange: (name: string, value: string) => void
  width?: string
  errors?: Record<string, string>
  touched?: Record<string, boolean>
}

export const Select = ({ label, options, name, onChange, errors, touched }: TSelect) => {
  const [open, setOpen] = useState(false)
  const [value, setValue] = useState('')

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

  const handleSelect = (option: string) => {
    setValue(option)
    onChange(name, option)
    setOpen(false)
  }

  const hasError = touched?.[name] && errors?.[name]

  return (
    <div className={css.wrapper} ref={wrapperRef}>
      <p className={css.label}>{label}</p>
      <div
        className={cn({
          [css.select]: true,
          [css.placeholder]: !value,
          [css.errors]: hasError,
        })}
        onClick={() => setOpen(!open)}
      >
        {value || 'Выберите'}
      </div>

      {hasError && <div className={css.error}>{errors[name]}</div>}

      {open && (
        <ul className={css.dropdown}>
          {options.map((option) => (
            <li key={option} className={css.item} onClick={() => handleSelect(option)}>
              {option}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
