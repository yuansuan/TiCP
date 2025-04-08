/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable, runInAction, computed } from 'mobx'
import Project, { BaseProject } from './Project'
import { history, Http } from '@/utils'
import { PageCtx } from '@/domain/common'
import { currentUser, env } from '@/domain'
import { LAST_PROJECT_ID } from '@/constant'

export class BaseProjectList {
  @observable list: Project[]
  @observable pageCtx: PageCtx = new PageCtx()
  @observable current = new Project()
  @observable icon_list: string[]
  @observable consume_limit_amount: number = 0
  @observable allow_submit_job_over_limit: boolean = true
  @observable consume_limit_enabled: boolean = false
}

export default class ProjectList extends BaseProjectList {
  @computed
  get consume() {
    return this.consume_limit_amount
  }

  @action
  async fetch() {
    let res
    if (!env.isPersonal && !env.isSpaceManager) {
      res = await Http.get('/project/list_by_company', {
        params: {
          company_id: env.company?.id,
          page_index: 1,
          page_size: 1000
        }
      })
    } else if (env.isSpaceManager) {
      res = await Http.get('/project/list_by_joined_user', {
        params: {
          page_index: 1,
          page_size: 1000
        }
      })
    } else {
      res = await Http.get('/project/list_by_user', {
        params: {
          user_id: currentUser.id,
          page_index: 1,
          page_size: 1000
        }
      })
    }

    const { data } = res
    runInAction(() => {
      this.updateData({
        ...data,
        pageCtx: data.page_ctx
      })
    })
  }

  async create(data) {
    if (!env.isPersonal) {
      return Http.post('/project/create_by_company', {
        ...data,
        account_id: env.company.account_id,
        remark: '工作空间描述。'
      })
    } else {
      return Http.post('/project/create_by_user', {
        ...data,
        account_id: currentUser.account_id,
        remark: '工作空间描述。'
      })
    }
  }

  // 指定project_id
  async delete(id) {
    let res
    if (!env.isPersonal) {
      res = await Http.delete(`/project/delete_by_company/${id}`)
    } else {
      await Http.delete(`/project/delete_by_user/${id}`)
    }

    this.fetch()

    return res.data
  }

  // 指定project_id
  async modify({ id, ...data }) {
    if (!env.isPersonal) {
      await Http.put(`/project/modify_by_company/${id}`, data)
    } else {
      await Http.put(`/project/modify_by_user/${id}`, data)
    }

    this.fetch()
  }

  // 指定project_id
  async getInfo({ id }) {
    if (!env.isPersonal) {
      return await Http.get(`/project/by_company/${id}`)
    } else {
      return await Http.get(`/project/by_user/${id}`)
    }
  }

  async getUsers({ id, page_index, page_size }) {
    if (!env.isPersonal) {
      return await Http.get(`/project/company/${id}/users`, {
        params: {
          page_index,
          page_size
        }
      })
    } else {
      return await Http.get(`/project/${id}/users`)
    }
  }

  @action
  updateList = list => {
    this.list = [...list]
  }

  @action
  updateData({
    list,
    icon_list,
    pageCtx,
    data
  }: Partial<BaseProjectList> & { data: BaseProject }) {
    if (list) this.list = list.map((props: BaseProject) => new Project(props))
    if (icon_list) this.icon_list = icon_list
    if (pageCtx) this.pageCtx.update(pageCtx)
    if (data) {
      this.allow_submit_job_over_limit = data.allow_submit_job_over_limit
      this.consume_limit_amount = data.consume_limit_amount
      this.consume_limit_enabled = data.consume_limit_enabled
      this.current.update(data)
    }
  }

  async change(id) {
    const { data } = await Http.get(`/project/${id}`)
    localStorage.setItem(LAST_PROJECT_ID, id)
    this.updateData({ data: data.project })
    return data
  }

  removeLastRecord = async () => {
    localStorage.removeItem(LAST_PROJECT_ID)
  }
}
