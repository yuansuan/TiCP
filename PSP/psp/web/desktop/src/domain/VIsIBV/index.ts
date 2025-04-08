import { observable, action } from 'mobx'
import { Software, SoftwareRequest } from './Software'
import { Hardware, BaseHardware } from './Hardware'

export { Hardware, Software }
export class BaseList {
  @observable softwareList: Software[] = []
  @observable hardwareList: Hardware[] = []
  @observable projectNames: string[]

  @observable page_ctx: {
    index: number
    size: number
    total: number
  } = {
    index: 1,
    size: 10,
    total: 0
  }
}

type IRequest = Omit<BaseList, 'list'> & {
  list: SoftwareRequest[] | BaseHardware[]
}

export class SoftwareList extends BaseList {
  constructor(props?: Partial<IRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<IRequest>) {
    Object.assign(this, props)

    if (list) {
      this.softwareList = list.map(item => new Software(item))
    }
  }
}

export class HardwareList extends BaseList {
  constructor(props?: Partial<IRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<IRequest>) {
    Object.assign(this, props)

    if (list) {
      this.hardwareList = list.map(item => new Hardware(item))
    }
  }
}

export const OPERATING_SYSTEM_PLATFORM_OF_ADD = {
  '1': 'LINUX',
  '2': 'WINDOWS'
}
