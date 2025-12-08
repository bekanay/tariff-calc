import { Formik, Form } from 'formik'
import * as Yup from 'yup'
import css from './index.module.scss'
import { Select } from '../Select'
import { Button } from '../Button'
import { Input } from '../Input'
import { Checkbox } from '../Checkbox'

const FormSchema = Yup.object().shape({
  shipmentType: Yup.string().required('Обьязательное поле'),
  dispatchType: Yup.string().required('Обьязательное поле'),
  speed: Yup.string().required('Обьязательное поле'),
  wagonOwnership: Yup.string().required('Обьязательное поле'),
  transportationType: Yup.string().required('Обьязательное поле'),
  originCountry: Yup.string().required('Обьязательное поле'),
  destinationCountry: Yup.string().required('Обьязательное поле'),
  originStation: Yup.string().min(2).required('Обьязательное поле'),
  destinationStation: Yup.string().min(2).required('Обьязательное поле'),
  wagonsCount: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  conductorsCount: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  wagonTare: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  axlesCount: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  wagonKind: Yup.string().required('Обьязательное поле'),
  wagonType: Yup.string().required('Обьязательное поле'),
  coordinates: Yup.string().required('Обьязательное поле'),
  distanceKm: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  carryingCapacity: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
  etsngCode: Yup.string().required('Обьязательное поле'),
  gngCode: Yup.string().required('Обьязательное поле'),
  weightTn: Yup.number().min(1, 'Минимум 1').required('Обьязательное поле'),
})

export const MyForm = () => {
  return (
    <Formik
      initialValues={{
        shipmentType: '',
        dispatchType: '',
        speed: '',
        wagonOwnership: '',
        transportationType: '',
        originCountry: '',
        originStation: '',
        destinationCountry: '',
        destinationStation: '',
        wagonsCount: 0,
        conductorsCount: 0,
        wagonTare: 0,
        axlesCount: 4,
        wagonKind: '',
        wagonType: '',
        coordinates: '',
        distanceKm: 0,
        carryingCapacity: 0,
        isEmptyWagon: false,
        hasSpecialConditions: false,
        hasDiscount: false,
        etsngCode: '',
        gngCode: '',
        weightTn: 0,
      }}
      validationSchema={FormSchema}
      onSubmit={(values) => {
        console.log('submit', values)
      }}
    >
      {({ errors, touched, setFieldValue, initialValues }) => (
        <Form className={css.form}>
          <h3 className={css.title}>1. Данные по вагону отправления</h3>
          <div className={css.formRow}>
            <Select
              options={['Внутриреспубликанская', 'Международная экспортная', 'Международная импортная', 'Транзит']}
              label="Вид отправки"
              name="shipmentType"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Повагонная', 'Групповая', 'Маршрутная', 'Контейнерная']}
              label="Тип отправки"
              name="dispatchType"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Грузовая', 'Большая']}
              label="Скорость"
              name="speed"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Собственный', 'Аренда', 'Инвентарный']}
              label="Признак вагона"
              name="wagonOwnership"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Груз ', 'Груз 2', 'Груз 3']}
              label="Перевозка"
              name="transportationType"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Казахстан', 'Китай', 'Россия']}
              label="Координаты"
              name="coordinates"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Казахстан', 'Китай', 'Россия']}
              label="Страна отправления"
              name="originCountry"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Казахстан', 'Китай', 'Россия']}
              label="Страна назначения"
              name="destinationCountry"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="search"
              label="Станция отправления"
              name="originStation"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="search"
              label="Станция назначения"
              name="destinationStation"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Количество вагонов"
              name="wagonsCount"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Количество проводников"
              name="conductorsCount"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Тара вагона, т"
              name="wagonTare"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Количество осей"
              name="axlesCount"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
              initialValues={initialValues.axlesCount}
            />
            <Select
              options={['Казахстан', 'Китай', 'Россия']}
              label="Род вагона"
              name="wagonKind"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Select
              options={['Казахстан', 'Китай', 'Россия']}
              label="Тип вагона"
              name="wagonType"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Расстояние, км"
              name="distanceKm"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Грузоподъемность, т"
              name="carryingCapacity"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
          </div>
          <div className={css.checkboxGroup}>
            <Checkbox
              id="isEmptyWagon"
              name="isEmptyWagon"
              label="Порожний вагон"
              onChange={(name, value) => setFieldValue(name, value)}
            />
            <Checkbox
              id="hasSpecialConditions"
              name="hasSpecialConditions"
              label="Особые условия перевозки"
              onChange={(name, value) => setFieldValue(name, value)}
            />
            <Checkbox
              id="hasDiscount"
              name="hasDiscount"
              label="Скидка / льгота"
              onChange={(name, value) => setFieldValue(name, value)}
            />
          </div>
          <div className={css.divider}></div>
          <h3 className={css.title}>2. Характеристики груза</h3>
          <div className={css.formRow}>
            <Input
              options={['Астана', 'Алматы', 'Караганда', 'Семей', 'Шымкент', 'Мангыстау']}
              label="Код ЕТСНГ"
              name="etsngCode"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              options={['Астана', 'Алматы', 'Караганда', 'Семей', 'Шымкент', 'Мангыстау']}
              label="Код ГНГ"
              name="gngCode"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
            <Input
              type="number"
              label="Вес, т"
              name="weightTn"
              onChange={(name, value) => setFieldValue(name, value)}
              errors={errors}
              touched={touched}
            />
          </div>
          <Button width="100%">Рассчитать</Button>
        </Form>
      )}
    </Formik>
  )
}
