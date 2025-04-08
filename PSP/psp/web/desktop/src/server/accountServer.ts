/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

export const accountServer = {
  get: () => Http.get('/account/detail')
}
