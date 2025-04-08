import CompanyList from './CompanyList'
import FundOperateList from './FundOperateList'
export {
  FundOperate,
  FundOperateStatus,
  FundOperateStatus_MAP,
  FundOperateType,
  FundOperateType_MAP
} from './FundOperate'
export {
  default as Company,
  CompanyStatus,
  COMPANY_STATUS_MAP
} from './Company'

export { default as CompanyList } from './CompanyList'
export { Account as CompanyAccount, AccountStatus } from './CompanyAccount'
export const companyList = new CompanyList()
export const fundOperateList = new FundOperateList()
