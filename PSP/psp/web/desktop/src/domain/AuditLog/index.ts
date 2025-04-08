import { LogList } from './LogList'

const operateUserType = {
  normal: 1,
  admin: 2,
  security: 3
}

export const normalLogList = new LogList(operateUserType.normal)
export const adminLogList = new LogList(operateUserType.admin)
export const securityLogList = new LogList(operateUserType.security)

export { LogList } from './LogList'
