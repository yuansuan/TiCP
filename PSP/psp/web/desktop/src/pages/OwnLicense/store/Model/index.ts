/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Item, BaseItem } from './Item'
import { byolServer } from '@/server'
export class BaseModel {
  @observable list: Item[] = []
}

type IRequest = Omit<BaseModel, 'list'> & {
  list: BaseItem[]
}

export type FetchParams = {
  key: string
}

export class Model extends BaseModel {
  constructor(props?: Partial<IRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<IRequest>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Item(item))
    }
  }

  fetch = async ({ key }: FetchParams) => {
    const { data } = await byolServer.get()
    runInAction(() => {
      this.update({
        list: data.filter(item =>
          item.app_name.toLowerCase().includes(key.toLowerCase())
        ),
      })
    })
  }
}
