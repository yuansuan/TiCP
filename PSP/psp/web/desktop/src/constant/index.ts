/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import moment, { Moment } from 'moment'

const PREFIX = 'PLATFORM_'

export const PLATFORM_CURRENT_ZONE_KEY = `${PREFIX}CURRENT_ZONE`

export const EDITABLE_SIZE = 3 * 1024 * 1024
export const UPLOAD_CHUNK_SIZE = 5 * 1024 * 1024

export const LAST_COMPANY_ID = PREFIX + 'LAST_COMPANY_ID'
export const LAST_PROJECT_ID = PREFIX + 'LAST_PROJECT_ID'

export const DeployMode = {
  Local: 'local',
  Cloud: 'cloud',
  Hybrid: 'hybrid'
}

export enum FeatureEnum {
  CLOUD_APP = 'cloud_app'
}

export enum JobDraftEnum {
  JOB_DRAFT_STORE_KEY = 'JOB_DRAFT_STORE_KEY',
  JOB_REDEPLOY_DRAFT_STORE_KEY = 'JOB_REDEPLOY_DRAFT_STORE_KEY',
  JOB_CONTINUOUS_DRAFT_STORE_KEY = 'JOB_CONTINUOUS_DRAFT_STORE_KEY'
}

export enum CloudShellPermissionEnum {
  NONE,
  PENDING,
  APPROVED
}

export const GeneralDatePickerRange: Record<string, [Moment, Moment]> = {
  当月: [moment().startOf('month'), moment().endOf('month')],
  上月: [
    moment().subtract(1, 'month').startOf('month'),
    moment().subtract(1, 'month').endOf('month')
  ],
  前三月: [
    moment().subtract(3, 'month').startOf('month'),
    moment().subtract(1, 'month').endOf('month')
  ],
  本年度: [moment().startOf('year'), moment().endOf('year')],
  上一年度: [
    moment().subtract(1, 'year').startOf('year'),
    moment().subtract(1, 'year').endOf('year')
  ]
}

export enum JobFileTypeEnum {
  all,
  result,
  model,
  log,
  middle,
  others
}

export enum RESOURCE_TYPE {
  UNKNOWN = 0,
  COMPUTE_APP = 1,
  VISUAL_APP = 2,
  CLOUD_STORAGE = 3,
  SC_TEMINAL_APP = 4,
  IBV_SOFTWARE = 5,
  IBV_HARDWARE = 6,
  STANDARD_COMPUTE_APP = 7,
  BUNDLE_VISUAL_APP_ALL = 101,
  CLOUD_APP_COMBO = 102, // 3D云应用套餐（购买套餐）
  //对前端来说，账单关注使用情况
  CLOUD_APP_COMBO_USAGE = 103 // 3D云应用使用套餐
}

export * from './job'
export * from './custom'

export const VISUAL_TASK_STATUS_NAME = {
  0: '未知',
  1: '排队中',
  2: '提交中',
  3: '运行中',
  4: '已失败',
  5: '已关闭',
  6: '离线'
}

export const SPACE_MGR_ROLE_NAMES = ['普通用户', '外协人员']

export const CAN_CHANGE_SPACE_OWNER_ROLE_NAMES = ['管理员', '子管理员']

export const INNER_COMPANY_IDS = ['4mhk5sRjJNd', '49S5JQbqeb9']

if (process.env.NODE_ENV === 'development') {
  INNER_COMPANY_IDS.push('3P6Yp9xVJZW')
}

// 资源中心 操作系统平台
export const OPERATING_SYSTEM_PLATFORM = {
  0: 'ALL',
  1: 'LINUX',
  2: 'WINDOWS'
}
// 资源中心 软件形式
export const SOFTWARE_FORM = {
  0: '未知', //  UNKNOWN
  1: '桌面形式', //  DESKTOP
  2: '应用形式' //  APPLICATION
}

// 那些作业状态下可以显示可视化分析:残差图，监控项，云图等
export const SHOW_JOB_VISIBLE = [
  'running',
  'completing',
  'completed',
  'failed',
  'canceled'
]

// 单价类型
export enum CHARGE_TYPE {
  ALL_TYPE = 0,
  REAL_QUATITY_TYPE = 1, // 按量
  MONTHLY_TYPE = 2, // 包年包月
  HOURLY_TYPE = 3 // 包时
}
export const BILLING_TYPE_MAP = {
  0: '未知',
  1: '按量计费',
  2: '包年包月',
  3: '包时长'
}

