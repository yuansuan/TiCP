import { observable, action, computed } from 'mobx'
import moment from 'moment'
interface IUsage {
  id: number
  company_name: string
  licenses: number
  app_name: string
  app_id: string
  job_id: string
  job_name: string
  create_time: string
}

export class Usage implements IUsage {
  @observable id
  @observable app_id
  @observable app_name: string
  @observable company_name: string
  @observable job_id: string
  @observable job_name: string
  @observable licenses
  @observable create_time

  constructor(obj: IUsage) {
    this.init(obj)
  }

  @action
  init(obj) {
    if (!obj) return
    Object.assign(this, {
      ...obj,
      create_time: moment(obj.create_time).format('YYYY-MM-DD HH:mm:ss')
    })
  }
}
