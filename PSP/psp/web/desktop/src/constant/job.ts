/* Copyright (C) 2016-present, Yuansuan.cn */
export const JOB_DISPLAY_STATE = {
  Pending: '等待中',
  Running: '运行中',
  Terminated: '已终止',
  Completed: '已完成',
  Failed: '已失败'
}

export const JOB_LIST_STATE = {
  Submitted: '已提交',
  Bursting: '爆发中',
  BurstFailed: '爆发失败'
}
export const JOB_STEP_STATE = {
  CloudRunning: '运行中',
  CloudTerminated: '已终止',
  CloudSubmitted: '已提交', //如果时间线中有LocalBursting，翻译为爆发完成
  CloudSuspended: '已暂停',
  CloudCompleted: '已完成',
  CloudFailed: '已失败',
  LocalBursting: '爆发中',
  LocalBurstFailed: '爆发失败',
  LocalRunning: '已运行',
  LocalTerminated: '已终止',
  LocalSuspended: '已暂停',
  LocalCompleted: '已完成',
  LocalFailed: '已失败',
  LocalSubmitted: '已提交',
  CloudDownloaded: '回传完成',
  CloudDownloadFailed: '回传失败'
}

export const ALL_JOB_STATES = {
  ...JOB_DISPLAY_STATE,
  ...JOB_STEP_STATE,
  ...JOB_LIST_STATE
}

export const DATASTATEMAP = {
  Downloading: {
    text: '回传中',
    type: 'primary'
  },
  Downloaded: {
    text: '回传完成',
    type: 'success'
  },
  DownloadFailed: { text: '回传失败', type: 'error' },
  Uploading: {
    text: '上传中',
    type: 'primary'
  },
  Uploaded: {
    text: '已上传',
    type: 'success'
  },
  UploadFailed: {
    text: '上传失败',
    type: 'error'
  }
}
