/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction, computed } from 'mobx'
import { AsyncParallelHook } from 'tapable'

import { Task } from '@/domain/Uploader'

interface IUploadingDirectory {
  is_dir: true
  name: string
  path: string
  size: number
  status: 'uploading' | 'done'
  percent: number
  _master: string
  _from: 'local' | 'server'

  tasks: Task[]
}

export default class VirtualDirectory implements IUploadingDirectory {
  readonly is_dir = true
  @observable _master
  @observable _from

  @observable name = ''
  @observable path = ''
  @observable isDone = false

  @observable tasks: Task[] = []
  hooks

  constructor(props: Partial<IUploadingDirectory>) {
    this.init(props)

    this.hooks = {
      aborted: new AsyncParallelHook(),
    }
  }

  @action
  init = (props: Partial<IUploadingDirectory>) => {
    Object.assign(this, props)
  }

  @action
  abort = () => {
    const { tasks } = this

    this.tasks = []
    tasks.forEach(task => task.abort())

    this.hooks.aborted.callAsync(() => {})
  }

  @computed
  get status() {
    return this.isDone ? 'done' : 'uploading'
  }

  @computed
  get size() {
    return this.tasks.reduce((total, item) => total + item.target.size, 0)
  }

  @computed
  get loaded() {
    return this.tasks.reduce((total, item) => total + item.loaded, 0)
  }

  @computed
  get percent() {
    return (this.loaded / this.size) * 100
  }

  @action
  addTask = task => {
    this.tasks.push(task)

    // when upload complete, replace the local file with server file
    const subscription = task.status$.subscribe(status => {
      if (['done', 'error', 'aborted'].includes(status)) {
        // check if all files are uploaded
        const isDone = this.tasks.every(task =>
          ['done', 'aborted', 'error'].includes(task.status)
        )
        if (isDone) {
          runInAction(() => {
            // set the UploadingDirectory is done
            this.isDone = true
          })
        }

        subscription.unsubscribe()
      }
    })
  }
}
