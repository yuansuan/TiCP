import { Http } from '@/utils'
import { action, observable } from 'mobx'
import UserSession from './UserSession'

export const defaultPage = {
  index: 1,
  size: 10
}
export type Page = typeof defaultPage

export class UserSessionList {
  @observable sessionList: UserSession[] = []
  @observable page_index: number = 1
  @observable page_size: number = 10
  @observable total: number = 0

  @action
  getUserSessionList = async name => {
    const res = await Http.post(`/auth/onlineListByUser`, {
      user_name: name,
      page: {
        index: this.page_index,
        size: this.page_size
      }
    })
    this.sessionList = res.data?.list?.map(item => {
      return new UserSession(item)
    })
    this.total = res?.data?.page?.total
    return res
  }
  logoutByJwtToken = (jtis: string[], name: string) => {
    return Http.post('/auth/offlineByJti', {
      jti_list: jtis,
      user_name: name
    })
  }
}
export default new UserSessionList()
