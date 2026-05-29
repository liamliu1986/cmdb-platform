import client from './client'

export interface LoginReq {
  username: string
  password: string
}

export interface LoginRes {
  code: number
  data: {
    token: string
    user_id: number
    username: string
  }
  message: string
}

export const authApi = {
  login: (data: LoginReq) => client.post<LoginRes>('/auth/login', data),
  register: (data: { username: string; password: string; email: string; nickname?: string }) =>
    client.post('/auth/register', data),
}
