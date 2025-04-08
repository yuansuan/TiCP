/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import _get from 'lodash.get'
import _set from 'lodash.set'
import { Http } from '@/utils'
import axios from 'axios'
type IComputeTypes = {
  compute_type: string
  show_name: string
}
type IGlobalConfig = {
  enable_visual: boolean
  compute_types: IComputeTypes[]
}

class SysConfig {
  threeMemberMgrConfig = {
    state: false,
    defaultApprover: {
      id: '',
      name: ''
    }
  }
  websiteConfig: any = {}
  JobBurstConfig: any = {}
  userConfig: any = {}
  mailInfoConfig: any = {}
  mailServerConfig: any = {}
  menuConfig: any = {}
  customConfig: any = {}
  firewallConfig: any = {}
  jobWorkSpacePath: string = ''
  globalConfig: IGlobalConfig
  get csrfToken() {
    return this.userConfig._ || ''
  }

  get enableThreeMembers() {
    return this.userConfig.enableThreeMembers
  }

  getIframeMap() {
    let res = {}
    const root = {
      path: 'root',
      schedulerName: '',
      children: this.menuConfig.root
    }

    function helper(root) {
      if (!root) {
        return
      }

      if (root.type === 'iframe') {
        res[root.key] = root
      }

      root?.children?.forEach(child => {
        helper(child)
      })
    }

    helper(root)

    // <key, {path,}>
    return res
  }

  getPathLinks() {
    let paths = [],
      currPath = window.location.hash.replace('#', '').split('?')[0] || ''

    const root = {
      path: 'root',
      schedulerName: '',
      children: this.menuConfig.root
    }

    function findPath(root, paths) {
      if (!root) {
        return false
      }

      if (root.path == currPath) {
        paths.push(root)
        return true
      }

      paths.push(root)

      if (root?.children?.some(child => findPath(child, paths))) {
        return true
      }

      paths.pop()
    }

    findPath(root, paths)

    paths.shift()

    return paths
  }

  getPageHeader() {
    let path = window.location.hash.replace('#', '').split('?')[0] || '',
      schedulerName = this.schedulerName,
      res = ''

    const root = {
      path: 'root',
      schedulerName: '',
      children: this.menuConfig.root
    }

    function helper(root, schedulerName, path) {
      if (!root) {
        return
      }

      if (
        root.path === path &&
        (root.schedulerName === schedulerName || root.schedulerName === '')
      ) {
        res = root.name
        return
      }

      root?.children?.forEach(child => {
        helper(child, schedulerName, path)
      })
    }

    helper(root, schedulerName, path)

    return res
  }

  get schedulerName() {
    return this.userConfig.schedulerName || 'lsf'
  }

  get installType() {
    return this.userConfig.installType || 'psp'
  }

  async fetchThreeMemberMgrConfig() {
    try {
      const res = await Http.get('/approve/threePersonManagement')
      this.threeMemberMgrConfig.state = res?.data?.state ?? false
    } catch (e) {
      console.error('获取三员管理配置失败', e)
    }

    try {
      const res = await Http.get('/sysconfig/getThreePersonManagementConfig')
      this.threeMemberMgrConfig.defaultApprover = {
        id: res?.data?.def_safe_user_id ?? '',
        name: res?.data?.def_safe_user_name ?? ''
      }
    } catch (e) {
      console.error('获取三员管理配置失败', e)
    }
  }

  async updateThreeMemberMgrConfig({ id, name }) {
    const res = await Http.post('/sysconfig/setThreePersonManagementConfig', {
      def_safe_user_id: id,
      def_safe_user_name: name
    })
    this.threeMemberMgrConfig.defaultApprover = {
      id,
      name
    }
    return res
  }

  get enableThreeMemberMgr() {
    return this.threeMemberMgrConfig.state
  }

  async fetchGlobalSysconfig() {
    const res = await Http.get('/sysconfig/global')

    if (res.data) {
      this.globalConfig = res.data
      localStorage.setItem('GlobalConfig', JSON.stringify(res.data))
    }

    return res
  }

  async fetchJobBurstConfig() {
    const res = await Http.get('/sysconfig/getJobBurstConfig')
    return res
  }
  async updateJobBurstConfig(body) {
    const res = await Http.post('/sysconfig/setJobBurstConfig', {
      ...body
    })
    return res
  }
  async fetchJobConfig() {
    const res = await Http.get('/sysconfig/getJobConfig')
    return res
  }
  async updateJobConfig(body) {
    const res = await Http.post('/sysconfig/setJobConfig', {
      ...body
    })
    return res
  }
  fetchFirewallConfig() {
    return Http.get('/firewall').then(res => {
      this.firewallConfig = res.data
      return res
    })
  }

