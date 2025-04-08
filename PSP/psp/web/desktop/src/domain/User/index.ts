import { Http } from '@/utils'
import { action, observable, computed } from 'mobx'
import { PermList } from '@/domain/UserMG'
import { Company } from '@/domain/CompanyMgr'
interface ItemPerm {
  id: string
  name: string
  key: string
  external_id: string
  has: boolean
}
interface IPerms {
  cloud_app: ItemPerm[]
  local_app: ItemPerm[]
  visual_software: ItemPerm[]
  system: ItemPerm[]
}
interface IUserConfig {
  ldap_enable: boolean
  openapi_switch: boolean
}
interface IRequest {
  id: number
  authType: string
  enabled: boolean
  name: string
  balance: number
  email: string
  mobile: string
  perms: IPerms
  config: IUserConfig
  permList: PermList
  roles: string[]
  users: string[]
  mountList: {
    id: string
    name: string
    path: string
  }[]
  user_id: string
  company: Company

  roleIds: number[]
}

interface ICurrentUser {
  id: number
  authType: string
  enabled: boolean
  name: string
  balance: number
  email: string
  mobile: string
  perms: IPerms
  config: IUserConfig
  permList: PermList
  roles: string[]
  mountList: {
    id: string
    name: string
    path: string
  }[]
  user_id: string
  company: Company
  roleIds: number[]
}

class CurrentUser implements ICurrentUser {
  readonly id
  @observable enabled = false
  @observable authType
  @observable name = 'yskj'
  @observable email = ''
  @observable mobile = ''
  @observable balance = 0
  @observable config: IUserConfig
  @observable perms: IPerms
  @observable permList = new PermList({})
  @observable mountList
  @observable roles = []
  @observable user_id
  @observable company: Company = new Company()

  @observable roleIds: number[] = []

  @computed
  get homePath() {
    const homePoint = this.mountList.find(item => item.id === 'home')
    return homePoint && homePoint.path
  }

  @action
  init = (props: Partial<IRequest>) => {
    Object.assign(this, {
      id: props.id,
      enabled: props.enabled,
      authType: props.authType,
      name: props.name,
      email: props.email,
      mobile: props.mobile,
      balance: props.balance,
      mountList: props.mountList,
      perms: props.perms,
      roles: props.roles,
      config: props.config,
      permList: new PermList(props.permList || {}),
      user_id: props.user_id || props.id,
      company: new Company(props.company),
      roleIds: props.roleIds
    })
  }

  @action
  fetch = async () => {
    return await Http.get('/user/current').then(res => {
      const SystemPerm = res.data.perm?.system?.filter(p => p?.has === true)
      localStorage.setItem('SystemPerm', JSON.stringify(SystemPerm))
      this.init({
        ...res?.data?.user_info,
        roles: res?.data?.role,
        config: res?.data?.conf,
        // ...res.data.auth,
        // mountList: res.data.mountList,
        perms: res.data.perm
        // groups: res.data.groups,
        // permList: res.data.permList,
        // company: res.data.company,
        // roleIds: res.data.roleIds
      })
    })
  }

  @action
  public updatePwd = (name, oldPwd, newPwd) => {
    return Http.put(`/user/ldap/pwd`, {
      name,
      password: oldPwd,
      new_password: newPwd
    })
  }

  public isPwdExpiredWarning = () => {
    return Http.get(`/user/ldap/pwdwarning/${this.id}`)
  }

  @action
  public update = ({
    name = '',
    password = '',
    email = this.email,
    mobile = this.mobile
  }) => {
    return Http.put(`/user/current`, {
      name,
      password,
      email,
      mobile,
      enabled: this.enabled
    }).then(res => {
      this.init({
        ...res.data.info,
        ...res.data.auth,
        mountList: res.data.mountList,
        perms: res.data.perms,
        roles: res.data.roles
      })
    })
  }

  get isLdapEnabled() {
    return this.config?.ldap_enable === false // 未开启ldap才可以创建用户
  }

  get isOpenapiSwitchEnable() {
    return this.config?.openapi_switch === true // 开启oepnapi才能给用户赋权openapi
  }

  get isPersonalJobManager() {
    const allJobManagerPerm = this.perms?.system?.filter(
      p => p.key.includes('job_manager') && p.has
    )
    if (allJobManagerPerm.length == 1) {
      return allJobManagerPerm[0].key === 'personal_job_manager'
    } else if (allJobManagerPerm.length === 2) {
      return false
    }

    return true
  }

  get hasSecurityApprovalPerm() {
    return this.perms?.system?.some(p => p.key === 'security_approval' && p.has)
  } 

  get hasNormalLogPerm() {
    return this.perms?.system?.some(p => p.key === 'normal_audit_log' && p.has)
  } 

  get hasSysAdminLogPerm() {
    return this.perms?.system?.some(p => p.key === 'system_admin_audit_log' && p.has)
  } 

  get hasSecurityAdminLogPerm() {
    return this.perms?.system?.some(p => p.key === 'security_admin_audit_log' && p.has)
  } 

  get hasSysMgrPerm() {
    return this.perms?.system?.some(p => p.key === 'sys_manager' && p.has)
  } 

  get hasProjectMgrPerm() {
    return this.perms?.system?.some(p => p.key === 'project_manager' && p.has)
  }

  get hasPersonalProjectMgrPerm() {
    return this.perms?.system?.some(
      p => p.key === 'personal_project_manager' && p.has
    )
  }

  get permKeys() {
    return this.perms.system.filter(p => p.has).map(p => p.key)
  }

  hasCreateWorkspacePerm = () => {
    return this.perms.includes('system-workspace_add')
  }

  hasMonitorPerm = () => {
    return this.perms.includes('internal-index_to_monitor')
  }

  get isSuperAdmin() {
    return this.roles.some((role) => role.id === 1)
  }

  get isSysAdmin() {
    return this.roles.includes('系统管理员')
  }

  get isOnlySecurityAdmin() {
    return this.roles.length === 1 && this.roles.includes('安全管理员')
  }

  get isOnlyAuditAdmin() {
    return this.roles.length === 1 && this.roles.includes('审计管理员')
  }

  get defaultVisitUrl() {
    let url = null,
      label = null

    if (this.hasMonitorPerm()) {
      url = '/dashboard'
      label = '集群监控'
    } else if (this.isOnlySecurityAdmin) {
      url = '/auditunapproved'
      label = '未审批'
    } else if (this.isOnlyAuditAdmin) {
      url = '/userlog'
      label = '普通用户日志'
    } else {
      url = '/localapps'
      label = '本地作业提交'
    }

    return {
      url,
      label
    }
  }
}

export default new CurrentUser()
