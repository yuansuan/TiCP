/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

//TODO 完善请求参数类型

export const departmentServer = {
  getList: querys =>
    Http.get('/company/department/list', {
      params: {
        ...querys
      }
    }),
  getUsers: querys =>
    Http.get('/company/department/users', {
      params: {
        ...querys
      }
    }),
  create: body => Http.post('/company/department', body),
  edit: body => Http.put('/company/department', body),
  delete: data =>
    Http.delete('/company/delete_by_department', {
      data
    })
}
