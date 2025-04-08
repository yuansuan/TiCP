import { observable, action } from 'mobx'

export class BaseHardware {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable instance_type: string
  @observable instance_family: string
  @observable network_bandwidth: number
  @observable number_of_cpu: number
  @observable number_of_mem: number
  @observable number_of_gpu: string
  @observable enabled: boolean
  @observable cpu_model: string
  @observable gpu_model: string
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
