/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, computed, action } from 'mobx'
export class IVisualTask {
  @observable id: number
  @observable status: number
  @observable template_name: string
  @observable user_id: string
  @observable start_time: Date
  @observable link: string
  @observable screenshot: string
}
export default class VisualTask extends IVisualTask {
  @computed get startTime() {
    return this.start_time ? new Date(this.start_time) : new Date()
  }
  constructor(props: IVisualTask) {
    super()
    Object.assign(this, props)
  }
  @action
  updateScreenShot(sc: string) {
    this.screenshot = sc
  }
}
