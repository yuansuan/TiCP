/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { FileTree } from './FileTree'
import { JobFile } from './JobFile'

interface FileListItem {
  is_dir?: boolean
  path?: string
  mod_time: number
  name: string
  size: number
}

export function fromJSON2Tree(list: FileListItem[] = []) {
  const root = new FileTree({
    name: '',
  })
  list.forEach(({ name, is_dir, ...rest }) => {
    if (is_dir) {
      root.ensureDir(name)
    } else {
      const names = name.split('/')
      const parentPath = names.slice(0, names.length - 1).join('/')
      const parent = parentPath ? root.ensureDir(parentPath) : root
      const file = new JobFile({ name: names[names.length - 1], ...rest })
      parent.push(file)
    }
  })
  return root
}

function toFileTree(node) {
  return new FileTree({
   ...node
  })
}

export function fromJSON2Tree2(list: FileListItem[] = [], rootPath: string, isTempDir) {
  const root = new FileTree({
    name: isTempDir ? '' : '.',
  })

  if (list.length === 0) {
    // root 已经有 '.' 了, rootPath 干掉开头的点
    root.ensureDir(rootPath.replace(/^\.\//, ''))
  } else {
    list.forEach(({ name, is_dir, path,  ...rest }) => {
      if (is_dir) {
        root.ensureDir(path)
      } else {
        const paths = path.split('/')
        const parentPath = paths.slice(0, paths.length - 1).join('/')
        const parent = parentPath ? root.ensureDir(parentPath) : root
        const file = new JobFile({ name, is_dir, ...rest })
        parent.push(file)
      }
    })
  }

  let currentNode = root.filterFirstNode(node => {
    return node.path === rootPath
  })

  return currentNode ? toFileTree(currentNode) : root
}


