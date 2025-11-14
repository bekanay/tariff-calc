export type TUserData = {
  username: string
  password: string
} 



export type TAuthState = {
  isAuth: boolean
  errorText: string | null
  isLoading: boolean
  isInit: boolean
  rememberMe: boolean
}
