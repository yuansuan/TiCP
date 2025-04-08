import { Http } from '@/utils'
import sysUserManager from '@/domain/SysUserMG/SysUserList'

export class Dashboard {
  getDashboardInfo(type, dates) {
    // type ClUSTER_INFO, RESOURCE_INFO, JOB_INFO, SOFTWARE_INFO, ONLINE_USERS, USER_JOB_INFO
    if (type === 'ClUSTER_INFO') {
      return Http.get(`/dashboard/clusterInfo?start=${dates[0]}&end=${dates[1]}`)
    } else if ( type === 'RESOURCE_INFO') {
      return Http.get(`/dashboard/resourceInfo?start=${dates[0]}&end=${dates[1]}`)
    } else if ( type === 'JOB_INFO') {
      return Http.get(`/dashboard/jobInfo?start=${dates[0]}&end=${dates[1]}`)
    } else if ( type === 'SOFTWARE_INFO') {
      return Http.get(`/job/appJobNum?start=${dates[0]}&end=${dates[1]}`)
    } else if (type === 'USER_JOB_INFO') {
      return Http.get(`/job/userJobNum?start=${dates[0]}&end=${dates[1]}`)
    } else if ( type === 'ONLINE_USERS') {
      return new Promise((reslove, reject) => {
        sysUserManager.getSysUserList().then(res => {
          reslove({
            // 保持结构一致
            data: res.data?.page?.total || 0
          })
        }).catch((e) => reject(e))
      })
    } else {
      return Promise.resolve(null)
    }
  }
}

export default new Dashboard()
