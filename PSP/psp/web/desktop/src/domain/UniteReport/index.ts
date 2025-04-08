import { Http } from '@/utils'

export class UniteReport {
  async getReportByType(
    type,
    dates,
    license_id?: string,
    license_type?: string
  ) {
    if (
      type === 'MEM_UT_AVG' ||
      type === 'CPU_UT_AVG' ||
      type === 'TOTAL_IO_UT_AVG'
    ) {
      const { data } = await Http.get(
        `/report/resourceUtAvg?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data || {}
    } else if (type === 'CPU_TIME_SUM') {
      const { data } = await Http.get(
        `/report/cpuTimeSum?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'CPU_TIME_SUM') {
      const { data } = await Http.get(
        `/report/cpuTimeSum?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'JOB_COUNT') {
      const { data } = await Http.get(
        `/report/jobCount?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'JOB_DELIVER_COUNT') {
      const { data } = await Http.get(
        `/report/jobDeliverCount?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'DISK_UT_AVG') {
      const { data } = await Http.get(
        `/report/diskUtAvg?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'JOB_WAIT_STATISTIC') {
      const { data } = await Http.get(
        `/report/jobWaitStatistic?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'LICENSE_APP_USED_UT_AVG') {
      const { data } = await Http.get(
        `/report/licenseAppUsedUtAvg?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'LICENSE_APP_MODULE_USED_UT_AVG') {
      const { data } = await Http.get(
        `/report/licenseAppModuleUsedUtAvg?license_id=${license_id}&license_type=${encodeURIComponent(
          license_type
        )}&type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'NODE_DOWN_STATISTIC') {
      const { data } = await Http.get(
        `/report/nodeDownStatistics?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'VISUAL_USAGE_DURATION') {
      const { data } = await Http.get(
        `/vis/statistic/report/duration?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else if (type === 'VISUAL_NUMBER_STATUS') {
      const { data } = await Http.get(
        `/vis/statistic/report/numberStatus?type=${type}&start=${dates[0]}&end=${dates[1]}`
      )
      return data
    } else {
      return {}
    }
  }
}

export default new UniteReport()
