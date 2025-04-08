import { action, observable, computed } from 'mobx'

export interface IPackage {
  id: string
  icon: string
  name: string // app name
  package_name: string
  author: string // uploader
  verison: string
  create_time: any // upload time
  description: string
  licence_type: number
  licence_status: number
  app_type: number
}

export const licence_type_map = {
  0: '开源',
  1: '商用',
}

export const app_type_map = {
  0: '计算',
  1: '可视化',
  2: '计算和可视化',
}

export const licence_status_map = {
  0: '未上传',
  1: '已上传',
}

export default class Package implements IPackage {
  @observable public id = null
  @observable public icon = null
  @observable public name = ''
  @observable public package_name = ''
  @observable public author = ''
  @observable public verison = ''
  @observable public create_time = ''
  @observable public description = ''
  @observable public licence_type = 0
  @observable public licence_status = 0
  @observable public app_type = 0

  constructor(request?: IPackage) {
    this.init(request)
  }

  @computed
  get upload_time() {
    return this.create_time
  }

  @computed
  get license_type_str() {
    return licence_type_map[this.licence_type]
  }

  @computed
  get license_status_str() {
    return licence_status_map[this.licence_status]
  }

  @computed
  get app_type_str() {
    return app_type_map[this.app_type]
  }

  @action
  public init = (request?: IPackage) => {
    request && Object.assign(this, { ...request })
  }
}
