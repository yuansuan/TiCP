/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, runInAction } from 'mobx'
import { Http } from '@/utils'
import VirtualMachine from './VirtualMachine'
import MachineNode from './MachineNode'
import GPU from './GPU'
export default class MachineNodeList {
  @observable public list: MachineNode[] = []

  find = (id: string) => {
    return this.list.find(node => {
      return node.id == id
    })
  }
  fetch = async () => {
    const res = await Http.get('/visual/machine/list', { baseURL: '' })
    // use new api: /visual/vm/info/list
    // old api: /visual/vm/list
    const res2 = await Http.get('/visual/vm/info/list', { baseURL: '' })
    const vmMap = {}
    res2.data = res2.data || []
    res2.data.forEach((item: any) => {
      const vm = new VirtualMachine(item)
      if (!vmMap[vm.agent_id]) {
        vmMap[vm.agent_id] = []
      }
      vmMap[vm.agent_id].push(vm)
    })
    res.data = res.data || []
    const list = res.data.map((item: any) => {
      let node = new MachineNode(item)
      if (vmMap[node.id]) {
        node.children = vmMap[node.id]
      }

      return node
    })
    for (var i = 0; i < list.length; i++) {
      const res = await Http.get(`/visual/machine/gpus/${list[i].id}`, {
        baseURL: '',
      })
      res?.data?.forEach((g: any) => {
        const gpu = new GPU(g)
        list[i].gpus.push(gpu)
      })
    }

    runInAction(() => {
      this.list = list
    })
  };
  [Symbol.iterator]() {
    return this.list[Symbol.iterator]
  }
}
