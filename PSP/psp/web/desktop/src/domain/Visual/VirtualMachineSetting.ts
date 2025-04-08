/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, runInAction } from 'mobx'
import { Http } from '@/utils'

type userTask = {
  enable: boolean
  number: number // 用户任务数
}
export default class VirtualMachineSetting {
  @observable public user_task: userTask
  @observable public default_vm_task_number: number = 0 //虚拟机任务数
  constructor(request?: any) {
    if (request) {
      this.default_vm_task_number = request.default_vm_task_number
      this.user_task = request.user_task
    }
  }

  async fetch() {
    const res = await Http.get('/visual/settings/list', { baseURL: '' })

    runInAction(() => {
      Object.assign(this, res.data)
    })
  }
  async update(data) {
    await Http.post('/visual/settings/update', data, {
      baseURL: '',
    })
  }
}
