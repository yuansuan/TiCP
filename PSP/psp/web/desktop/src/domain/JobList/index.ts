/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Job, JobRequest, CloudJobStatus, CloudJobStatusMap } from './Job'
export { CloudJobStatus, CloudJobStatusMap }
import { NewJobFileList } from '@/domain/JobList/NewJobFileList'

export const jobStatusColumnFields = [
  { key: 'run', status: 'Running', label: '运行', color: '#1890FF' },
  { key: 'susp', status: 'Suspend', label: '暂停', color: '#F7B500' },
  { key: 'pend', status: 'Pending', label: '等待', color: '#9013FE' },
  { key: 'exited', status: 'Exited', label: '退出', color: '#E02020' },
  { key: 'done', status: 'Done', label: '完成', color: '#52C41A' }
]

export const jobStatusColumnMap = jobStatusColumnFields.reduce(
  (p, s) => ((p[s.key] = s), p),
  {}
)
export class BaseJobList {
  @observable list: Job[] = []
  @observable page_ctx: {
    index: number
    size: number
    total: number
  } = {
    index: 1,
    size: 10,
    total: 0
  }
}

type JobListRequest = Omit<BaseJobList, 'list'> & {
  list: JobRequest[]
}

export class JobList extends BaseJobList {
  constructor(props?: Partial<JobListRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<JobListRequest>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Job(item))
    }
  }
}
export enum JOB_API_TYPE {
  NORMAL = 'NORMAL',
  SUBARRAY = 'SUBARRAY',
  HOST = 'HOST'
}

export enum JobTableType {
  LiveTable = 'LiveTable',
  FinishedTable = 'FinishedTable'
}
export const hostJobList = new NewJobFileList()
