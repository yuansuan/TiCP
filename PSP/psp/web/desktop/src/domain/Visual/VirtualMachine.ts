/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable, computed } from 'mobx'
import { Http } from '@/utils'

const colorMap = new Map([
  ['online', '#52C41A'],
  ['offline', '#9B9B9B'],
  ['changing', '#F5A623'],
])

export default class VirtualMachine {
  @observable public id: number = 0
  @observable public os_name: string = ''
  @observable public name: string = ''
  @observable public path: string = ''
  @observable public status: string = ''
  @observable public agent_id: string = ''
  @observable public agent_ip: string = ''
  @observable public editing: boolean = false
  @observable public image_id: number = 0
  @observable public running_task_num: number = 0
  @observable public allocate_cpu: number = 0
  @observable public allocate_mem: number = 0
  @observable public restart_type: number = 0
  @observable public editing_task_id: number = 0
  @observable public cpu_usage: number = 0
  @observable public disk_usage: number = 0
  @observable public mem_usage: number = 0
  @observable public gpu_usage: string = '0%'
  @observable public gpu_memory: string = '0%'
  @observable public update_time: string = ''
  @observable public work_task: any = null
  @observable public children: any = []

  // 新增字段
  @observable public gpu_name: string = ''
  @observable public gpu_domain: string = ''
  @observable public gpu_bus: string = ''
  @observable public gpu_slot: string = ''
  @observable public gpu_function: string = ''
  @observable public image_os_type: string = ''
  @observable public image_path: string = ''
  @observable public resource_status: string = ''
  @observable public resource_number: number = 0
  @observable public vm_ip: string = ''

  constructor(request?: any) {
    this.init(request)
  }
  @computed
  get statusColor() {
    const s = this.status || 'offline'
    return colorMap.get(s)
  }
  @computed
  get resourceStatusColor() {
    const s = this.resource_status || 'offline'
    return colorMap.get(s)
  }
  @computed
  get allocate_mem_giga() {
    return Math.round(this.allocate_mem / Math.pow(2, 20)) + 'GB'
  }

  @action
  public edit = async () => {
    const res = await Http.post(
      `/visual/vm/${this.id}/image`,
      {},
      { baseURL: '' }
    )
    this.work_task = res.data
    this.editing_task_id = res.data.id
    this.editing = true
    this.work_task.status = 0
  }

  @action
  public close = async () => {
    await Http.put(
      `/visual/vm/${this.id}/image/${this.editing_task_id}`,
      {},
      { baseURL: '' }
    )
    this.work_task = null
    this.editing = false
  }

  @action
  public online = async () => {
    await Http.put(
      `/visual/vm/${this.id}`,
      { status: 'online' },
      { baseURL: '' }
    )
    this.status = 'changing'
  }

  @action
  public offline = async () => {
    await Http.put(
      `/visual/vm/${this.id}`,
      { status: 'offline' },
      { baseURL: '' }
    )
    this.status = 'changing'
  }

  @action
  public init = (request?: any) => {
    const {
      id,
      name,
      os_name,
      path,
      status,
      agent_id,
      agent_ip,
      editing,
      editing_task_id,
      image_id,
      running_task_num,
      allocate_cpu,
      allocate_mem,
      restart_type,
      work_task,
      children,
      cpu_usage,
      disk_usage,
      mem_usage,
      gpu_usage,
      gpu_memory,
      update_time,
      gpu_name,
      gpu_domain,
      gpu_slot,
      gpu_bus,
      gpu_function,
      image_path,
      image_os_type,
      resource_status,
      resource_number,
      vm_ip,
    } = request
    request &&
      Object.assign(this, {
        id,
        name,
        os_name,
        path,
        status,
        agent_id,
        agent_ip,
        editing,
        editing_task_id,
        image_id,
        running_task_num,
        allocate_cpu,
        allocate_mem,
        restart_type,
        work_task,
        children,
        cpu_usage,
        disk_usage,
        mem_usage,
        gpu_usage,
        gpu_memory,
        update_time,
        gpu_name,
        gpu_domain,
        gpu_slot,
        gpu_bus,
        gpu_function,
        image_path,
        image_os_type,
        resource_status,
        resource_number,
        vm_ip,
      })
  }
}
