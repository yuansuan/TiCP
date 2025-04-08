/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, runInAction } from 'mobx'

import { Http } from '@/utils'
import VirtualMachine from './VirtualMachine'

export default class VirtualMachineList {
  @observable childrenMap = {}
  @observable list = new Map()

  get = (id, isRoot) => {
    if (isRoot) {
      return this.list.get(id)
    } else {
      return this.childrenMap[id]
    }
  }

  fetch = async () => {
    const res = await Http.get('/visual/vm/list', { baseURL: '' })
    const res2 = await Http.get('/visual/worktask/editings', { baseURL: '' })
    const res3 = await Http.get('/visual/machine/list', { baseURL: '' })
    let machineMap = {}
    res3.data.forEach(t => {
      machineMap[t.id] = t.name
    })
    let taskMap = {}
    res2.data.tasks.forEach(t => {
      taskMap[t.task_id] = t
    })
    let treeData: any = []
    let treeMap = {}
    res.data = res.data || []
    for (let i = 0; i < res.data.length; i++) {
      let vm = new VirtualMachine(res.data[i])
      this.childrenMap[vm.id] = vm
      let task = taskMap[vm.editing_task_id]
      if (vm.editing && vm.editing_task_id && task) {
        vm.work_task = {
          id: task.task_id,
          user_id: task.user_id,
          link: task.link,
          status: task.status,
        }
      }
      let ipMachine = treeMap[vm.agent_id]
      if (ipMachine) {
        ipMachine.children.push(vm)
      } else {
        treeMap[vm.agent_id] = new VirtualMachine({
          id: vm.id,
          name: machineMap[vm.agent_id],
          agent_ip: vm.agent_ip,
          children: [vm],
        })
        treeData.push(treeMap[vm.agent_id])
      }
    }
    runInAction(() => {
      this.list = new Map(
        treeData.map(item => [item.id, new VirtualMachine(item)])
      )
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
