import { observable, action } from 'mobx'

export class Preset {
  @observable id: string
  @observable name: string
  @observable default_preset: boolean
}
class BaseSoftware {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable icon: string
  @observable platform: number
  @observable display: number
  @observable image_id: string
  @observable init_script: string
  @observable gpu_desired: boolean
  @observable state: string
  @observable enabled: boolean
  @observable presets: Preset[]
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
