import { action, observable } from 'mobx'
import { Http } from '@/utils'

interface ISoftwareApp {
  id: number
  name: string
  version: string
  icon_data: string
  path: string
  os_type: string
  gpu_support: boolean
  app_param: string
  app_param_paths: string
  WS_list: Array<any>

  only_show_desktop: boolean
  process_name: string
}

export default class SoftwareApp implements ISoftwareApp {
  @observable public id: number = 0
  @observable public name: string = ''
  @observable public version: string = ''
  @observable public path: string = ''
  @observable public icon_data: string

  @observable public os_type: string = ''
  @observable public gpu_support: boolean = false
  @observable public app_param: string = ''
  @observable public app_param_paths: string = ''
  @observable public WS_list: Array<any> = []

  @observable public only_show_desktop: boolean
  @observable public process_name: string

  constructor(request?: ISoftwareApp) {
    this.init(request)
  }

  @action
  public init = (request?: ISoftwareApp) => {
    request &&
      Object.assign(this, {
        id: request.id,
        name: request.name,
        version: request.version,
        path: request.path,
        icon_data: request.icon_data,
        os_type: request.os_type,
        gpu_support: request.gpu_support,
        app_param: request.app_param,
        app_param_paths: request.app_param_paths,

        only_show_desktop: request.only_show_desktop,
        process_name: request.process_name,
      })
  }
}
