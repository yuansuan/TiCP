import CompanyMerchandiseList from './CompanyMerchandiseList'
export {
  CompanyMerchandise,
  CompanyMerchandiseStatus,
} from './CompanyMerchandise'
export const companyMerchandiseList = new CompanyMerchandiseList()

export enum CompanyMerchandiseStatus_MAP {
  STATE_UNKNOWN = '未知',
  STATE_ONLINE = '启用',
  STATE_OFFLINE = '未启用',
}
