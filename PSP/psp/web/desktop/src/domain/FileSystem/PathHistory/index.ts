// FileManagement domain state model
import { computed } from 'mobx'

import Path from './Path'
import Store from '../Store'
import { Point, FavoritePoint } from '../Points'
import History from './History'

export default class PathHistory extends History {
  constructor() {
    super()

    Store.hooks.afterUpdate.tapAsync(
      'Store.rename: sync path history',
      this.onStoreUpdate
    )

    Store.hooks.afterDelete.tapAsync(
      'Store.delete: sync path history',
      this.onStoreDelete
    )
  }

  private onStoreDelete = paths => {
    paths.forEach(path => {
      this.list.forEach((item, index) => {
        if (item.path.startsWith(path)) {
          this.delete(index)
        }
      })
    })
  }

  private onStoreUpdate = nodes => {
    nodes.forEach(({ path, oldProps, newProps }) => {
      // rename
      if (oldProps.name !== newProps.name) {
        this.list.forEach(item => {
          if (item.path.startsWith(path)) {
            item.updatePath(`${path.replace(/[^\\/]+$/, '')}${newProps.name}`)
          }
        })
      }
    })
  }

  @computed
  get currentPath() {
    if (this.current) {
      return this.current.path
    } else {
      return ''
    }
  }

  push({ source, path }: { source: Point | FavoritePoint; path: string }) {
    // ignore duplicated history
    const { current } = this
    if (current) {
      if (current.source === source && current.path === path) {
        return null
      }
    }

    const item = new Path({
      source,
      path,
    })
    super.push(item)

    return item
  }
}
