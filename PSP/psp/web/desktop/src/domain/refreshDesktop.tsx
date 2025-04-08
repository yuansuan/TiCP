import { appList, vis, sysConfig } from '@/domain'
import store from '@/reducers'

import {
  allApps,
  generateCalcApp,
  GenerateDesktopApps,
  generateCloudApp
} from '@/utils'

export const fetchDesktopAPP = async () => {
  await fetchSoftware()

  store.dispatch({ type: 'DESKHIDE' })

  setTimeout(() => {
    store.dispatch({ type: 'DESKSHOW' })
  }, 300)
}

// 获取最新的桌面软件
export const fetchSoftware = async (): Promise<any> => {
  const collectFetchSession = []
  try {
    await appList.fetch(true)
    // enable_visual才能调用vis 相关的服务
    if (sysConfig.globalConfig.enable_visual) {
      const res = await vis.getUsingSoftware()
      const newStatuses  = res.data?.using_statuses.filter(s => s.status === 'STARTED') || []
      collectFetchSession.push(newStatuses)
    }
    const SESSIONS = await Promise.all(collectFetchSession)
    // 已经开启的会话
    const { softWareType } = generateCloudApp(SESSIONS?.flat())
    const { generateAppInfo, appType } = generateCalcApp(
      appList.publishedAppList || []
    )
    const appsAll = [...generateAppInfo, ...allApps]
    const allType = [...appType.map(a => `${a.type}`), ...softWareType]
    store.dispatch({ type: 'apps', data: appsAll, payload: 'generateApp' })

    const generaApps = GenerateDesktopApps(appsAll, allType)
    store.dispatch({
      type: 'desktop',
      data: generaApps,
      payload: 'generateApp'
    })

    return generaApps
  } catch (error) {
    console.error('fetchSoftware error: ', error)
    return []
  }
}
