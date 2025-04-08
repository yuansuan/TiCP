import { observable, runInAction, action } from 'mobx'
import currentUser from '@/domain/User'
import { Http } from '@/utils'
import Package from './Package'

interface IPackageList {
  list: Map<string, Package>
}

interface IBody {
  icon?: string // base64 string
}

export default class PackageList implements IPackageList {
  @observable packagePath = ''
  @observable packageName = ''
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0
  @observable orderBy = 'create_time'
  @observable orderAsc = false // 升序降序

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = current
    this.size = size
  }

  @action
  updateOrder(orderBy, orderAsc) {
    this.orderBy = orderBy
    this.orderAsc = orderAsc
  }

  get = id => this.list.get(id)

  fetch = () => {
    let url = `/localpackage/list?index=${this.index}&size=${this.size}&orderBy=${this.orderBy}&orderAsc=${this.orderAsc}&package_name=${this.packageName}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.packages?.map(item => {
            return [item.id, new Package(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  }

  getPackagePath = async () => {
    const res = await Http.get(`/localpackage/path`)
    this.packagePath = res.data.path
  }

  checkPackageInstallStatus = packageName => {
    return Http.get(`/localpackage/status?packageName=${packageName}`)
  }

  // check zip
  checkPackage = packageName => {
    return Http.get(
      `/localpackage/check?packageName=${packageName}&userName=${currentUser.name}`
    )
  }

  // 删除
  delete = (id, packageName) => {
    return Http.delete(`/localpackage/${id}?packageName=${packageName}`)
  }

  // 更新
  update = (id, body: IBody) => {
    return Http.put(`/localpackage/${id}`, { ...body })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
