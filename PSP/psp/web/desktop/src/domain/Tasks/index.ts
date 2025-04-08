/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Task, TaskRequest } from './Task'

class BaseTasks {
  @observable list: Task[] = []
}

export type TasksRequest = Omit<BaseTasks, 'list'> & {
  list: Partial<TaskRequest>[]
}

export class Tasks extends BaseTasks {
  constructor(props?: Partial<TasksRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({ list, ...props }: Partial<TasksRequest>) => {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Task(item))
    }
  }
}
