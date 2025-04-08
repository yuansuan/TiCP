/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'
import { VisualHttp } from '@/domain/VisualHttp'

let bindCurrentUserCalled = false

export const visualServer = {
  fetch: async () => {
    const { data: visualizeSetting } = await Http.get('/visual/setting')
    const { data: terminalInfo } = await Http.get('/visual/terminal')
    return {
      data: {
        isOpen: visualizeSetting.isOpen,
        ...terminalInfo
      }
    }
  },
  bindCurrentUser: () => {
    if (bindCurrentUserCalled) {
      console.log('visualServer.bindCurrentUser already called, skip')
    } else {
      console.log('visualServer.bindCurrentUser')
      Http.post('/visual/bindCurrentUser')
      bindCurrentUserCalled = true
    }
  },
  fetchBundleUsages: async () => {
    const { data } = await Http.get('/visual/bundleUsages')
    return data
  },
  stop: async ({ user_id, work_task_id }) =>
    VisualHttp.post('/worktask/stop', { user_id, work_task_id })
}
