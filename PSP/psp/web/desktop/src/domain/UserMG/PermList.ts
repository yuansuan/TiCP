/**
 * @module PermissionList
 * @description total permission list
 */

import { action, computed, observable } from 'mobx'

import { Perm } from '@/domain/UserMG'

export interface IRequest {}

export default class PermList {
  @observable perm = {}

  @computed
  get systemPerms() {
    return this.perm['system'] || []
  }

  @computed
  get subAppPerms() {
    return this.perm['local_app'] || []
  }

  @computed
  get remoteAppPerms() {
    return this.perm['visual_software'] || []
  }
  @computed
  get cloudAppPerms() {
    return this.perm['cloud_app'] || []
  }

  constructor(props: IRequest) {
    props && this.init(props)
  }

  @action
  init = (props: Partial<IRequest>) => {
    let p = {}

    Object.keys(props || {}).forEach(k => {
      p[k] = props[k]?.map(item => new Perm(item))
    })

    Object.assign(this, {
      perm: p
    })
  }
}
