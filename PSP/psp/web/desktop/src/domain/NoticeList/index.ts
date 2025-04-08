/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Notice, NoticeRequest } from './Notice'
import { Http } from '@/utils'
import moment from 'moment'
import { formatUnixTime } from '@/utils'

export class BaseNoticeList {
  @observable list: Notice[] = []
}

type IRequest = Omit<BaseNoticeList, 'list'> & {
  list: NoticeRequest[]
}

export class NoticeList extends BaseNoticeList {
  @action
  update({ list, ...props }: Partial<IRequest>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Notice(item))
    }
  }

  fetch = async () => {
    const { data } = await Http.get('/notice', {
      params: {
        filter: [
          'published||eq||true',
          `start_time||lte||${formatUnixTime(moment().unix())}`,
          `end_time||gte||${formatUnixTime(moment().unix())}`,
        ],
      },
    })

    runInAction(() => {
      this.update({
        list: data || [],
      })
    })
  }
}
