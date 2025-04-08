import { observable, runInAction } from 'mobx'

import { Http } from '@/utils'
import WorkTask from './WorkTask'

interface IWorkTaskList {
  list: Map<number, WorkTask>
}

export default class WorkTaskList implements IWorkTaskList {
  @observable list = new Map()

  get = id => this.list.get(id)

  fetchWorkStationNames = () => {
    return Http.get('/visual/workstation', { baseURL: '' }).then(res => {
      let map = {}
      res.data.workstation_list.map(ws => {
        map[ws.id] = ws.name
      })
      return map
    })
  }
  fetchUserNames = identities => {
    return Http.post('/user/batch', { userIdentities: identities }).then(
      res => {
        let map = {}
        res.data.list.map(user => {
          map[user.id] = user.name
        })
        return map
      }
    )
  }
  fetch = async () => {
    const res = await Http.get('/visual/worktask', { baseURL: '' })
    let worktask_list = res.data.worktask_list
    let identities = []
    let identMap = {}
    worktask_list.map(wt => {
      if (!identMap[wt.user_id]) {
        identMap[wt.user_id] = true
        identities.push({ id: wt.user_id })
      }
    })
    /*
    let umap = {}
    if (identities.length !== 0) {
      umap = await this.fetchUserNames(identities)
    }*/

    const wmap = await this.fetchWorkStationNames()
    worktask_list.map(wt => {
      wt.workstation_name = wmap[wt.workstation_id]
      //wt.user_name = umap[wt.user_id]
    })

    runInAction(() => {
      this.list = new Map(
        worktask_list.map(item => [item.id, new WorkTask(item)])
      )
    })
    return worktask_list
  }
  remove = worktask => {
    Http.post(
      '/visual/worktask/stop',
      {
        user_id: worktask.userId,
        work_task_id: worktask.id
      },
      { baseURL: '' }
    ).then(() => {
      this.list.delete(worktask.id)
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
