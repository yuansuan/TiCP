import { observable, action } from 'mobx'
import { Http } from '@/utils'
import Project from './Project'
import { currentUser } from '..'

export const PROJECT_STATE_MAP = {
  'Init': '初始化', 'Running': '进行中', 'Terminated': '已终止', 'Completed': '已结束'
}

export enum PROJECT_STATE_ENUM {
  'Init' = 'Init', 
  'Running' ='Running', 
  'Terminated' = 'Terminated', 
  'Completed' = 'Completed'
}

type PROJECT_STATE = keyof typeof PROJECT_STATE_MAP | ''

export class ProjectMG {
  @observable list: Project[] = []
  
  @observable projectName: string
  @observable startTime
  @observable endTime
  @observable state: PROJECT_STATE = ''
  @observable is_sys_menu: boolean = currentUser.hasSysMgrPerm

  @observable pageIndex: number = 1
  @observable pageSize: number = 10

  @observable total: number = 0

  @action
  async getList(isAdmin?) {
    const res = await Http.post('/project/list', {
      project_name: this.projectName,
      start_time: this.startTime?.unix(),
      end_time: this.endTime?.unix(),
      state: this.state ? [this.state] : null,
      is_sys_menu: typeof isAdmin === 'undefined' ? this.is_sys_menu : isAdmin,
      page: {
        index: this.pageIndex,
        size: this.pageSize
      }
    })

    this.total = res.data?.total || 0;

    this.list = res.data?.project_list?.map(item => {
      const project = new Project(item)
      return project
    })

    return res
  }

  async termination(project_id: string) {
    return Http.post(`/project/terminate`, {
      project_id,
    })
  }

  async delete(project_id: string) {
    return Http.post(`/project/delete`, {
      project_id,
    })
  }

  async add(data) {
    return Http.post(`/project/save`, {
      project_name: data.project_name,
      comment: data.comment,
      members: data.members, 
      project_owner: currentUser.id,
      start_time: data.start_time,
      end_time: data.end_time,
    })
  }

  async edit(data) {
    return Http.post(`/project/edit`, {
      project_name: data.project_name,
      comment: data.comment,
      members: data.members,
      project_id: data.project_id,
      project_owner: data.project_owner_id,
      start_time: data.start_time,
      end_time: data.end_time,
    })
  }

  async changeMembers(data) {
    return Http.post(`projectMember/save`, {
      project_id: data.project_id,
      user_ids: data.user_ids
    }) 
  }

  async changeOwner(data) {
    return Http.post(`/project/modifyOwner`, {
      project_ids: [data.project_id],
      target_project_owner_id: data.new_owner_id
    })
  }

  async getALLUsers() {
    const res = await Http.get(`/user/optionList`)
    return res
  }

  async getALLUsersWithProjectMgrPerm() {
    const res = await Http.get(`/user/optionList?filterPerm=6`)
    return res
  }
}

const adminProjectMG = new ProjectMG()
const personalProjectMG = new ProjectMG()

export { personalProjectMG, adminProjectMG }

