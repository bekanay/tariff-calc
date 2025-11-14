import cn from 'classnames'
import css from './index.module.scss'
import { Link } from 'react-router-dom'

export type ButtonProps = {
  type?: 'button' | 'submit' | 'reset'
  loading?: boolean
  children: React.ReactNode
  color?: 'red' | 'green' | 'blue'
  handler?: React.MouseEventHandler<HTMLButtonElement>
  disabled?: boolean
  width?: string
}

export type ButtonLinkProps = {
  to: string
  color?: 'red' | 'green'
  handler?: React.MouseEventHandler
  children: React.ReactNode
}

export const Button = ({ type = 'submit', loading = false, color = 'blue', handler, disabled, width, children }: ButtonProps) => {
  return (
    <button
      type={type}
      className={cn({
        [css.button]: true,
        [css.disabled]: loading || disabled,
        [css.loading]: loading,
        [css[`color-${color}`]]: color,
      })}
      onClick={handler}
      disabled={loading || disabled}
      style={{width: width}}
    >
      <span className={css.text}>{children}</span>
    </button>
  )
}

export const ButtonLink = ({ to, color = 'green', handler, children }: ButtonLinkProps) => {
  return (
    <Link className={cn({ [css.link]: true, [css[`color-${color}`]]: color })} onClick={handler} to={to}>
      {children}
    </Link>
  )
}
