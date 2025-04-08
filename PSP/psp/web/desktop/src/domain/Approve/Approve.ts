import { observable, action, computed, runInAction } from 'mobx'
import { Http } from '@/utils'
import { AUDIT_REQUEST_TYPE } from '@/constant'

export const OPT_TYPE_MAP = Object.values(AUDIT_REQUEST_TYPE).reduce((p, c) => (p[c.approve_type] = c.name, p), {})

export const RESULT_MAP = ['', '等待审批', '审批已通过', '审批已拒绝', '审批已撤销', '审批操作执行失败']

export const ALL_RESULT_MAP = {
  1: RESULT_MAP[1],
  2: RESULT_MAP[2],
  3: RESULT_MAP[3],
  4: RESULT_MAP[4],
  5: RESULT_MAP[5]
}

export const APPROVED_MAP = {
  2: RESULT_MAP[2],
  3: RESULT_MAP[3],
  5: RESULT_MAP[5]
}

export const statusColors = ['', 'blue', 'green', 'red', 'gray', 'red']

interface IApprove {
  id?: string | null
  application_name?: string // 申请人
  create_time?: string // 申请时间
  approve_time?: string // 审批时间
  approve_user_name?: string // 审批人
  result?: number // 审批结果
  type?: number // 操作类型
  record_id?: string 
  content?: string // 审批内容
  status?: number // 审批结果
  suggest?: string  // 审批备注
}

export class Approve implements IApprove {
  @observable id: string
  @observable application_name: string
  @observable create_time: string
  @observable approve_time: string
  @observable approve_user_name?: string | null
  @observable result?: number | null
  @observable type?: number | null
  @observable record_id?: string | null
  @observable content?: string | null
  @observable status?: number | null
  @observable suggest?: string | null


  @computed
  get result_str() {
    return RESULT_MAP[this.status] || '--'
  }

  @computed
  get opt_type_str() {
    return OPT_TYPE_MAP[this.type] || '--'
  }

  @computed
  get create_time_str() {
    return this.create_time ? this.create_time.substring(0, 19) : '--'
  }

  @computed
  get approve_time_str() {
    return this.approve_time ? this.approve_time.substring(0, 19) : '--'
  }

  constructor(props) {
    if (props) {
      this.init(props)
    }
  }

  @action
  init(props) {
    Object.assign(this, props)
  }

  cancel = async reason => {
    const res = await Http.post('/approve/cancel', {
      id: this.id
    })

    runInAction(() => {
      this.status = 4
      this.suggest = reason
    })

    return res
  }

  accept = async reason => {
    const res = await Http.post('/approve/pass', {
      id: this.id,
      approve_type: this.type,
      record_id: this.record_id,
      suggest: reason,
    })

    runInAction(() => {
      this.status = 2
      this.suggest = reason
    })

    return res
  }

  reject = async reason => {
    const res = await Http.post('/approve/refuse', {
      id: this.id,
      approve_type: this.type,
      record_id: this.record_id,
      suggest: reason,
    })

    runInAction(() => {
      this.status = 3
      this.suggest = reason
    })

    return res
  }
}
