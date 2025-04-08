/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

export const appServer = {
  list: () => Http.get('/application'),
  get: id => Http.get(`/application/${id}`),
  spareResource: ({
    app_type,
    sc_ids
  }: {
    app_type: string
    sc_ids: string[]
  }) =>
    Http.post('/application/remainingResource', {
      app_type,
      sc_ids
    }),
  listSC: id => Http.get(`/application/cores/${id}`)
}
