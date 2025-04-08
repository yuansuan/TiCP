import { action, observable } from 'mobx'
import { formatTime } from '@/utils/formatter'

export enum FundOperateStatus {
  FundOperateStatusUnknow = 'FundOperateStatusUnknow',
  FundOperateStatusSuccess = 'FundOperateStatusSuccess',
  FundOperateStatusFail = 'FundOperateStatusFail'
}

export const FundOperateStatus_MAP = {
  FundOperateStatusSuccess: '成功',
  FundOperateStatusFail: '失败'
}

export const FundOperateType_MAP = {
  FundOperateTypeAdd: '加款',
  FundOperateTypeSub: '扣款'
}

export enum FundOperateType {
  FundOperateTypeUnknow = 'FundOperateTypeUnknow',
  FundOperateTypeAdd = 'FundOperateTypeAdd',
  FundOperateTypeSub = 'FundOperateTypeSub'
}

interface IFundOperate {
  id?: string | null
  company_id?: string | null
  account_id?: string | null
  type?: FundOperateType | null
  amount?: number | null
  operator_uid?: string | null
  operator_name?: string | null
  status?: FundOperateStatus | null
  status_desc?: string | null
  remark?: string | null
  update_time?: any | null
  create_time?: any | null
}

export class FundOperate implements IFundOperate {
  @observable id?: string | null
  @observable company_id?: string | null
  @observable account_id?: string | null
  @observable type?: FundOperateType | null = FundOperateType.FundOperateTypeAdd
  @observable amount?: number | null = 0
  @observable operator_uid?: string | null
  @observable operator_name?: string | null
  @observable status?: FundOperateStatus | null
  @observable status_desc?: string | null
  @observable remark?: string | null
  @observable update_time?: any | null
  @observable create_time?: any | null

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

  constructor(request?: IFundOperate) {
    this.init(request)
  }

  @action
  public init = (request?: IFundOperate) => {
    request && Object.assign(this, { ...request })
  }
}
