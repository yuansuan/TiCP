import { action, observable } from 'mobx'

interface IWorkStation {
  id: number
  name: string
  os: string
  type: string
  cloud_type: string
  queuing: number
  free: number
  up_limit: number
  used: number
  software_list: Array<any>
}

export default class WorkStation implements IWorkStation {
  @observable public id: number = 0
  @observable public name: string = ''
  @observable public os: string = ''
  @observable public type: string = ''
  @observable public cloud_type: string = ''
  @observable public queuing: number = 0
  @observable public free: number = 0
  @observable public up_limit: number = 0
  @observable public used: number = 0
  @observable public software_list: Array<any> = []

  constructor(request?: IWorkStation) {
    this.init(request)
  }

  @action
  public init = (request?: IWorkStation) => {
    request &&
      Object.assign(this, {
        id: request.id,
        name: request.name,
        os: request.os,
        type: request.type,
        cloud_type: request.cloud_type,

        queuing: request.queuing,
        free: request.free,
        up_limit: request.up_limit,
        used: request.used,
        software_list: request.software_list,
      })
  }
}
