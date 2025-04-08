/**
 * @module FileList
 * visible file list
 */
import { observable, transaction } from 'mobx'

import { Tree } from '@/utils/Tree'
import { BaseFile, BaseDirectory } from '@/utils/FileSystem'
import Store from './Store'

export default class List extends Tree {
  @observable parentPath: string

  constructor() {
    super()

    Store.hooks.afterDelete.tapAsync(
      'Store.delete: sync FileList',
      this.onAfterDelete
    )

    Store.hooks.afterUpdate.tapAsync(
      'Store.update: sync FileList',
      this.onAfterUpdate
    )

    Store.hooks.afterAdd.tapAsync('Store.add: sync HomePoint', this.onAfterAdd)
  }

  private onAfterAdd = nodes => this.mount(nodes)

  private onAfterDelete = paths => {
    this.removeNodes(item => paths.includes(item.path))
  }

  private onAfterUpdate = nodes => {
    // use transaction to prevent frequent update
    transaction(() => {
      nodes.forEach(({ path, newProps }) => {
        if (this.isChild(path)) {
          const targetNode = this.filterFirstNode(item => item.path === path)
          targetNode && targetNode.update(newProps)
        }
      })
    })
  }

  // mount exist files
  mount = nodes => {
    this.add(
      nodes
        .filter(item => this.isChild(item.path))
        .map(props => {
          if (!props.isFile) {
            return new BaseDirectory(props)
          } else {
            return new BaseFile(props)
          }
        })
    )
  }

  isChild = path => {
    const tailReg = /[\\/][^\\/]+$/
    const parentPath = path.replace(tailReg, '')

    return parentPath === this.parentPath
  }

  update = (path: string) => {
    if (path !== this.parentPath) {
      // freshen file list
      transaction(() => {
        this.parentPath = path
        this.clear()
        this.mount([...Store.nodeMap.values()])
      })
    }
  }
}
