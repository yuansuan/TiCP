import { action, observable } from 'mobx'
import { Http } from '@/utils'
import SysUser from './SysUser'

export class SysUserList {
  @observable sysuserList: SysUser[] = []
  @observable page_index: number = 1
  @observable page_size: number = 10
  @observable total: number = 0
  @observable filter_name: string = ''

  @observable query = {
    order_by: 'name',
    page: {
      index: this.page_index,
      size: this.page_size
    },
    sort_by: true
  }

  @action
  async getSysUserList(queryKey?: string) {
    const res = await Http.post('/auth/onlineList', {
      ...this.query,
      filter_name: queryKey
    })
    this.sysuserList = res.data?.list?.map(item => {
      return new SysUser(item)
    })
    this.total = res.data?.page?.total
    return res
  }

  logoutByUserName = (user_names: string[]) => {
    return Http.post('/auth/offlineByUserName', {
      user_name_list: user_names
    })
  }
}
export default new SysUserList()
