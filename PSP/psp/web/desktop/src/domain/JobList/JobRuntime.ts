/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Timestamp } from '@/domain/common'

export class BaseJobRuntime {
  @observable id?: string
  @observable download_task_id?: string
  @observable download_finished?: string
  @observable cpu_time?: number
  @observable start_time?: Timestamp
  @observable end_time?: Timestamp
  @observable have_residual?: boolean
  @observable resource_assign?: {
    cpus?: number
  }
  @observable server_params: any
}

export type JobRuntimeRequest = Omit<
  BaseJobRuntime,
  'start_time' | 'end_time'
> & {
  start_time: {
    seconds: number
    nanos: number
  }
  end_time: {
    seconds: number
    nanos: number
  }
}

export class JobRuntime extends BaseJobRuntime {
  constructor(props?: Partial<JobRuntimeRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({ start_time, end_time, ...props }: Partial<JobRuntimeRequest>) => {
    Object.assign(this, props)

    if (start_time) {
      this.start_time = new Timestamp(start_time)
    }

    if (end_time) {
      this.end_time = new Timestamp(end_time)
    }
  }
}
