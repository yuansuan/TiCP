import { observable } from 'mobx'

interface InfoProps {
  name: string 
  version: string 
  expiry: string 
  available_days: number
}
export default class Info {
  @observable name: string
  @observable version: string
  @observable expiry: string
  @observable available_days: number

  constructor(props?: InfoProps) {
    Object.assign(this, props)
  }
}