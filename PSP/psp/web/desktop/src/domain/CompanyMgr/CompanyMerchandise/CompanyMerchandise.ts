import { action, observable, computed } from 'mobx'
import { formatTime } from '@/utils/formatter'

export enum CompanyMerchandiseStatus {
  STATE_UNKNOWN = 'STATE_UNKNOWN',
  STATE_ONLINE = 'STATE_ONLINE',
  STATE_OFFLINE = 'STATE_OFFLINE',
}

interface ICompanyMerchandise {
  id?: string | null
  company_id?: string | null
  company_name?: string | null
  merchandise_id?: string | null
  merchandise_name?: string | null
  out_resource_type?: number | null
  out_resource_id?: string | null
  license_type?: string | null
  license_active?: string | null
  state?: CompanyMerchandiseStatus | null
  create_uid?: string | null
  create_time?: any | null
  update_time?: any | null
  product_id?: string | null
}

export class CompanyMerchandise implements ICompanyMerchandise {
  @observable id
  @observable company_id: string | null
  @observable company_name: string | null
  @observable merchandise_id: string | null
  @observable merchandise_name: string | null
  @observable out_resource_type: number | null
  @observable out_resource_id: string | null
  @observable license_type: string | null
  @observable license_active
  @observable state
  @observable create_uid
  @observable create_time = null
  @observable update_time = null
  @observable product_id
  @observable price_active_time = null
  @observable price_unit = null

  @computed
  get price_active_time_string() {
    return this.price_active_time !== null
      ? formatTime(this.price_active_time.seconds)
      : null
  }

  @computed
  get create_time_string() {
    return this.create_time !== null
      ? formatTime(this.create_time.seconds)
      : null
  }

  @computed
  get update_time_string() {
    return this.update_time !== null
      ? formatTime(this.update_time.seconds)
      : null
  }

  constructor(request?: ICompanyMerchandise) {
    this.init(request)
  }

  @action
  public init = (request?: ICompanyMerchandise) => {
    request && Object.assign(this, { ...request })
  }
}
