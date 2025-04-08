import { observable } from 'mobx'

interface UserSessionProps {
  jti:string //JWT Token
  expire_time:string // 过期时间
  ip: string //IP地址
}
export default class UserSession {
  @observable jti: string
  @observable expire_time: string
  @observable ip: string

  constructor(props?: UserSessionProps) {
    Object.assign(this, props)
  }
}
