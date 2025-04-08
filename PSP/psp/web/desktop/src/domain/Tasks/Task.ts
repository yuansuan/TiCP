/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { v4 as uuidv4 } from 'uuid'

class BaseTask {
  id = uuidv4()
  @observable workTaskId: number
  @observable userId: string
  @observable appName: string
}

export type TaskRequest = BaseTask

export class Task extends BaseTask {
  constructor(props?: Partial<TaskRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<TaskRequest>) => {
    Object.assign(this, props)
  }

  @action
  reset = () => {
    this.update({
      workTaskId: undefined,
      userId: undefined,
      appName: undefined,
    })
  }
}
