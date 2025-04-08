import { observable } from 'mobx'

import { BaseDirectory } from '@/utils/FileSystem'

export default class FavoriteDirectory extends BaseDirectory {
  @observable favoriteId

  constructor(props) {
    super(props)

    // unique id
    this.favoriteId = window.btoa(
      encodeURIComponent(`${props.name}::${props.path}`)
    )
  }
}
