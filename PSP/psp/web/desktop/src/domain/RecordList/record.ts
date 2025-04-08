/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action,computed} from 'mobx'
import { Http, history } from '@/utils'



export enum InvitedUserType {
  // 未知
  INVITE_TO_UNKNOW = 0,
  // 非管理员
  INVITE_NOT_ADMIN = 1,
  // 管理员
  INVITE_IS_ADMIN = 2
}

interface IRequest {
  share_time: string
  content: string
  id: string
  state: string
 
}

export class Record {
  @observable id
  @observable content
  @observable share_time
  @observable state


  @computed
  get timeTitle() {
    return new Date(this.share_time).toLocaleString()
  }
  constructor(props: IRequest) {
    if (props) {
      this.init(props)
    }
  }


  @action
  init = (props: IRequest) => {
    Object.assign(this, props)
  }

  readShare = async id => {
    await Http.put('/storage/share/updateRecordState', {
      record_ids: [this.id]
    })
  }

}
