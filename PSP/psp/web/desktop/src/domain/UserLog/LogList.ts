import { observable, runInAction, action } from 'mobx'
import { Http } from '@/utils'
import Log from './Log'
import { LOG_TYPE_MAP } from '@/constant'

interface ILogList {
  list: Map<string, Log>
}

export class LogList implements ILogList {
  @observable role_id = ''
  @observable user_name = ''
  @observable ip_address = ''
  @observable operate_type = ''
  @observable start_time = ''
  @observable end_time = ''

  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0
  @observable orderBy = 'operate_time'
  @observable orderAsc = false // 升序降序

  @observable LogOptionsTypesMap = LOG_TYPE_MAP

  @action
  updateUsername(username) {
    this.user_name = username
    this.index = 1
  }

  @action
  updateIPAddress(ip) {
    this.ip_address = ip
    this.index = 1
  }

  @action
  updateOptType(type) {
    this.operate_type = type
    this.index = 1
  }

  @action
  updateOptTime(dates) {
    this.start_time = dates?.[0]
    this.end_time = dates?.[1]
  }

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = 1
    this.size = size
  }

  @action
  updateOrder(orderBy, orderAsc) {
    this.orderBy = orderBy
    this.orderAsc = orderAsc
  }

  get = id => this.list.get(id)

  fetch = () => {
    return Http.post('/auditlog/list', {
      page: {
        index: this.index,
        size: this.size
      },
      operate_type: this.operate_type ? +this.operate_type : null,
      user_name: this.user_name,
      ip_address: this.ip_address,
      start_time: this.start_time,
      end_time: this.end_time
    }).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.id, new Log(item)]
          })
        )
        this.totals = res.data?.page?.total || 0
      })
      return res
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
