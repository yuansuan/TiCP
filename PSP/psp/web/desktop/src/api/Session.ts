/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Http } from '@/utils'
interface SessionParams {
  user_name: string
  hardware_id: string
  software_id: string
  status: string
}

export interface SessionListParams extends SessionParams {
  page_index: number
  page_size: number
}

export const getSessionList: (
  params?: Partial<SessionListParams>
) => Promise<any> = (
  params = {
    hardware_id: '',
    software_id: '',
    status: '',
    page_index: 1,
    page_size: 10
  }
) => {
  return Http.get('/vis/session', {
    params: {
      ...params
    }
  })
}

export const closeSession: (data?: {
  session_id: string
  reason: string
}) => Promise<any> = (
  data = {
    session_id: '',
    reason: ''
  }
) => {
  return Http.post('/vis/session/close', { ...data })
}
