import { observable, action } from 'mobx'

import { Http } from '@/utils'
import FavoriteDirectory from './FavoriteDirectory'
import Store from '../../Store'

export default class FavoritePoint extends FavoriteDirectory {
  @observable rootPath = ''

  constructor(props) {
    super(props)

    this.rootPath = props.path

    this.fetch()

    Store.hooks.afterDelete.tapAsync(
      'Store.delete: sync Point',
      this.onAfterDelete
    )
  }

  onAfterDelete = paths => {
    this.removeNodes(item => paths.includes(item.path))
  }

  fetch = () =>
    Http.get('/file/favorite').then((res: any) => {
      this.mount(this, res.data || [])

      return res.data
    })

  // mount fileList to specific node
  @action
  mount = (parentNode: FavoriteDirectory, files: any[] = []) => {
    // clear children
    parentNode.clear()
    files.forEach((item: any) => {
      const node = new FavoriteDirectory(item)
      parentNode.add(node)
    })
  }
}
