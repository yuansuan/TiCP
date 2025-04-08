import { observable, action, runInAction } from 'mobx'
import { companyList } from '@/domain'
import { Http } from '@/utils'

export class BaseVisIBVConfig {
  @observable isOpen: boolean
}

export default class VisIBVConfig extends BaseVisIBVConfig {
  constructor(props?) {
    super()
    !!props && Object.assign(this, props)
  }

  @action
  update(props?) {
    Object.assign(this, props)
  }

  @action
  async fetch() {
    const { data } = await Http.get('/vis_ibv/setting')
    runInAction(() => {
      this.update(data)
    })
  }

  get showVisIBVApp() {
    //在菜单中显示3D云应用（新的）的条件：1。企业用户。2。在OMS中开启了可视化应用
    return !!companyList.current && !!this.isOpen
  }
}
