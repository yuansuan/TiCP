/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

type ServerProp = (params?: any) => Promise<any>
export type FileServer = {
  delete: ServerProp
  list: ServerProp
  move: ServerProp
  mkdir: ServerProp
  getContent: ServerProp
  read: ServerProp
  stat: ServerProp
  getFileUrl: ServerProp
  download: ServerProp
  upload: ServerProp
  rename: ServerProp
  linkToCommon: ServerProp
  getUserCompressStatus: ServerProp
  compress: ServerProp
}

export { boxServer } from './boxServer'
export { newBoxServer } from './newBoxServer'
export { jobServer } from './jobServer'
export { jobCenterServer } from './jobCenterServer'
export { userServer } from './userServer'
export { accountServer } from './accountServer'
export { companyServer } from './companyServer'
export { billServer } from './billServer'
export { visualServer } from './visualServer'
export { byolServer } from './byolServer'
export { billUserServer } from './billForUserServer'
export { dashboardServer } from './dashboardServer'
export { appServer } from './appServer'
export { departmentServer } from './departmentServer'
export { fileRecordServer } from './fileRecordServer'
export { standardJobMGTServer } from './standardJobMGTServer'
