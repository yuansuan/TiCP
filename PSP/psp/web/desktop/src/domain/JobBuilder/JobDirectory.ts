/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { DirectoryFactory } from '@/utils/FileSystem'
import { JobFile } from './JobFile'
import { observable } from 'mobx'
import { computed } from 'mobx'

export type ChildNode = JobDirectory | JobFile

export class JobDirectory extends DirectoryFactory<JobDirectory, JobFile> {
  @observable realCommonPathPrefix = null

  constructor(props: any) {
    super(props)
    Object.assign(this, props)
  }

  get isRoot() {
    return !this.parent
  }

  ensureDir = (path: String, common_path_prefix = '') => {
    const paths = path.split('/').filter(s => !!s)
    let currentNode: JobDirectory = this
    paths.forEach(name => {
      let node = currentNode.children.find(
        item => !item.isFile && item.name === name
      ) as JobDirectory

      if (!node) {
        node = new JobDirectory({
          name,
          parent: this,
          realCommonPathPrefix: common_path_prefix
        })
        currentNode.push(node)
      }
      currentNode = node
    })

    return currentNode
  }

  @computed
  get path() {
    let paths = [this.name]
    let parent = this.parent
    while (parent) {
      paths.unshift(parent.name)
      parent = parent.parent
    }
    return paths.filter(Boolean).join('/')
  }
}
