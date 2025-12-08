import { useState } from 'react'
import css from './index.module.scss'

type TCheckbox = {
  id: string
  name: string
  label: string
  onChange: (name: string, value: boolean) => void
}

export const Checkbox = ({ id, name, label, onChange }: TCheckbox) => {
  const [checked, setChecked] = useState(false)

  const handleSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    setChecked(e.target.checked)
    onChange(name, e.target.checked)
  }

  return (
    <div className={css.checkboxWrapper}>
      <input type="checkbox" id={id} checked={checked} name={name} onChange={(e) => handleSelect(e)} />
      <label htmlFor={id} className={css.label}>
        {label}
      </label>
    </div>
  )
}
