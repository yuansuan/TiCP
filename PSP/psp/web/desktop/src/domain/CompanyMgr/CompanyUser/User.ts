import { action, observable } from 'mobx'
import { formatTime } from '@/utils/formatter'

export enum CompanyUserStatus {
  UNKNOWN = 'UNKNOWN',
  NORMAL = 'NORMAL',
  DELETED = 'DELETED',
}

interface ICompanyUser {
  user_id?: string | null
  company_id?: string | null
  real_name?: string | null
  phone?: string | null
  email?: string | null
  status?: CompanyUserStatus | null
  is_admin?: boolean | null
  create_time?: number | null
  update_time?: number | null
  user_name?: string | null
}

export class User implements ICompanyUser {
  @observable user_id
  @observable company_id?: string | null
  @observable real_name?: string | null
  @observable phone?: string | null
  @observable email?: string | null
  @observable status?: CompanyUserStatus | null
  @observable is_admin?: boolean | null = false
  @observable create_time?: number | null
  @observable update_time?: number | null
  @observable user_name?: string | null

  get role() {
    return this.is_admin ? 'admin' : 'normal'
  }

  set role(value) {
    this.is_admin = value === 'admin' ? true : false
  }

  get create_time_string() {
    return this.create_time !== null ? formatTime(this.create_time) : null
  }

  get update_time_string() {
    return this.update_time !== null ? formatTime(this.update_time) : null
  }

  constructor(request?: ICompanyUser) {
    this.init(request)
  }

  @action
  public init = (request?: ICompanyUser) => {
    request && Object.assign(this, { ...request })
  }
}
