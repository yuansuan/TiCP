import { observable, runInAction } from 'mobx'
import { Http } from '@/utils'
import OpenedApp from './OpenedApp'

interface IAppList {
  list: Map<string, OpenedApp>
}

export default class OpenedAppList implements IAppList {
  *[Symbol.iterator]() {
    yield* this.list.values()
  }
  @observable list = new Map()

  get = id => this.list.get(id)
  remove = app => {
    Http.post(
      '/visual/worktask/stop',
      {
        user_id: app.userId,
        work_task_id: app.id,
      },
      { baseURL: '' }
    ).then(() => {
      this.list.delete(app.id)
    })
  }
  fetch = () =>
    Http.get('/visual/worktask/my', {
      params: {
        page_size: 10,
        page_num: 1,
      },
      baseURL: '',
    }).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data.worktask_list.map(item => [item.id, new OpenedApp(item)])
        )
      })
      return res
    })
}
