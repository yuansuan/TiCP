/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { computed, observable, runInAction } from 'mobx'

import { Http } from '@/utils'
import RemoteApp from './RemoteApp'
import { sysConfig } from '@/domain'
import { INSTALL_TYPE } from '@/utils/const'

interface IAppList {
  list: Map<string, RemoteApp>
}

export default class RemoteAppList implements IAppList {
  @observable hasLikedList = false
  @observable list = new Map()

  get = name => this.list.get(name)

  @computed
  get isAIO() {
    return sysConfig.installType === INSTALL_TYPE.aio
  }

  fetchTemplates = (params?: { state?: string }) =>
    Http.get('/app/template/list', {
      params: {
        compute_type: 'cloud'
      }
    }).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data.apps?.map(item => [item.name, new RemoteApp(item)])
        )
      })

      return res
    })

  // type local | cloud | all
  fetch = async (type: 'local' | 'cloud' | 'all' = 'all') => {
    const defaultRes = {
      data: {
        apps: []
      }
    }

    const requests = {
      all: ['/app/list']
    }

    const likedRequest = [
      type === 'local' ? '/app/liked/job_submission' : '/app/liked/cloud'
    ]

    const allRequests = requests[type].concat(likedRequest)

    const [resNative, resRemote, resLiked] = await Promise.all(
      allRequests.map(async url => {
        return url ? Http.get(url) : Promise.resolve(defaultRes)
      })
    )

    const likedSet = new Set(resLiked.data.favorites.map(f => f.name))

    this.hasLikedList = likedSet.size !== 0

    runInAction(() => {
      this.list = new Map(
        resNative.data.apps
          .concat(
            resRemote.data.apps.sort((a, b) =>
              a.name.toLowerCase() > b.name.toLowerCase() ? 1 : -1
            )
          )
          .map(item => {
            item.isLiked = likedSet.has(item.name)
            return [item.name, new RemoteApp(item)]
          })
      )
    })
  }

  addLiked = (name, type) => {
    return Http.post('/app/liked', {
      name,
      type: type === 'local' ? 'job_submission' : 'cloud'
    })
  }

  deleteLiked = (name, type) => {
    return Http.delete(
      `/app/liked/${type === 'local' ? 'job_submission' : 'cloud'}/${name}`
    )
  }

  delete = name =>
    Http.delete('/app/template', {
      params: { name, compute_type: 'cloud' }
    }).then(res => {
      this.fetchTemplates()

      return res
    })

  add = ({
    baseName,
    baseState,
    icon,
    description,
    baseVersion,
    newType,
    newVersion
  }) => {
    return Http.post('/app/template', {
      new_type: newType,
      base_name: baseName,
      state: baseState,
      description: description,
      icon: icon,
      base_version: baseVersion,
      new_version: newVersion,
      compute_type: 'cloud'
    }).then(res => {
      this.fetchTemplates()

      return res
    })
  }

  // batch publish
  publish = (names: string[]) => {
    const apps = names.map(name => this.get(name))
    return Http.put('/app/template/publish', {
      names: apps.map(app => app.name),
      state: 'published',
      compute_type: 'cloud'
    }).then(res => {
      this.fetchTemplates()

      return res
    })
  }

  // batch unpublish
  unpublish = (names: string[]) => {
    const apps = names.map(name => this.get(name))
    return Http.put('/app/template/publish', {
      names: apps.map(app => app.name),
      state: 'unpublished',
      compute_type: 'cloud'
    }).then(res => {
      this.fetchTemplates()

      return res
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
