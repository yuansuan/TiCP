/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable, runInAction } from 'mobx'
import { currentUser } from '@/domain'
import {Http} from '@/utils';
import { userServer } from '@/server'
class BaseEnv {
  @observable productId: string = ''
  @observable YSEnv: string = ''
  @observable isVisible: boolean = false
}

export class Env extends BaseEnv {
  @action
  async initEnv() {}

  @action
  logout = async () => {
    await Http.post('/auth/logout')
    // don't use history.push which will be intercepted
    location.reload()
    localStorage.removeItem('userId')
    localStorage.removeItem('SystemPerm')
    localStorage.removeItem('CURRENTROUTERPATH')
    localStorage.removeItem('GlobalConfig')
    document.cookie = 'access_token=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
    document.cookie = 'refresh_token=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
  }
  @action
  async init() {
    await Promise.all([
      // Http.get('/sysconfig/userconfig'),
      userServer.current().then(({ data }) => {
        runInAction(() => {
          currentUser.update(data?.info)
        })
      })
    ])

    // localStorage.setItem('userId', '1002') // 2001 1002
    // observer user change
    window.onstorage = async e => {
      if (e.key === 'userId') {
        window.location.replace('/')
      }
    }
  }
}
