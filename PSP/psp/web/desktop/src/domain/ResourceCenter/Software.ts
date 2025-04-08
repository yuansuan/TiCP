import { observable, action } from 'mobx'
import { BaseHardware } from './Hardware'

export class Preset {
  @observable id: string
  @observable hardware: BaseHardware
  @observable defaulted: boolean
}

class BaseSoftware {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable icon: string
  @observable platform: number
  @observable display: number
  @observable preset: Preset[]
  @observable zone: {
    name: string
    desc: string
  }
}

export type SoftwareRequest = Omit<BaseSoftware, ''> & {}

export class Software extends BaseSoftware {
  constructor(props?: Partial<SoftwareRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<SoftwareRequest>) {
    Object.assign(this, props)
  }
}
