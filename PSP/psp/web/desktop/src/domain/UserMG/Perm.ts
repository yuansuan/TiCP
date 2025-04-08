import { action, observable } from 'mobx'

export interface IRequest {
  id: number
  name: string
  resource_id: number
  type: string
  has: boolean
}

interface IPermission {
  id: number
  name: string
  resourceId: number
  type: string
  has: boolean
}

export default class Perm implements IPermission {
  public readonly id
  @observable public name = ''
  @observable public resourceId
  @observable public type = ''
  @observable public has = false

  constructor(props?: Partial<IRequest>) {
    props && this.init(props)
  }

  @action
  public init = (props: Partial<IRequest>) => {
    Object.assign(this, {
      id: props.id,
      name: props.name,
      resourceId: props.resource_id,
      type: props.type,
      has: props.has,
    })
  }
}
