import { observable } from 'mobx'

interface INodeProps {
  id: string
  node_name: string
  scheduler_status: string
  status: string
  queue_name: string
  node_type: string
  total_core_num: number
  used_core_num: number
  free_core_num: number
  total_mem: number
  used_mem: number
  free_mem: number
  available_mem: number
  create_time: string
}

export default class Node {
  @observable id: string
  @observable node_name: string
  @observable scheduler_status: string
  @observable status: string
  @observable queue_name: string
  @observable node_type: string
  @observable total_core_num: number
  @observable used_core_num: number
  @observable free_core_num: number
  @observable total_mem: number
  @observable used_mem: number
  @observable free_mem: number
  @observable available_mem: number
  @observable create_time: string

  constructor(props?: INodeProps) {
    Object.assign(this, props)
  }
}
