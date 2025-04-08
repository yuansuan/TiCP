import { observable } from 'mobx'

interface SysUserProps {
  name: string //系统用户名
  count: number //会话连接数
 
}
export default class SysUser {
  @observable name: string
  @observable count: number

  constructor(props?: SysUserProps) {
    Object.assign(this, props)
  }
}
