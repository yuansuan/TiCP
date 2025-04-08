import { action, observable, computed } from 'mobx'
import moment from 'moment'


export interface ILog {
  id: string | null
  role_id: string | null
  user_name: string | null
  ip_address: string | null
  operate_type: string | null
  operate_content: string | null
  operate_time: any | null
}

export default class Log implements ILog {
  @observable public id = null
  @observable public role_id = null
  @observable public user_name = ''
  @observable public ip_address = ''
  @observable public operate_type = ''
  @observable public operate_content = ''
  @observable public operate_time

  constructor(request?: ILog) {
    this.init(request)
  }

  @computed
  get operate_time_str() {
    return this.operate_time ? moment(this.operate_time).format("YYYY-MM-DD HH:mm:ss") : '--'
  }

  @action
  public init = (request?: ILog) => {
    request && Object.assign(this, { ...request })
  }
}
