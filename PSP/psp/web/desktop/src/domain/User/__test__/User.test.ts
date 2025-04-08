/* Copyright (C) 2016-present, Yuansuan.cn */

import { User } from '../index'
import { dissoc } from 'ramda'

const initialUser = {
  ysid: 'ysid',
  name: 'name',
  account_id: 'account_id',
  email: 'email',
  display_user_name: 'display_user_name',
  user_name: 'user_name',
  real_name: 'real_name',
  public_box_domain: 'public_box_domain',
  phone: 'phone',
  headimg_url: 'headimg_url',
  wechat_union_id: 'wechat_union_id',
  wechat_nick_name: 'wechat_nick_name',
  wechat_open_id: 'wechat_open_id',
  wechat_nickname: 'wechat_nickname',
}

function getFieldsFromUser(model) {
  return {
    id: model.id,
    name: model.name,
    account_id: model.account_id,
    email: model.email,
    display_user_name: model.display_user_name,
    user_name: model.user_name,
    real_name: model.real_name,
    public_box_domain: model.public_box_domain,
    phone: model.phone,
    headimg_url: model.headimg_url,
    wechat_union_id: model.wechat_union_id,
    wechat_nick_name: model.wechat_nick_name,
    wechat_open_id: model.wechat_open_id,
    wechat_nickname: model.wechat_nickname,
  }
}

describe('@domain/User', () => {
  it('constructor with no params', () => {
    const model = new User()

    expect(getFieldsFromUser(model)).toEqual({
      id: undefined,
      name: undefined,
      account_id: undefined,
      email: undefined,
      display_user_name: undefined,
      user_name: undefined,
      real_name: undefined,
      public_box_domain: undefined,
      phone: undefined,
      headimg_url: undefined,
      wechat_union_id: undefined,
      wechat_nick_name: undefined,
      wechat_open_id: undefined,
      wechat_nickname: undefined,
    })
  })

  it('constructor with params', () => {
    const model = new User(initialUser)

    expect(getFieldsFromUser(model)).toEqual(
      dissoc('ysid', {
        ...initialUser,
        id: initialUser.ysid,
      })
    )
  })

  it('update', () => {
    const model = new User()
    model.update(initialUser)

    expect(getFieldsFromUser(model)).toEqual(
      dissoc('ysid', {
        ...initialUser,
        id: initialUser.ysid,
      })
    )
  })

  it('displayName', () => {
    const model = new User(initialUser)

    expect(model.displayName).toEqual(model.real_name)

    model.update({
      real_name: undefined,
    })
    expect(model.displayName).toEqual(model.display_user_name)

    model.update({
      display_user_name: undefined,
    })
    expect(model.displayName).toEqual(model.phone)
  })
})
