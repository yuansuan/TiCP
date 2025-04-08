/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { byolServer } from '@/server'
export class BaseLicense {
  @observable id: string
  @observable ip: string
  @observable license_port: string
  @observable extra_port: string
  @observable provider_port: string
}

export class License extends BaseLicense {
  constructor(props?: Partial<BaseLicense>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<BaseLicense>) {
    Object.assign(this, props)
  }

  fetch = async (ID?: string) => {
    const { data } = await byolServer.getlicense(ID)
    const license = data.license && JSON.parse(data.license)

    runInAction(() => {
      this.update({
        ...license,
        id: data.id,
      })
    })
  }
}
