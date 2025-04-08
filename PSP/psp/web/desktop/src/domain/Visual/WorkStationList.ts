import { observable, runInAction } from 'mobx'

import { Http } from '@/utils'
import WorkStation from './WorkStation'

interface IWorkStationList {
  list: Map<number, WorkStation>
}

export default class WorkStationList implements IWorkStationList {
  @observable list = new Map()

  get = id => this.list.get(id)

  fetch = () => {
    return Http.get('/visual/workstation', { baseURL: '' }).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data.workstation_list.map(item => [
            item.id,
            new WorkStation(item),
          ])
        )
      })
      return res
    })
  }

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
