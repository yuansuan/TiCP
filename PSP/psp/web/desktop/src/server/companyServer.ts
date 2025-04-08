/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http, v2Http } from '@/utils'
import { UserQueryOrderBy } from '@/domain/CompanyUsers'

export const companyServer = {
  getList: () => Http.get('/company/user/company_list'),
  get: id => Http.get(`company/${id}`),
  invite: user_id => Http.post('/platform_user/invite', { user_id }),
  batchInvite: (params: {
    role_id: string
    phone_list: any
    department_id?: string
  }) => Http.post('/platform_user/batch_invite', params),
  delete: (params: {
    user_id: string
    company_id: string
    company_name: string
  }) =>
    Http.delete('/company/user', {
      params
    }),
  delete_perm_user_default: (params: {
    user_id: string
    company_id: string
    company_name: string
  }) =>
    Http.delete('/company/user_in_personal_setting', {
      params
    }),
  configUser: (params: {
    user_id: string
    role_id: string
    consume_limit: number
    department_id?: string
  }) => Http.put('/company/user/role', params),
  getInviteList: (params: {
    company_id: string
    status?: number
    page_index: number
    page_size: number
  }) => Http.get('/company/company_invite_list', { params }),
  queryUsers: ({
    company_id,
    key = '',
    order_by = UserQueryOrderBy.ORDERBY_NULL,
    status,
    page_index,
    page_size
  }) =>
    Http.get('/company/query_users', {
      params: {
        company_id: company_id,
        key,
        status,
        page_index,
        page_size,
        order_by
      }
    }),
  getUserRole: () => Http.get('/company/user_role'),
  getUserRoleInSetting: () => Http.get('/company/user_role_in_setting'),
  listConfig: configKeys =>
    Http.get('/company/config/list', {
      params: { configKeys }
    }),
  getDomain: (company_id: string) =>
    Http.post('/storage/domain', { company_id }),
  getSCList: () => v2Http.get('/sc_id')
}