  updateFireWallConfig(level) {
    return Http.put(`/firewall/${level}`).then(res => {
      this.firewallConfig = { ...this.firewallConfig, level }
      return res
    })
  }

  async fetchWebsiteConfig() {
    const metas = document.getElementsByTagName('meta')
    const vender = metas?.['vender']?.content || ''
    const configUrl = vender ? `/config.${vender}.json` : `/config.json`

    await axios.get(configUrl).then(res => {
      this.websiteConfig = res.data
      document.title = res.data?.title

      var link =
        document.querySelector('link[rel=icon]') ||
        document.createElement('link')
      link['href'] = res.data?.favicon || '/favicon.svg'
      link['rel'] = 'shortcut icon'
      document.head.appendChild(link)
    })
  }
  async fetchMenuConfig() {
    const url = '/sysconfig/menuconfig'
    // let data = await this.cacheResponseData(url)
    // this.menuConfig = data
    return []
  }

  checkHomedir(path) {
    return Http.get(`/sysconfig/checkhomedir?path=${path}`)
  }

  fetchUserConfig() {
    return Http.get('/sysconfig/userconfig').then(res => {
      this.userConfig = res.data
      return res
    })
  }

  async updateUserConfig(body) {
    const res = await Http.put('/sysconfig/userconfig', body)
    // 更新成功 update 数据
    this.userConfig = {
      // disabledFunction, scheduler_name, show_verify_code, installType 属性，
      // 暂时不支持配置修改
      disabledFunction: this.userConfig.disabledFunction,
      schedulerName: this.userConfig.schedulerName,
      installType: this.userConfig.installType,
      show_verify_code: this.userConfig.show_verify_code,
      _: this.userConfig._,
      ...body
    }
    return res
  }

  async updateUserVCode(show_verify_code) {
    await this.fetchUserConfig()

    const { password, homedir } = this.userConfig

    await this.updateUserConfig({
      homedir,
      password,
      show_verify_code,
      type: 'show_verify_code'
    })
  }

  async updateUserHomeDir(homedir) {
    await this.fetchUserConfig()

    const { password, show_verify_code } = this.userConfig

    await this.updateUserConfig({
      homedir,
      password,
      show_verify_code,
      type: 'homedir'
    })
  }

  async updateUserPwdConf(password) {
    await this.fetchUserConfig()

    const { homedir, show_verify_code } = this.userConfig

    await this.updateUserConfig({
      homedir,
      password,
      show_verify_code,
      type: 'password'
    })
  }

  fetchJobWorkSpacePath() {
    return Http.get('job/workspace').then(res => {
      this.jobWorkSpacePath = res.data?.name
    })
  }

  async updateJobResumbmitConfig(config) {
    const body = {
      ...this.JobBurstConfig.job,
      resubmit: {
        use_identical_job_dir: config.use_identical_job_dir
      }
    }

    const res = await Http.put('/sysconfig/jobconfig', body)
    // 更新成功 update 数据
    this.JobBurstConfig = { job: { ...body } }
    return res
  }

  async updateJobActionConfig(config) {
    const body = {
      ...this.JobBurstConfig.job,
      action: config.action
    }

    const res = await Http.put('/sysconfig/jobconfig', body)
    // 更新成功 update 数据
    this.JobBurstConfig = { job: { ...body } }
    return res
  }

  async updateJobListConfig(config) {
    const body = {
      ...this.JobBurstConfig.job,
      list_default_filter: {
        states: config.states
      }
    }
    const res = await Http.put('/sysconfig/jobconfig', body)
    // 更新成功 update 数据
    this.JobBurstConfig = { job: { ...body } }
    return res
  }

  async fetchMailInfoConfig() {
    return await Http.get('sysconfig/getEmailConfig').then(res => {
      let obj = res.data
      this.mailInfoConfig = obj
      return obj
    })
  }

  async updateMailNotifation(notification) {
    await this.fetchMailInfoConfig()

    let obj = JSON.parse(JSON.stringify(notification))

    const res = await Http.post('/sysconfig/setEmailConfig', {
      notification: obj
    })

    this.mailInfoConfig = notification
    return res
  }
  async fetchMailServerConfig() {
    return await Http.get('sysconfig/globalEmail').then(res => {
      let obj = res.data
      this.mailServerConfig = obj
      return obj
    })
  }
  async updateMailServer(server) {
    await this.fetchMailServerConfig()

    let obj = JSON.parse(JSON.stringify(server))

    const res = await Http.post('/sysconfig/globalEmail', {
      email_config: obj
    })
    this.mailInfoConfig = obj
    return res
  }

  async sendMail() {
    return await Http.post('/sysconfig/email/testSend')
  }
}

export default new SysConfig()
