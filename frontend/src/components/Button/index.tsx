import cn from 'classnames'
import css from './index.module.scss'
import { Link } from 'react-router-dom'

export type ButtonProps = {
  type?: 'button' | 'submit' | 'reset'
  loading?: boolean
  children: React.ReactNode
  color?: 'red' | 'green' | 'blue'
  onClick?: React.MouseEventHandler
  disabled?: boolean
  className?: string
  width?: string
}

export type ButtonLinkProps = {
  to: string
  color?: 'red' | 'green'
  handler?: React.MouseEventHandler
  width?: string
  children: React.ReactNode
}

export const Button = ({
  type = 'submit',
  loading = false,
  color = 'blue',
  onClick,
  disabled,
  width,
  className,
  children,
}: ButtonProps) => {
  return (
    <button
      type={type}
      className={cn({
        [css.button]: true,
        ...(className && { [css[className]]: true }),
        [css.disabled]: loading || disabled,
        [css.loading]: loading,
        [css[`color-${color}`]]: color,
      })}
      onClick={onClick}
      disabled={loading || disabled}
      style={{ width: width }}
    >
      <span className={css.text}>{children}</span>
    </button>
  )
}

export const ButtonLink = ({ to, color = 'green', handler, width, children }: ButtonLinkProps) => {
  return (
    <Link
      className={cn({ [css.link]: true, [css[`color-${color}`]]: color })}
      onClick={handler}
      to={to}
      style={{ width: width }}
    >
      {children}
    </Link>
  )
}
