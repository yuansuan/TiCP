import { observable } from 'mobx'

class Paging {
  @observable index: number = 1
  @observable size: number = 10
}

class AppInfo {
  @observable app_id: string
  @observable app_name: string
  @observable duration: number
}

class ReportOverview {
  @observable list: AppInfo[] = []
  @observable total: number = 10
}

class AppDetail {
  @observable id: string
  @observable u_id: string
  @observable app_name: string
  @observable workstation_name: string
  @observable start_time:string
  @observable end_time: string
  @observable duration: number
}

class ReportDetail {
  @observable list: AppDetail[] = []
  @observable total: number = 10
}

export const overviewDataApp = new ReportOverview()

export const overviewPagingDataApp = new Paging()

export const detailDataApp = new ReportDetail()

export const detailPagingDataApp = new Paging()

export const overviewDataUser = new ReportOverview()

export const overviewPagingDataUser = new Paging()

export const detailDataUser = new ReportDetail()

export const detailPagingDataUser = new Paging()

export const DATE_FORMAT = "YYYY-MM-DD HH:mm:ss"