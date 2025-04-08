/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { IUser } from '@/domain/CompanyUsers/User'

export type RowData = Omit<
  IUser,
  'create_time' | 'update_time' | 'last_login_time'
> & {
  roles: string
  create_time: string
  update_time: string
  last_login_time: string
}
