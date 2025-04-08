/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, computed } from 'mobx'
import VirtualMachine from './VirtualMachine'
import GPU from './GPU'

export default class MachineNode {
  @observable public id: string = ''
  @observable public name: string = ''
  @observable public agent_ip: string = ''
  @observable public gpus: GPU[] = []
  @observable public children: Array<VirtualMachine> = []
  constructor(request?: any) {
    if (request) {
      this.id = request.id
      this.name = request.name
      this.agent_ip = request.agent_ip
    }
  }
  @computed public get agentIps() {
    if (!this.agent_ip) {
      return []
    } else {
      return this.agent_ip.split(',')
    }
  }
}
