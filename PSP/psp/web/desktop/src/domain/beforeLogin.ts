import sysConfig from './SysConfig'

export default async () => {
  await Promise.all([sysConfig.fetchWebsiteConfig()])
}
