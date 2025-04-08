/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import moment from 'moment'
import { currentUser } from '@/domain'

export enum CloudJobStatus {
  init = 'init',
  uploading = 'uploading',
  upload_done = 'upload_done',
  cloud = 'cloud_status', // use job status 因为后端已经将云作业状态同步到该字段
  done = 'back_to_local',
  passbacking = 'passbacking',
  passback_done = 'passback_done'
}

export const CloudJobStatusMap = {
  init: '准备中',
  uploading: '上传中',
  upload_done: '上传完成',
  cloud_status: '云作业状态', // use job status 因为后端已经将云作业状态同步到该字段
  back_to_local: '爆发结束',
  passbacking: '回传中',
  passback_done: '回传完成'
}

export class BaseJob {
  @observable id: string //作业ID
  @observable name: string // 作业名称
  @observable app_name: string // 应用名称
  @observable app_id: string // 应用ID
  @observable user_name: string // 用户名称
  @observable cluster_name: string // 集群名称
  @observable type: string // '作业类型：local：本地作业；cloud：云端作业',
  @observable real_job_id: string // 调度器作业ID
  @observable out_job_id: string // 外部接口作业ID
  @observable queue: Array<any> // 作业队列
  @observable project_name: string // 项目名称
  @observable project_id: string //
  @observable state: string // 作业状态
  @observable raw_state: string // 作业原始状态
  @observable exit_code: string // 作业退出码
  @observable file_filter_regs: string[] // 过滤文件正则
  @observable priority: string // 作业优先级
  @observable cpus_alloc: string // 作业已分配核数
  @observable data_state: string // 数据状态
  @observable mem_alloc: string // 作业已分配内存(MB)
  @observable exec_duration: string // 作业实际运行时长(秒)
  @observable exec_host_num: string // 作业执行节点数量
  @observable reason: string // 作业状态原因
  @observable work_dir: string // 作业工作目录
  @observable exec_hosts: string // 作业执行节点名称
  @observable submit_time: string // 作业提交时间
  @observable pend_time: string // 作业等待时间
  @observable start_time: string // 作业开始时间
  @observable end_time: string // 作业结束时间
  @observable terminate_time: string // 作业终止时间
  @observable suspend_time: string // 作业暂停时间
  @observable create_time: string // 创建时间
  @observable update_time: string // 更新时间
  @observable isCloud: boolean // 是否云上作业
  @observable isSyncToLocal: boolean // 云上作业是否同步本地
  @observable enable_residual: boolean // 是否开启残差图
  @observable enable_snapshot: boolean // 是否开启云图
  @observable timelines: Array<{ name: string; progress: number; time: string }>
  @observable user_id: string // 用户ID
}

export type JobRequest = BaseJob

export class Job extends BaseJob {
  constructor(props?: Partial<JobRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({
    create_time,
    start_time,
    end_time,
    submit_time,
    ...props
  }: Partial<JobRequest>) {
    this.isCloud = props.type === 'cloud' ? true : false
    this.isSyncToLocal = props.data_state === 'Downloaded'
    Object.assign(this, props)

    if (create_time) {
      this.create_time = moment(create_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (start_time) {
      this.start_time = moment(start_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (end_time) {
      this.end_time = moment(end_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (submit_time) {
      this.submit_time = moment(submit_time).format('YYYY-MM-DD HH:mm:ss')
    }
  }

  get terminalable() {
    const canTerminalByState =
      this.state === 'Pending' ||
      this.state === 'Running' ||
      this.state === 'Suspended'

    const canTerminalBySysAdmin = currentUser.hasSysMgrPerm ? true : false

    const canTerminalByYou = this.user_id === currentUser.id

    return canTerminalByState && (canTerminalBySysAdmin || canTerminalByYou)
  }

  get resubmittable() {
    let canResubmitByState =
      this.state === 'Completed' ||
      this.state === 'Failed' ||
      this.state === 'Terminated' ||
      this.state === 'BurstFailed'

    if (this.type === 'cloud') {
      canResubmitByState =
        this.data_state === 'Downloaded' || this.data_state === 'DownloadFailed'
    }

    const canResubmitByYou = this.user_id === currentUser.id

    return canResubmitByState && canResubmitByYou
  }
}
