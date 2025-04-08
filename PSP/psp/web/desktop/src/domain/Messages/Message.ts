/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed, runInAction } from 'mobx'
import { Http } from '@/utils'

const messageTemplate = {
  job_event_type: params => params,
  share_file_event_type: params => params,
  project_event_type: params => params,
  session_event_type: params => params,
}

export const messageTypeMap = {
  job_event_type: '作业管理通知',
  share_file_event_type: '文件分享通知',
  project_event_type: '项目管理通知',
  session_event_type: '会话管理通知',
  approve_event_type: '审批通知',
}

interface IRequest {
  id: string
  type: string
  content: string
  state: number
  create_time: string
}

export class Message implements IRequest {
  @observable id: string
  @observable type: string
  @observable state: number // 1未读 2已读
  @observable content: string
  @observable title: string
  @observable create_time: string

  constructor(props?: IRequest) {
    if (props) {
      this.init(props)
    }
  }
  @computed
  get body() {
    return this.content
  }

  @computed
  get message() {
    if (messageTemplate[this.type]) {
      return messageTemplate[this.type](this.body)
    } else {
      return this.body
    }
  }

  @computed
  get timeTitle() {
    return new Date(this.create_time).toLocaleString()
  }

  @computed
  get typeToString() {
    return messageTypeMap[this.type]
  }

  @action
  init(props: IRequest) {
    Object.assign(this, props)
    this.title = messageTypeMap[this.type]
  }

  read = async () => {
    await Http.put('/notice/read', {
      message_ids: [this.id]
    })

    runInAction(() => {
      this.state = 2
    })
  }
}
