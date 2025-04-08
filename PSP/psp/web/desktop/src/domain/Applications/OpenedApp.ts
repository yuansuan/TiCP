import { action, observable } from 'mobx'
export interface IRequest {
  id: string
  status: number
  type: string
  user_name: string
  user_id: string
  start_time: Date
  link: string
  full_url: string
  template_name: string
  software_list: Array<any>
}

interface IOpenedApp {
  id: string
  name: string
  status: number
  type: string
  userName: string
  userId: string
  link: string
  template_name: string
}

export default class OpendedApp implements IOpenedApp {
  @observable public id = ''
  @observable public name = ''
  @observable public status = 1
  @observable public icon = ''
  @observable public userName = ''
  @observable public userId = ''
  @observable public startTime = null
  @observable public fullUrl = ''
  @observable link = ''
  @observable type = ''
  @observable template_name = ''
  constructor(request?: IRequest) {
    this.init(request)
  }

  @action
  public init = (request?: IRequest) => {
    request &&
      Object.assign(this, {
        id: request.id,
        status: request.status,
        link: request.link,
        userName: request.user_name,
        userId: request.user_id,
        startTime: new Date(request.start_time),
        fullUrl: request.full_url,
        template_name: request.template_name,
        name: request.software_list[0] && request.software_list[0].name,
      })
  }
}
