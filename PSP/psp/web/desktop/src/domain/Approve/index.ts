/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction, computed } from 'mobx'
import { AUDIT_REQUEST_TYPE } from '@/constant'
import { Approve, ALL_RESULT_MAP, APPROVED_MAP, statusColors } from './Approve'
import { Http } from '@/utils'
import { PageCtx } from '../common'

type LIST_TYPE = 'application' | 'pending' | 'complete'
export class ApproveList {
  optTypes = Object.values(AUDIT_REQUEST_TYPE)
  
  resultTypes = {
    application: ALL_RESULT_MAP,
    pending: ALL_RESULT_MAP,
    complete: APPROVED_MAP
  }

  @observable list: Approve[] = []
  @observable page_ctx: PageCtx = new PageCtx()

  @observable application_name = ''
  @observable approve_user_name = ''
  @observable type = ''
  @observable result = ''
  @observable start_time
  @observable end_time

  @observable page_index = 1
  @observable page_size = 10

  listType: LIST_TYPE = 'application'

  constructor(type) {
    this.listType = type
  }

  @action
  updateCurrentIndex = current => {
    this.page_index = current
  }

  @action
  updatePageSize = (current: number, size: number) => {
    this.page_size = size
    this.page_index = 1
  }

  @action
  updateList = list => {
    this.list = [...list]
  }

  find = id => this.list.find(item => item.id === id)

  fetch = () => {
    return Http.post(`/approve/list/${this.listType}`, {
      application_name: this.application_name === '' ? undefined : this.application_name,
      type: this.type === '' ? undefined : this.type,
      status: this.result === '' ? undefined : this.result,
      start_time: this.start_time,
      end_time: this.end_time,
      page: {
        index: this.page_index,
        size: this.page_size,
      }
    }).then(({ data: { list, page } }) => {
      runInAction(() => {
        this.updateList(list.map(item => new Approve(item)))
        this.page_ctx.update(page)

        let maxPageIndex = Math.ceil(this.page_ctx.total / this.page_size)

        if (this.page_index > maxPageIndex) {
          this.updateCurrentIndex(1)
        }
      })
    })
  }

  @action
  updateApplicationName(val) {
    this.application_name = val
  }

  @action
  uodateApproveUserName(val) {
    this.approve_user_name = val
  }

  @action
  updateResult(val) {
    this.result = val
  }

  @action
  updateTime(dates) {
    this.start_time = dates?.[0]
    this.end_time = dates?.[1]
  }

  @action
  updateOptType(type) {
    this.type = type
  }

  @action
  batchAccept = async (ids: string[]) => {
    // TODO
  }

  @action
  batchReject = async (ids: string[]) => {
     // TODO
  }
}

export const allApproveList = new ApproveList('application')
export const unapproveList = new ApproveList('pending')
export const approvedList = new ApproveList('complete')
export const statusColorMap = statusColors
