import { observable, action } from 'mobx'

export interface IRequest {
  from: string
  name: string
  state: string
  user_name: string
}

interface IFavorite {
  from: string
  name: string
  state: string
  userName: string
}

export default class Favorite implements IFavorite {
  @observable from = ''
  @observable name = ''
  @observable state = ''
  @observable userName = ''

  constructor(props: IRequest) {
    this.init(props)
  }

  @action
  init = (props: IRequest) => {
    Object.assign(this, {
      from: props.from,
      name: props.name,
      state: props.state,
      userName: props.user_name
    })
  }
}
