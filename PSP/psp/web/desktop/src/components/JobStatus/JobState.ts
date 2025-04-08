/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

export enum JOB_STATE_ENUM {
  UNKNOWN = 0,
  UPLOAD = 1,
  COMPUTE = 2,
  BACK = 3,
  FINISHED = 4,
}

export class BaseJobState {
  @observable id: string
  @observable state: JOB_STATE_ENUM
  @observable progress: number
  @observable speed: number
  @observable description: string
  @observable total_size: number
}

export class JobState extends BaseJobState {
  constructor(props?: Partial<BaseJobState>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<BaseJobState>) => {
    Object.assign(this, props)
  }
}
