/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Env } from './Env'
import { Messages } from './Messages'
import { CompanyList } from './CompanyList'
import { PermList } from './PermList'
import { NoticeList } from './NoticeList'
import { RecordList } from './RecordList'

export * from './socket'
export { InvitationList } from './InvitationList'
export { Messages } from './Messages'

import { InvitationList } from './InvitationList'
import ProjectList from './ProjectList'
import { Domain } from './domain'
import Uploader from '@/components/Uploader'
import ApplicationList from './_Application'
import { SCList } from './SCList'
import CloudAppList from './Visualization/CloudAppList'
import VisualTaskList from './Visualization/VisualTaskList'
import VisualConfig from './Visualization/VisualConfig'
import VisIBVConfig from './Vis/VisIBVConfig'
import { Tasks } from './Tasks'
import { Account } from './Account'
import { WebConfig } from './WebConfig'
import Vis from './Vis'
import { LicenseMgrList, LicenseUsageList } from './LicenseMgr'

export const messages = new Messages()
export { default as currentUser } from './User'
export const companyList = new CompanyList()
export const permList = new PermList()
export const noticeList = new NoticeList()
export const env = new Env()
// NOTE 上传前置接口配置url
export const uploader = new Uploader({
  preUploadUrl: '/storage/preUpload'
})
export const lastMessages = new Messages()
export const lastInvitations = new InvitationList()
export const recordList = new RecordList()
export const appList = new ApplicationList()
export const scList = new SCList()
export const projectList = new ProjectList()
export const domainList = new Domain()
export const cloudAppList = new CloudAppList()
export const visualTaskList = new VisualTaskList()
export const visualConfig = new VisualConfig()
export const tasks = new Tasks()
export const account = new Account()
export const webConfig = new WebConfig()
export { default as BoxHttp } from './Box/BoxHttp'
export { NewBoxHttp } from './Box/NewBoxHttp'
export { useResize } from './useResize'
export const vis = new Vis()
export const visIBVConfig = new VisIBVConfig()
export { default as beforeLogin } from './beforeLogin'
export { default as sysConfig } from './SysConfig'
export { default as init } from './init'
export { default as Auth } from './Auth'
export const lmList = new LicenseMgrList()
export const lmUsageList = new LicenseUsageList()

