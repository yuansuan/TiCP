/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed, runInAction } from 'mobx'
import { Company, BaseCompany } from './Company'
import { LAST_COMPANY_ID } from '@/constant'
import { companyServer } from '@/server'

export class BaseList {
  @observable list: Company[] = []
  @observable currentId: string
}

type IRequest = Omit<BaseList, 'list'> & {
  list: BaseCompany[]
}

export class CompanyList extends BaseList {
  constructor(props?: Partial<IRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({ list, ...props }: Partial<IRequest>) => {
    Object.assign(this, props)

    if (list) {
      this.list = [...list].map(item => new Company(item))
    }
  }

  @computed
  get current() {
    return this.list.find(item => item.id === this.currentId)
  }

  fetch = async () => {
    const { data } = await companyServer.getList()

    runInAction(() => {
      this.update({
        list: data,
      })
    })
  }

  change = async id => {
    const { data } = await companyServer.get(id)

    localStorage.setItem(LAST_COMPANY_ID, data?.id)
    runInAction(() => {
      this.update({
        currentId: id,
      })
      // update company
      const company = this.list.find(item => item.id === id)
      company && company.update(data)
    })
  }

  removeLastRecord = async () => {
    localStorage.removeItem(LAST_COMPANY_ID)
  }
}
