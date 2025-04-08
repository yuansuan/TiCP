import { observable, action } from 'mobx'

export type IRequest = IHelpDoc

export interface IHelpDoc {
  type: string
  value: string
}

export default class HelpDoc implements IHelpDoc {
  @observable type
  @observable value = ''

  constructor(request?: IRequest) {
    request && this.init(request)
  }

  @action
  init = (request: IRequest) => {
    Object.assign(this, {
      type: request.type || 'url',
      value: request.value
    })
  }

  toRequest = (): IRequest => ({
    type: this.type,
    value: this.value
  })
}