export const AppNameMap = {
  FileManage: '文件管理',
  Calculator: '作业提交',
  JobManage: '作业管理',
  '3dcloudApp': '3D云应用',
  userlog: '用户日志',
  EnterpriseManage: '系统管理',
  'Starccm+': '3D-StarCCM+',
  'Starccm+Calc': '求解-StarCCM+'
}

export const RouterTogg = {
  files: 'FILEMANAGE',
  'new-job-creator': 'CALCUAPP',
  'new-jobs': 'JOBMANAGE',
  sys: 'ENTERPRISEMANAGE',
  messages: 'MESSAGES',
  'vis-session': '3DCLOUDAPP',
  'new-job': 'NEWJODETAIL',
  'new-job-set': 'NEWJOBSETDETAIL'
}

export const RouterWappMap = {
  files: 'fileManage',
  'new-job-creator': 'calculator',
  'new-jobs': 'jobManage',
  'vis-session': '3dcloudApp',
  'new-job': 'jobDetail',
  'new-job-set': 'jobSetDetail'
}

function format(str) {
  return str.replace(/([A-Z1-9])/g, '-$1').toLowerCase()
}

const fontSize = ['12px', '14px', '16px', '20px', '24px', '30px', '38px']

const _antTheme = {
  fontFamily: 'PingFangSC-Regular',
  btnBorderRadiusBase: '2px',

  primaryColor: '#005dfc',
  linkColor: '#3182ff',

  errorColor: '#f5222d',
  warningColor: '#ffa726',
  successColor: '#52c41a',
  infoColor: '#3182ff',

  borderColorBase: 'rgba(0,0,0,0.10)',
  borderColorSplit: 'rgba(0,0,0,0.10)',

  // disabled bg color
  backgroundColorBase: '#f5f5f5',
  disabledColor: 'rgba(0,0,0,0.25)',

  fontSizeSm: fontSize[0],
  fontSizeBase: fontSize[1],
  fontSizeLg: fontSize[2],
  heading1Size: fontSize[6],
  heading2Size: fontSize[5],
  heading3Size: fontSize[4],
  heading4Size: fontSize[3],
  heading5Size: fontSize[2]
}

export const antTheme = Object.entries(_antTheme).reduce(
  (res, [key, value]) => {
    res[format(key)] = value
    return res
  },
  {}
)

export const theme = {
  ..._antTheme,
  backgroundColorHover: '#F6F8FA',
  secondaryColor: '#3182FF',
  cancelColor: '#BFBFBF',
  cancelHighlightColor: '#8C8C8C',

  fontSize
}

export type Theme = typeof theme

export const DatePicker_FORMAT = 'YYYY-MM-DD HH:mm:ss'
export const DatePicker_SHOWTIME_FORMAT = 'HH:mm:ss'

export const AUDIT_REQUEST_TYPE = {
  USER_ADD: {
    type: 'USER',
    url: '/user/add',
    method: 'POST',
    name: '新建用户',
    approve_type: 1
  },
  USER_DEL: {
    type: 'USER',
    url: '/user/delete',
    method: 'DELETE',
    name: '删除用户',
    approve_type: 2
  },
  USER_EDIT: {
    type: 'USER',
    url: '/user/update',
    method: 'PUT',
    name: '编辑用户',
    approve_type: 3
  },
  USER_ACTIVE: {
    type: 'USER',
    url: '/user/active',
    method: 'PUT',
    name: '启用用户',
    approve_type: 4
  },
  USER_INACTIVE: {
    type: 'USER',
    url: '/user/inactive',
    method: 'PUT',
    name: '禁用用户',
    approve_type: 5
  },
  ROLE_ADD: {
    type: 'ROLE',
    url: '/role/add',
    method: 'POST',
    name: '新建角色',
    approve_type: 6
  },
  ROLE_DEL: {
    type: 'ROLE',
    url: '/role/delete',
    method: 'DELETE',
    name: '删除角色',
    approve_type: 7
  },
  ROLE_EDIT: {
    type: 'ROLE',
    url: '/role/update',
    method: 'PUT',
    name: '编辑角色',
    approve_type: 8
  },
  ROLE_SET_LDAP_ROLE: {
    type: 'ROLE',
    url: '/role/setLdapUserDefRole',
    method: 'PUT',
    name: '设置默认LDAP角色',
    approve_type: 9
  }
}

export const LOG_TYPE_MAP = {
  1: '文件管理',
  2: '用户管理',
  3: '权限管理',
  4: '作业管理',
  5: '计算应用',
  6: '集群管理',
  7: '许可证管理',
  8: '项目管理',
  9: '3D云应用',
  10: '安全审批'
}
