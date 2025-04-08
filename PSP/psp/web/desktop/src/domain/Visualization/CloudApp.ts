/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable } from 'mobx'
export class ICloudApp {
  @observable id: number
  @observable app_id: string
  @observable name: string
  @observable icon_data: string
  @observable app_param: string
  @observable app_param_paths: string
}
export default class CloudApp extends ICloudApp {
  constructor(props) {
    super()
    Object.assign(this, props)
    this.app_id = props?.ysid
  }
}
