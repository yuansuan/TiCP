import { action, observable, computed } from 'mobx'

const statusMap = new Map([
  [1, '排队'],
  [2, '提交'],
  [3, '运行'],
  [4, '失败'],
  [5, '关闭'],
])
const colorMap = new Map([
  [1, '#F5A623'],
  [2, '#52C41A '],
  [3, '#4A90E2'],
  [4, '#D0021B '],
  [5, '#9B9B9B'],
])

export default class WorkTask {
  @observable public id: number = 0
  @observable public os: string = ''
  @observable public status: number = 1
  @observable public template_name: string = ''
  @observable public start_time: Date = new Date()
  @observable public close_time: Date = new Date()
  @observable public link: string = ''
  @observable public user_id: number = 0
  @observable public user_name: string = ''
  @observable public workstation_id: number = 0
  @observable public workstation_name: string = ''
  @observable public software_list: Array<any> = []

  @computed
  get statusName() {
    return statusMap.get(this.status)
  }
  @computed
  get statusColor() {
    return colorMap.get(this.status)
  }
  constructor(request?: any) {
    this.init(request)
  }

  @action
  public init = (request?: any) => {
    request &&
      Object.assign(this, {
        id: request.id,
        os: request.os,
        status: request.status,
        template_name: request.template_name,
        start_time: new Date(request.start_time),
        close_time: new Date(request.close_time),
        link: request.link,
        user_id: request.user_id,
        user_name: request.user_name,
        workstation_id: request.workstation_id,
        workstation_name: request.workstation_name,
        software_list: request.software_list,
      })
  }
}
