import { observable, action } from 'mobx'

export class BaseHardware {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable charge_type: number
  @observable network_bandwidth: number
  @observable number_of_cpu: number
  @observable number_of_mem: number
  @observable number_of_gpu: string
  @observable cpu_model: string
  @observable gpu_model: string
  @observable zone: {
    name: string
    desc: string
  }
}

export class Hardware extends BaseHardware {
  constructor(props?: Partial<BaseHardware>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<BaseHardware>) {
    Object.assign(this, props)
  }
}
