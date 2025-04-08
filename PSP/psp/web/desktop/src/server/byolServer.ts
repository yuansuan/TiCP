/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

export const byolServer = {
  get: () => Http.get('/ownLicense/list'),
  getlicense: (merchandise_id?: string) =>
    Http.get('/ownLicense/license', {
      params: {
        merchandise_id,
      },
    }),
}
