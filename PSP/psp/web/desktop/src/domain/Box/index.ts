/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed } from 'mobx'

import Draft from './Draft'
import NewDraft from './NewDraft'
import Job from './Job'
import NewJob from './NewJob'
import Result from './Result'
import NewResult from './NewResult'
import { Http } from '@/utils'
import { DraftType } from './Draft'
export { NewDraft, DraftType }

export type BoxStatus = 'success' | 'error' | 'notfound' | 'warning'

interface IBox {
  url: string
  rawUrl: string
  token: string
  status: BoxStatus
  diskUsed: number
  diskTotal: number
}

export default class Box implements IBox {
  @observable url: string
  @observable rawUrl: string
  @observable token: string
  @observable status: BoxStatus = 'notfound'
  @observable diskUsed: number
  @observable diskTotal: number

  jobDraft = new Draft()
  redeployDraft = new Draft(DraftType.Redeploy)
  continuousDraft = new Draft(DraftType.Continuous)
  job = new Job()
  newJob = new NewJob()
  result = new Result()
  newResult = new NewResult()

  @action
  updateToken = token => {
    this.token = token
  }

  @computed
  get exhaust(): boolean {
    return this.diskUsed > this.diskTotal
  }

  init = async () => {
    const { data: token } = await Http.get('/box/token')
    this.updateToken(token)
  }

  update = (props: Partial<IBox>) => {
    if (props.hasOwnProperty('url')) {
      props.rawUrl = props.url
    }
    Object.assign(this, props)
  }
}
