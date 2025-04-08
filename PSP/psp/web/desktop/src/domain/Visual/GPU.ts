import { observable, computed } from 'mobx'
export default class GPU {
  @observable public id: number
  @observable public name: string
  @observable public agent_id: string
  @observable public used: boolean
  @observable public used_by: string

  @observable public domain: string
  @observable public bus: string
  @observable public slot: string
  @observable public function: string

  @observable public vendor_id: string
  @observable public device_id: string

  @computed get domain_name() {
    return `${this.domain}:${this.bus}:${this.slot}:${this.function}`
  }
  @computed get pci_id() {
    return `${this.vendor_id}:${this.device_id}`
  }

  constructor(request?: any) {
    if (request) {
      Object.assign(this, request)
    }
  }
}
