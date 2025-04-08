/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { VisualHttp } from '@/domain/VisualHttp'
import VisualTask from './VisualTask'

function getBase64(img, callback) {
  const reader = new FileReader()
  reader.addEventListener('loadend', () => callback(reader.result))
  reader.readAsDataURL(img)
}
export default class VisualTaskList {
  private scCache = new Map()
  @observable list: Array<VisualTask> = []
  @action
  async fetch() {
    const { data } = await VisualHttp.get('/worktask/my')
    const visualTaskList = data.worktask_list.map((item: any) => {
      const visualTask = new VisualTask(item)
      visualTask.updateScreenShot(this.scCache.get(visualTask.id) || '')
      return visualTask
    })
    runInAction(() => {
      this.list = visualTaskList
      this.fetchScreenShots()
    })
  }
  async exit(task: VisualTask) {
    await VisualHttp.post('/worktask/stop', {
      user_id: task.user_id,
      work_task_id: task.id,
    })

    await this.fetch()
  }

  @action
  onTaskMessage(message: any) {
    const task = new VisualTask(JSON.parse(message))
    const exist = this.list.find(t => {
      return t.id === task.id
    })
    if (!exist && task.status === 3) {
      //running task
      this.list.push(task)
    } else if (task.status === 4 || task.status === 5) {
      //failed & closed task
      this.list = this.list.filter(t => {
        return t.id != task.id
      })
    }
  }
  @action
  async fetchScreenShots() {
    await Promise.all(
      this.list.map(async task => {
        try {
          const res = await VisualHttp.get(
            `/worktask/${task.id}/thumbnail?${Date.now()}`,
            { responseType: 'blob', byPassInterceptor: true }
          )
          getBase64(res, (data: string) => {
            runInAction(() => {
              task.updateScreenShot(data)
              this.scCache.set(task.id, data)
            })
          })
        } catch (e) {}
      })
    )
  }
}
