import { observable } from 'mobx'
import { formatDateFromMilliSec } from '@/utils/formatter'

interface ListProps {
  id: string
  ip_address: string
  username: string
  create_time: number
}

export default class WhiteList {
  @observable id: string
  @observable ip: string
  @observable username: string
  @observable time: number

  constructor(props?: ListProps) {
    Object.assign(this, {
      id: props.id,
      ip: props.ip_address,
      username: props.username,
      time: formatDateFromMilliSec(props.create_time)
    })
  }
}
