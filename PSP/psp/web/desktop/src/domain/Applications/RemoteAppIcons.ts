import { observable } from 'mobx'
import { Http } from '@/utils'

export class RemoteAppIcons {
  @observable list = new Map()

  get = name => this.list.get(name)

  fetchIcon = name => {
    const icon = this.get(name)
    if (icon) {
      return Promise.resolve(icon)
    }
    const url = '/remote_app/icon'
    return Http.get(url, {
      params: { name },
    }).then(res => {
      const icon = res.data.icon_base64
      this.list.set(name, icon)
      return icon
    })
  }

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}

export default new RemoteAppIcons()
