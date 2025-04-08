import { Http } from '@/utils'
import { action, observable } from 'mobx'
import WhiteListInfo from './WhiteListInfo'

export class WhiteList {
  @observable list: WhiteListInfo[] = []

  @action
  async getWhileList() {
    const res = await Http.get('/sysconfig/whitelist')
    this.list = res.data.whitelist?.map(list => {
      return new WhiteListInfo(list)
    })
    return res
  }
}
export default new WhiteList()
