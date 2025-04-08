import currentUser from './User'
import { lastMessages, recordList } from '.'

import sysConfig from './SysConfig'
export default async () => {
  await Promise.all([
    sysConfig.fetchWebsiteConfig(),
    sysConfig.fetchGlobalSysconfig(),
    currentUser.fetch(),
    sysConfig.fetchThreeMemberMgrConfig(),
    // sysConfig.fetchUserConfig()
    sysConfig.fetchJobWorkSpacePath()
  ]).then(() => {
    lastMessages.fetchLast()
    recordList.fetchLast()
  })

  const SystemPerm =
    currentUser?.perms?.system?.filter(p => p?.has === true) || ''
  localStorage.setItem('SystemPerm', JSON.stringify(SystemPerm))
  localStorage.setItem('userId', currentUser.id || '3McPKsZzSrp')
}
