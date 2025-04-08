/**
 * @module Uploader
 * @description patch upload file/upload progress/task dispatcher
 */
import { observable, action, computed } from 'mobx'
import { filter } from 'rxjs/operators'
import { AsyncParallelHook } from 'tapable'

import { createMobxStore } from '@/utils'
import Task, { TaskStatus } from './Task'

interface ITask {
  target: Task
  timestamp: number
}

interface IUploader {
  tasks: Map<number, ITask>
}

const RUNNING_THRESHOLD = 5

export class Uploader implements IUploader {
  @observable tasks = new Map()
  hooks

  @computed
  get running() {
    return [...this.tasks.values()].filter(({ target: { status } }) =>
      [TaskStatus.uploading].includes(status)
    )
  }

  @computed
  get queuing() {
    return [...this.tasks.values()].filter(
      ({ target: { status } }) => status === TaskStatus.inited
    )
  }

  @computed
  get runnable() {
    return this.running.length < RUNNING_THRESHOLD
  }

  constructor() {
    this.hooks = {
      onStartUpload: new AsyncParallelHook(),
      onEndUpload: new AsyncParallelHook()
    }

    // autorun tasks
    createMobxStore(() => ({
      queuingCount: this.queuing.length,
      runningCount: this.running.length,
      runnable: this.runnable
    }))
      .pipe(filter(({ runnable }) => !!runnable))
      .subscribe(({ queuingCount, runningCount }) => {
        const num = Math.min(queuingCount, RUNNING_THRESHOLD - runningCount)

        if (num) {
          const waitingList = [...this.tasks.values()]
            .filter(({ target: { status } }) => status === TaskStatus.inited)
            .sort(
              ({ timestamp: prevTimestamp }, { timestamp }) =>
                prevTimestamp - timestamp
            )
            .slice(0, num)
            .map(item => item.target)

          waitingList.forEach(task => task.start())
        }
      })
  }

  get = taskId => this.tasks.get(taskId)
  clear = () => [...this.tasks.values()].forEach(({ target }) => target.abort())

  // patch upload
  // Antd.Upload.customRequest
  @action
  upload = ({ file, data, action }, isDir) => {
    file.customPath = file.webkitRelativePath || file.customPath || file.name
    if (data && data.dirPath) {
      file.uploadPath = `${data.dirPath}/${file.customPath}`
    }
    const task = new Task({
      file,
      action,
      params: data && data.params,
      dirPath: data && data.dirPath,
      isDir
    })

    // monitor the task done and aborted
    const subscription = task.status$.subscribe((status: TaskStatus) => {
      if (status === TaskStatus.done || status === TaskStatus.aborted) {
        this.tasks.delete(task.id)
        subscription.unsubscribe()
        this.hooks.onEndUpload.callAsync(() => {})
      }
    })

    this.tasks.set(task.id, {
      target: task,
      timestamp: Date.now()
    })

    // trigger onstartUpload hook
    this.hooks.onStartUpload.callAsync(() => {})

    return task
  };

  *[Symbol.iterator]() {
    yield* [...this.tasks.values()].map(item => item.target).values()
  }
}

export default new Uploader()

export { default as Task } from './Task'
