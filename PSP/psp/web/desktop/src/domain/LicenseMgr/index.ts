import { observable, runInAction, action } from 'mobx'

import { AppInfo } from './AppInfo'
import { Http } from '@/utils'
import { Usage } from './Usage'

export class LicenseMgrList {
  @observable list: AppInfo[] = []
  @observable map: Map<string, AppInfo> = new Map()
  @observable total: number = 0

  // filter params
  @observable license_type: string

  // paging
  @observable index: number = 1
  @observable size: number = 10

  @action
  setFilterParams(licenseType) {
    this.license_type = licenseType
  }

  @action
  onPageChange(index, size) {
    this.index = index
    this.size = size
  }

  @action
  onSizeChange(index, size) {
    this.index = index
    this.size = size
  }

  get params() {
    return {
      license_type: this.license_type,
      // index: this.index,
      // size: this.size
    }
  }

  clear() {
    this.list = []
    this.map = new Map()
  }

  async fetch() {
    this.clear()

    const { data } = await Http.get('/licenseManagers', { params: this.params })

    runInAction(() => {
      this.list = data.license_managers.map(item => {
        this.map.set(item.id, new AppInfo(item))
        return new AppInfo(item)
      })
      this.total =  data.total
    })
  }

  async add(data) {
    return Http.post('/licenseManagers', data)
  }

  async edit(id, data) {
    return Http.put(`/licenseManagers/${id}`, data)
  }

  async publish(id) {
    return Http.put(`/licenseManagers/publish/${id}`)
  }

  async unpublish(id) {
    return Http.put(`/licenseManagers/unpublish/${id}`)
  }
}

export class LicenseUsageList {
  @observable list: Usage[] = []
  @observable map: Map<string, Usage> = new Map()
  @observable total: number = 0

  // filter params
  @observable company: string

  @observable index: number = 1
  @observable size: number = 10

  @action
  setFilterParams(company) {
    this.company = company
  }

  @action
  onPageChange(index, size) {
    this.index = index
    this.size = size
  }

  @action
  onSizeChange(index, size) {
    this.index = index
    this.size = size
  }

  get params() {
    return {
      company: this.company,
      index: this.index,
      size: this.size
    }
  }

  clear() {
    this.list = []
    this.map = new Map()
  }

  async fetch(id: string) {
    this.clear()

    const { data } = await Http.get(`/licenseMgr/usage/${id}`, {
      params: this.params
    })

    runInAction(() => {
      this.list = data?.result.map(item => {
        this.map.set(item.id, new Usage(item))
        return new Usage(item)
      })
      let { total } = data.page
      this.total = total
    })
  }
}
