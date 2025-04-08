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
  @observable app_name: string
  @observable workstation_name: string
  @observable start_time: string
  @observable end_time: string
  @observable duration: number
}

class ReportDetail {
  @observable list: AppDetail[] = []
  @observable total: number = 10
}

export const overviewData = new ReportOverview()

export const overviewPagingData = new Paging()

export const detailData = new ReportDetail()

export const detailPagingData = new Paging()

export const DATE_FORMAT = 'YYYY-MM-DD HH:mm:ss'
