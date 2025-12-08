export type TLoginData = {
  username: string
  password: string
}

export type TStation = {
  id: number
  stan_kod: string
  stan_name: string
  stan_priznak: number
  paragraph: string
}

export type TMetaData = {
  current_page: number
  page_size: number
  first_page: number
  last_page: number
  total_records: number
}
