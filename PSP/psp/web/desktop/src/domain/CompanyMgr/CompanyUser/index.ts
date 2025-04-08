import UserList from './UserList'
export { User as CompanyUser } from './User'
export const companyUserList = new UserList()
export enum CompanyUserStatus_MAP {
  UNKNOWN = '未知',
  NORMAL = '正常',
  DELETED = '已删除',
}
