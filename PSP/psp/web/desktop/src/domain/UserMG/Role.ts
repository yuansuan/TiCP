import { action, observable } from 'mobx'

import { Http } from '@/utils'
import { RoleList, PermList } from '@/domain/UserMG'

export enum RoleType {
  CUSTOM,
  ROLE_ADMIN,
  ROLE_NORMAL_USER,
  ROLE_COMPANY_ADMIN,
  ROLE_SECURITY_ADMIN,
  ROLE_AUDIT_ADMIN
}

export const getRoleLevel = roleIds => {
  if (roleIds.includes(RoleType.ROLE_ADMIN)) {
    return RoleType.ROLE_ADMIN
  } else if (roleIds.includes(RoleType.ROLE_SECURITY_ADMIN)) {
    return RoleType.ROLE_SECURITY_ADMIN
  } else if (roleIds.includes(RoleType.ROLE_AUDIT_ADMIN)) {
    return RoleType.ROLE_AUDIT_ADMIN
  } else {
    return RoleType.ROLE_NORMAL_USER
  }
}

export const isRoleConflict = roleIds => {
  return (
    (roleIds.includes(RoleType.ROLE_ADMIN) &&
      roleIds.includes(RoleType.ROLE_SECURITY_ADMIN) &&
      roleIds.includes(RoleType.ROLE_AUDIT_ADMIN)) ||
    (roleIds.includes(RoleType.ROLE_ADMIN) &&
      roleIds.includes(RoleType.ROLE_SECURITY_ADMIN)) ||
    (roleIds.includes(RoleType.ROLE_ADMIN) &&
      roleIds.includes(RoleType.ROLE_AUDIT_ADMIN)) ||
    (roleIds.includes(RoleType.ROLE_SECURITY_ADMIN) &&
      roleIds.includes(RoleType.ROLE_AUDIT_ADMIN))
  )
}

export enum RoleTypeName {
  CUSTOM = '用户自定义',
  INTERNAL = '内置'
}

export interface IRequest {
  id: number
  name?: string
  comment?: string
  type?: RoleType
  typeName?: RoleTypeName
  perm?: Object
  is_internal: Boolean
  is_default: Boolean
  has_perm?: any
  permIds?: number[]
}

interface IRole {
  id: number
  name: string
  comment: string
  type: RoleType
  is_internal: Boolean
  is_default: Boolean
  typeName: RoleTypeName
  permList: PermList
  permIds: number[]
  perm: any
  has_perm: []
}

export default class Role implements IRole {
  public id = -1
  @observable public name = ''
  @observable public comment = ''
  @observable public type
  @observable public typeName
  @observable public permList = new PermList({})
  @observable public permIds = []
  @observable public perm
  @observable public is_internal
  @observable public is_default
  @observable public has_perm: []

  @observable public permNames = []

  constructor(props?: Partial<IRequest>) {
    props && this.init(props)
  }

  @action
  public init = (props: Partial<IRequest>) => {
    Object.assign(this, {
      id: props?.id,
      name: props?.name,
      type: props?.type,
      isInternal: props?.is_internal,
      isDefault: props?.is_default,
      // typeName: props.is_internal ? '内置' : '自定义',
      comment: props?.comment
    })
  }

  @action
  public initPerm = (props: Partial<IRequest>) => {
    const permIds = props?.has_perm?.map(item => item.id) || []
    Object.assign(this, {
      permIds: permIds,
      permList: new PermList(props.perm || {})
    })
  }

  @action
  public fetch = () =>
    Http.get(`/role/detail`, { params: { id: this.id } }).then(res => {
      this.init(res.data.role)
      this.initPerm(res.data)
      return res
    })

  @action
  public update = () =>
    Http.put(`/role/update`, {
      id: this.id,
      name: this.name,
      comment: this.comment,
      perms: this.permIds,
      // permNames: this.permNames,
      type: this.type
    }).then(res => {
      RoleList.fetch()
      return res
    })

  public toRequest = (): IRequest => ({
    id: this.id,
    name: this.name,
    comment: this.comment,
    type: this.type,
    typeName: this.typeName,
    perm: this.perm,
    is_internal: this.is_internal,
    is_default: this.is_default,
    has_perm: this.has_perm,
    permIds: this.permIds
  })
}
