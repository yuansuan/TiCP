/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'

import { Message } from './Message'
import { Http } from '@/utils'
import { PageCtx } from '../common'
export class Messages {
  @observable list: Message[] = []
  @observable page_ctx: PageCtx = new PageCtx()
  @observable unreadCount = 0
  @action
  updateListByNotification = (message: Message) => {
    this.list.unshift(message)
    this.list.pop()
  }

  @action
  updateList = list => {
    this.list = [...list]
  }
  @action
  updateUnreadCount = count => {
    this.unreadCount = count
  }

  find = id => this.list.find(item => item.id === id)

  fetchUnreadCount = () =>
    Http.get('/notice/count', {
      params: {
        state: 1,
      }
    }).then(res => {
      let { total } = (res as any).data
      this.updateUnreadCount(total)
      return res
    })
  fetch = (params: {
    page_index
    page_size
    filter?: { state: number; content: string }
  }) => {
    return Http.post(
      '/notice/list',
      {
        page: {
          index: params.page_index,
          size: params.page_size
        },
        filter: params.filter
      },
      {}
    ).then(({ data: { messages, page } }) => {
      this.updateList(messages.map(item => new Message(item)))
      this.page_ctx.update(page)
    })
  }

  fetchLast = () => {
    this.fetchUnreadCount()
    return this.fetch({
      page_index: 1,
      page_size: 5
    })
  }

  // 处理单独记录
  read = async (ids: string[]) => {
    await Http.put('/notice/read', {
      message_id: ids,
    })

    runInAction(() => {
      ids.forEach(id => {
        const item = this.find(id)
        if (item) {
          item.state = 2
        }
      })
    })

    // 数字减掉1
    this.updateUnreadCount(this.unreadCount - 1)
  }

  @action
  batchRead = async (ids: string[]) => {
    let count = 0

    ids.forEach(id => {
      const item = this.find(id)
      if (!item.state) {
        count++
      }
    })

    await Http.put(
      '/notice/read',
      { message_ids: ids }
    )

    // 将数字减掉
    this.updateUnreadCount(this.unreadCount - count)
  }

  readAll = async () => {
    await Http.put('/notice/readAll')

    runInAction(() => {
      this.list.forEach(item => {
        item.state = 2
      })
    })

    // 设置未读
    this.updateUnreadCount(0)
  }
}
