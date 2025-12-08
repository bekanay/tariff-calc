import { memo, useEffect, type FC, type ReactNode } from 'react'
import ReactDOM from 'react-dom'
import close from '../../assets/images/close.svg'
import css from './index.module.scss'

export type TModalProps = {
  title: string
  isOpen: boolean
  onClose: () => void
  children?: ReactNode
}

export const Modal: FC<TModalProps> = memo(({ title, isOpen, onClose, children }) => {
  const modalRoot = document.getElementById('modals')

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    }

    return () => {
      document.body.style.overflow = ''
    }
  }, [isOpen])

  if (!isOpen) return null

  return ReactDOM.createPortal(
    <div className={css.backdrop}>
      <div className={css.window}>
        <div className={css.header}>
          <h3 className={css.title}>{title}</h3>
          <button className={css.button} type="button" onClick={onClose}>
            <img className={css.close} src={close} />
          </button>
        </div>
        {children}
      </div>
    </div>,

    modalRoot as HTMLDivElement
  )
})
