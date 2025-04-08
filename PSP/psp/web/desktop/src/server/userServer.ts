/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

export const userServer = {
  current: () => Http.get('/auth/getLoginUserList'),
  checkPassword: (password: string) =>
    Http.post('/user/check_password', {
      password
    }),
  updatePassword: ({ name, password, newPassword }) =>
    Http.put('/user/updatePassword', {
      name,
      password,
      newPassword
    }),
  checkWxBind: (notification_type: 'job' | 'balance') =>
    Http.get('/user/wx/bind/check', {
      params: { notification_type }
    }),
  // wx解绑
  unbindWx: (notification_type: 'job' | 'balance', wechat_open_id) =>
    Http.delete('/user/wx/bind', {
      params: {
        notification_type,
        wechat_openid: wechat_open_id
      }
    }),
  // 获取微信二维码
  getWxCode: (notification_type: 'job' | 'balance') =>
    Http.post('/user/wx/qrcode', {
      notification_type
    }),
  updateRealName: real_name =>
    Http.put('/user/realname', {
      real_name
    }),
  updatePhone: ({ phone, oldPhone, captcha, token }) =>
    Http.put('/user/phone', {
      phone,
      captcha,
      token,
      oldPhone
    }),
  // 获取用户的企业邀请列表
  getInviteList: (params: {
    status?: number
    page_index: number
    page_size: number
  }) =>
    Http.get('/platform_user/user_invite_list', {
      params
    }),


  getShareUnreadCount: () => Http.get('storage/share/count', {
    params: {
      state:  1
    }
  }),
  getShareList: (params: {
    index: number
    size: number
  }) =>
    Http.post('/storage/share/recordList', {
        page: params
      }
    )
}
