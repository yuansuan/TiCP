import { observable, runInAction, action } from 'mobx'
import currentUser from '@/domain/User'
import { Http } from '@/utils'
import { FundOperate } from './FundOperate'

interface IFundOperateList {
  list: Map<string, FundOperate>
}

interface IFundOperateBody {
  company_id?: string | null
  account_id?: string | null
  type?: string | null
  amount?: number | null
  remark?: string | null
}

export default class CompanyList implements IFundOperateList {
  @observable company_id
  @observable company_name
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = current
    this.size = size
  }

  get = id => this.list.get(id)

  fetch = () => {
    let url = `/company/fundoperate/list?index=${this.index}&size=${
      this.size
    }&company_id=${this.company_id || ''}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.id, new FundOperate(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  }

  // 创建
  create = (body: IFundOperateBody) => {
    return Http.post('/company/fundoperate', {
      ...body,
      operator_uid: currentUser.user_id,
      operator_name: currentUser.name
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
