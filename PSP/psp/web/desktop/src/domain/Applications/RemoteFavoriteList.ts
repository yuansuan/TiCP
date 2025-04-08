import { observable, action } from 'mobx'

import { Http } from '@/utils'
import Favorite, { IRequest as IFavoriteRequest } from './Favorite'

interface IRequest {
  favorites: IFavoriteRequest[]
}

interface IFavoriteList {
  list: Favorite[]
}

export default class FavoriteList implements IFavoriteList {
  *[Symbol.iterator]() {
    yield* this.list.values()
  }

  @observable list = []

  @action
  init = (props: IRequest) => {
    Object.assign(this, {
      list: (props.favorites || []).map(item => new Favorite(item)),
    })
  }

  @action
  fetch = () =>
    Http.get('/remote_app/favorite').then(res => this.init(res.data))

  @action
  add = name =>
    Http.post('/remote_app/favorite/add', { name }).then(() => this.fetch())

  @action
  delete = name =>
    Http.delete('/remote_app/favorite/delete', {
      params: { name },
    }).then(() => this.fetch())
}
