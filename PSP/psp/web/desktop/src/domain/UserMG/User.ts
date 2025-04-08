import { Http } from '@/utils'
import { action, computed, observable, runInAction } from 'mobx'

import { UserList, PermList } from '@/domain/UserMG'
import { Timestamp } from '@/utils'

import { approve_status_map } from './const'

interface IRequest {
  id: number
  name: string
  email: string
  enabled: boolean
  mobile: string
  perm: Object
  roles: number[]
  is_internal: boolean
  created_at: number
  approve_status: number
  enable_openapi: boolean
  openapi_certificate: string
}

interface IUser {
  id: number
  name: string
  enabled: boolean
  mobile: string
  email: string
  isInternal: boolean
  roles: number[]
  created_at: number
  approve_status: number
  enable_openapi: boolean
}

export default class User implements IUser {
  public readonly id
  @observable public name = ''
  @observable public password = ''
  @observable public email = ''
  @observable public enabled = true
  @observable public mobile = ''
  @observable public created_at = null
  @observable public isInternal = false
  @observable public roles = []
  @observable public roleNames: string[]
  @observable public enable_openapi = false
  @observable approve_status = -1
  @observable public openapi_certificate

  constructor(props?: Partial<IRequest>) {
    props && this.init(props)
  }

  @action
  public init = (props: Partial<IRequest>) => {
    Object.assign(this, {
      id: props.id,
      name: props.name,
      enabled: props.enabled,
      email: props.email,
      mobile: props.mobile,
      isInternal: props.is_internal,
      created_at: props.created_at,
      permList: new PermList(props.perm || {}),
      roles: props.roles || [],
      approve_status: props.approve_status,
      enable_openapi: props.enable_openapi,
      openapi_certificate: props.openapi_certificate,
    })
  }

  @computed
  get approve_status_str() {
    return approve_status_map.get(this.approve_status) || '--'
  }

  @action
  public fetch = () => {
    if (!this.id) return
    Http.get(`/user/get`, {
      params: {
        id: this.id
      }
    }).then(res => {
      if (res.data) {
        const newObj = {
          ...res.data.user_info,
          roles: res.data.role.map(u => u.id),
          perm: res.data.perm,
          openapi_certificate: res.data.openapi_certificate
        }
        this.init(newObj)
      }

      return res
    })
  }

  @action
  public update = () =>
    Http.put(`/user/update`, {
      id: this.id,
      name: this.name,
      roleNames: this.roleNames,
      roles: this.roles,
      email: this.email,
      mobile: this.mobile,
      enabled: this.enabled,
      enable_openapi: this.enable_openapi
    }).then(res => {
      runInAction(() => {
        UserList.fetch()
      })
      return res
    })

  @action
  public toRequest = (): IRequest => ({
    id: this.id,
    name: this.name,
    enabled: this.enabled,
    email: this.email,
    mobile: this.mobile,
    roles: this.roles,
    perm: {},
    created_at: this.created_at,
    is_internal: this.isInternal,
    approve_status: this.approve_status,
    enable_openapi: this.enable_openapi,
    openapi_certificate: this.openapi_certificate,
  })

  @action
  public genCert = () =>
    Http.put(`/user/genOpenapiCertificate`, {
      user_id: this.id,
    }).then(res => {
      this.openapi_certificate = res.data
  })
}
