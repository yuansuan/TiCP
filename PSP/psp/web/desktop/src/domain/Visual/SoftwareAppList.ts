import { observable, runInAction } from 'mobx'

import { Http } from '@/utils'
import SoftwareApp from './SoftwareApp'

interface ISoftwareAppList {
  list: Map<number, SoftwareApp>
  loading: boolean
}

export default class SoftwareList implements ISoftwareAppList {
  @observable list = new Map()
  @observable loading = false

  get = id => this.list.get(id)

  fetch = () => {
    this.loading = true
    return Http.get('/visual/app/list', { baseURL: '' }).then(res => {

      runInAction(async () => {
        this.list = new Map(await Promise.all(
          res.data?.map(async item => {
            let resWSList = null
            
            try {
              resWSList = await Http.get(`/visual/app/${item.id}/workstationList`, {
                baseURL: '',
              })
            } catch (e) {
              console.error('获取软件关联工作站失败')
            }
            
            const app = new SoftwareApp(item)
            app.WS_list = resWSList?.data || []

            return [
              item.id,
              app,
            ]
          }))
        )
        this.loading = false
      })
      return res
    })
  }

  delete = (rowData) => {
    return Http.delete(`/visual/app/${rowData['id']}`, { baseURL: '' })
  }

  add = (data) => {
    return Http.post(`/visual/app`, data , { baseURL: '' })
  }

  edit = (data) => {
    return Http.put(`/visual/app/${data['id']}`, data , { baseURL: '' })
  }

  bindOrUnBind = async (data) => {
    const { oldBindWSIds = [], newBindWSIds = [], id } = data
    // new - old 需要新绑定
    // old - new 需要解绑的
    const bindIds = newBindWSIds.filter(i => !oldBindWSIds.includes(i))
    const unbindIds = oldBindWSIds.filter(i => !newBindWSIds.includes(i))

    const bindRes = await Promise.all(bindIds.map(async bindId => {
      return Http.post(`/visual/app/${id}/bind/${bindId}`, {} , { baseURL: '' })
    }))

    const unbindRes = await Promise.all(unbindIds.map(async bindId => {
      return Http.delete(`/visual/app/${id}/unbind/${bindId}`, { baseURL: '' })
    }))
    
    return [...bindRes, ...unbindRes]
  }

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
