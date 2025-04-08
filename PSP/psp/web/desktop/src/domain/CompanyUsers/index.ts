/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { User } from './User'

export enum UserQueryOrderBy {
  'ORDERBY_NULL' = 0,
  'ORDERBY_JOINTIMEDESC' = 1,
  'ORDERBY_JOINTIMEASC' = 2,
  'ORDERBY_LASTLOGINTIMEDESC' = 3,
  'ORDERBY_LASTLOGINTIMEDASC' = 4,
}

export class CompanyUsers {
  constructor(props?) {
    props && this.update(props)
  }

  @action
  update({ list, ...props }) {
    Object.assign(this, props)

    if (list) this.list = list.map(item => new User(item))
  }

  @observable list: User[] = []

  @action
  updateList = list => {
    this.list = [...list]
  }
}
