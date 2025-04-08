/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Session } from './Session'
import { getSessionList, SessionListParams } from '@/api/Session'

export class BaseSessionList {
  @observable list: Session[] = []
  @observable total: number = 0
}

export class SessionList extends BaseSessionList {
  @action
  update({ total, sessions }) {
    if (sessions) {
      this.list = sessions.map(session => new Session(session))
      this.total = total || 0
    }
  }

  fetch = async (params: Partial<SessionListParams>) => {
    const { data } = await getSessionList(params)

    runInAction(() => {
      this.update(data)
    })
  }
}
