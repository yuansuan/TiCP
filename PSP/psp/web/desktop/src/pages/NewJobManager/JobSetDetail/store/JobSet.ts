/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Timestamp } from '@/domain/common'

class BaseJobSet {
  @observable id: string
  @observable project_id: string
  @observable name: string
  @observable state: number
  @observable display_state: number
  @observable has_failed: boolean
  @observable creator: string
  @observable is_batch_job: boolean
  @observable workflow_id: string
  @observable count: number
  @observable user_name: string
  @observable create_time: Timestamp
  @observable update_time: Timestamp
  @observable finish_time: Timestamp
}

export type JobSetRequest = Omit<
  BaseJobSet,
  'create_time' | 'update_time' | 'finish_time'
> & {
  create_time: {
    seconds: number
    nanos: number
  }
  update_time: {
    seconds: number
    nanos: number
  }
  finish_time: {
    seconds: number
    nanos: number
  }
}

export class JobSet extends BaseJobSet {
  constructor(props?: Partial<BaseJobSet>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({
    create_time,
    update_time,
    finish_time,
    ...props
  }: Partial<BaseJobSet>) => {
    Object.assign(this, props)

    if (create_time) {
      this.create_time = new Timestamp(create_time)
    }
    if (update_time) {
      this.update_time = new Timestamp(update_time)
    }
    if (finish_time) {
      this.finish_time = new Timestamp(finish_time)
    }
  }
}
