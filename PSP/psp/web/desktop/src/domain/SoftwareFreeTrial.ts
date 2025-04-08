/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable, runInAction } from 'mobx'

export enum TrialApplyStatus {
  HAVE_NOT_APPLIED = 0, // 未申请
  APPLY_APPROVED, // 申请了，并至少有一个申请已通过
  APPLY_REJECTED // 申请了，没有任何申请被通过
}

class BaseSoftwareFreeTrial {
  @observable status: TrialApplyStatus = TrialApplyStatus.HAVE_NOT_APPLIED
}

export class SoftwareFreeTrial extends BaseSoftwareFreeTrial {
  constructor(props?: Partial<BaseSoftwareFreeTrial>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props?: Partial<BaseSoftwareFreeTrial>) {
    Object.assign(this, props)
  }

  @action
  async fetch() {
    const { data: status } = await softwareFreeTrialServer.getApplyStatus()
    runInAction(() => {
      this.update({
        status
      })
    })
  }
}
