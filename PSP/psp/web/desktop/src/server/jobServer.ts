/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'
import { currentUser } from '@/domain'
export type ListParams = {
  page_index: number
  page_size: number
  user_names?: string[]
  app_names?: string[]
  states?: string[]
  queues?: string[]
  project_ids?: string[]
  is_admin?: boolean
  job_name?: string
  start_time?: string
  end_time?: string
  job_id?: string
}

export const jobServer = {
  get(id) {
    return Http.get(`/job/detail`, {
      params: {
        job_id: id
      }
    })
  },
  getJobSet(id) {
    return Http.get(`/job/jobSetDetail`, {
      params: {
        job_set_id: id
      }
    })
  },
  getStatus(id) {
    return Http.get(`job/status/${id}`)
  },
  list({ page_index = 1, page_size = 10, user_names, ...params }: ListParams) {
    const { isPersonalJobManager } = currentUser
    return Http.post(
      '/job/list',
      {
        page: {
          index: page_index,
          size: page_size
        },
        filter: {
          ...params,
          user_names: isPersonalJobManager ? [currentUser.name] : user_names
        }
      },
      {}
    )
  },
  getJobSetDetail({ id, page_index = 1, page_size = 10 }) {
    return Http.get(`/job/set/${id}`, {
      params: {
        page: {
          index: page_index,
          size: page_size
        }
      }
    })
  },
  async getAllFilters(): Promise<any> {
    const [res1, res2, res3, res4, res5] = await Promise.all([
      Http.get('/job/appNames', {}),
      Http.get('/job/userNames', {}),
      Http.get('/job/queueNames', {}),
      Http.get('/project/listForParam', {
        params: {
          is_admin: currentUser.hasSysMgrPerm
        }
      }),
      Http.get('/job/jobSetNames', {})
    ])

    return {
      app_names: res1?.data,
      user_names: res2?.data,
      queue_names: res3?.data,
      projects: res4?.data?.projects || [],
      job_set_names: res5?.data?.job_set_names || []
    }
  },
  async getAppNames(): Promise<any> {
    return await Http.get('/job/appNames', {})
  },
  async getUserNames(): Promise<any> {
    return await Http.get('/job/userNames', {})
  },
  async getQueueNames(): Promise<any> {
    return await Http.get('/job/queueNames', {})
  },
  async getProjects(): Promise<any> {
    return Http.get('/project/listForParam', {
      params: {
        is_admin: currentUser.hasSysMgrPerm
      }
    })
  },
  async getJobSetNames(): Promise<any> {
    return await Http.get('/job/jobSetNames', {})
  },

  terminate({ out_job_id, compute_type }) {
    return Http.post('/job/terminate', {
      out_job_id,
      compute_type
    })
  },
  resubmit(id) {
    return Http.post('/job/resubmit', {
      job_id: id
    })
  },
  cancel(ids: string[]) {
    return Http.post('/job/cancel', {
      ids
    })
  },
  delete(ids: string[]) {
    return Http.post('/job/delete', {
      ids
    })
  },
  getResidualData(id) {
    return Http.get(`/job/residual`, {
      params: {
        job_id: id
      }
    })
  },
  pauseOrResume(id, bool: Boolean) {
    return Http.put(`/job/${id}/sync_info`, {
      action: bool ? 'resume' : 'pause'
    })
  }
}
