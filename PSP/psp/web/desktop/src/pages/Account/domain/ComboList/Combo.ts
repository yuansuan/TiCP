import { observable, action } from 'mobx'
import { Timestamp } from '@/domain/common'

export class BaseCombo {
  @observable combo_id: string
  @observable combo_name: string
  @observable chargeType: number
  @observable ticket_id: string
  @observable is_free: number
  @observable zone: string
  @observable softwares: any[]
  @observable hardwares: any[]
  @observable valid_begin_time: Timestamp
  @observable valid_end_time: Timestamp
  @observable remain_time: number // 剩余时长 在包小时级别的时候才会出现
  @observable used_time: number // 已经使用的时长 在包小时级别的时候才会出现
  @observable total_time: number
}

export type ComboRequest = Omit<
  BaseCombo,
  'valid_begin_time' | 'valid_end_time'
> & {
  valid_begin_time: {
    seconds: number
    nanos: number
  }
  valid_end_time: {
    seconds: number
    nanos: number
  }
}

export class Combo extends BaseCombo {
  constructor(props: ComboRequest) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({
    valid_begin_time,
    valid_end_time,
    ...props
  }: Partial<ComboRequest>) => {
    Object.assign(this, props)

    if (valid_begin_time) {
      this.valid_begin_time = new Timestamp(valid_begin_time)
    }
    if (valid_end_time) {
      this.valid_end_time = new Timestamp(valid_end_time)
    }
  }
}
