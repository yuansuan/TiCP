import { action, observable } from 'mobx'
import { formatTime } from '@/utils/formatter'

export enum CompanyStatus {
  COMPANY_UNKNOWN = 'COMPANY_UNKNOWN',
  COMPANY_NORMAL = 'COMPANY_NORMAL',
  COMPANY_DELETED = 'COMPANY_DELETED'
}

export const COMPANY_STATUS_MAP = {
  COMPANY_UNKNOWN: '全部',
  COMPANY_NORMAL: '正常',
  COMPANY_DELETED: '已删除'
}

interface ICompany {
  id?: string | null
  name?: string | null
  biz_code?: string | null
  is_ys_cloud?: number | null
  contact?: string | null
  phone?: string | null
  remark?: string | null
  status?: CompanyStatus | null
  account_id?: string | null
  modify_uid?: string | null
  modify_name?: string | null
  update_time?: any | null
  create_uid?: string | null
  create_name?: string | null
  create_time?: any | null
}

export default class Company implements ICompany {
  @observable id = null
  @observable name = ''
  @observable biz_code = null
  @observable is_ys_cloud = null
  @observable contact = null
  @observable phone = ''
  @observable remark = null
  @observable status = null
  @observable update_time = null
  @observable create_time = null

  get create_time_string() {
    return this.create_time !== null
      ? formatTime(this.create_time.seconds)
      : null
  }

  get update_time_string() {
    return this.update_time !== null
      ? formatTime(this.update_time.seconds)
      : null
  }

  constructor(request?: ICompany) {
    this.init(request)
  }

  @action
  public init = (request?: ICompany) => {
    request && Object.assign(this, { ...request })
  }
}
