/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { newBoxServer } from '@/server'

export default class NewJob {
  async downloadJobs(jobs: { jobId: string; jobName: string }[]) {
    const paths = jobs.map(job => job.jobId)
    const jobNamesCountMap = {}
    const path_rewrite = jobs.reduce((prev, curr) => {
      // 处理同名作业
      if (!jobNamesCountMap[curr.jobName]) {
        jobNamesCountMap[curr.jobName] = 1
      } else {
        jobNamesCountMap[curr.jobName] += 1
      }
      const count = jobNamesCountMap[curr.jobName]
      prev[curr.jobId] =
        `${curr.jobName}` + (count - 1 > 0 ? `(${count - 1})` : '')
      return prev
    }, {})
    newBoxServer.download({
      base: '.',
      path_rewrite,
      bucket: 'result',
      paths
    })
  }
}
