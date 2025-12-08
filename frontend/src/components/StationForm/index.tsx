import { Field, Form, Formik } from 'formik'
import * as Yup from 'yup'
import css from './index.module.scss'

export type TStationForm = {
  stan_kod: string
  stan_name: string
  stan_priznak: number | string
  paragraph: string
}

type StationFormProps = {
  initialValues?: TStationForm
  onSubmit: (values: TStationForm) => void
  submitText?: string
}

const StationSchema = Yup.object().shape({
  stan_kod: Yup.string().min(6, 'Необходимо 6 символов').max(6, 'Необходимо 6 символов').required('Обязательное поле'),
  stan_name: Yup.string().required('Обязательное поле'),
  stan_priznak: Yup.number().typeError('Введите число').min(0, 'Минимум 0').required('Обязательное поле'),
  paragraph: Yup.string(),
})

export const StationForm = ({
  initialValues = {
    stan_kod: '',
    stan_name: '',
    stan_priznak: '',
    paragraph: '',
  },
  onSubmit,
  submitText = 'Сохранить',
}: StationFormProps) => {
  return (
    <Formik initialValues={initialValues} validationSchema={StationSchema} enableReinitialize onSubmit={onSubmit}>
      {({ errors, touched }) => (
        <Form className={css.form}>
          <div className={css.inputWrapper}>
            <label className={css.label} htmlFor="stan_kod">
              Код станции
            </label>
            <Field className={css.input} id="stan_kod" name="stan_kod" placeholder="Введите" />
            {errors.stan_kod && touched.stan_kod && <div className={css.error}>{errors.stan_kod}</div>}
          </div>

          <div className={css.inputWrapper}>
            <label className={css.label} htmlFor="stan_name">
              Название станции
            </label>
            <Field className={css.input} id="stan_name" name="stan_name" placeholder="Введите" />
            {errors.stan_name && touched.stan_name && <div className={css.error}>{errors.stan_name}</div>}
          </div>

          <div className={css.inputWrapper}>
            <label className={css.label} htmlFor="stan_priznak">
              Признак
            </label>
            <Field
              className={css.input}
              type="number"
              id="stan_priznak"
              name="stan_priznak"
              min="0"
              placeholder="Введите"
            />
            {errors.stan_priznak && touched.stan_priznak && <div className={css.error}>{errors.stan_priznak}</div>}
          </div>

          <div className={css.inputWrapper}>
            <label className={css.label} htmlFor="paragraph">
              Параграф
            </label>
            <Field className={css.input} id="paragraph" name="paragraph" placeholder="Введите" />
            {errors.paragraph && touched.paragraph && <div className={css.error}>{errors.paragraph}</div>}
          </div>

          <button className={css.button} type="submit">
            {submitText}
          </button>
        </Form>
      )}
    </Formik>
  )
}
