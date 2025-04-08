/**
 * @module FileStore
 * @description
 * 1. store all files and directories;
 * 2. dispatch curd events;
 */

import { AsyncParallelHook } from 'tapable'
import isEqual from 'lodash/isEqual'
import cloneDeep from 'lodash/cloneDeep'

export class FileStore {
  nodeMap: Map<string, any> = new Map()
  hooks

  constructor() {
    this.hooks = {
      afterAdd: new AsyncParallelHook(['nodes']),
      afterDelete: new AsyncParallelHook(['paths']),
      afterUpdate: new AsyncParallelHook(['nodes']),
    }
  }

  get(path: string) {
    return cloneDeep(this.nodeMap.get(path))
  }

  delete = (paths: string | string[]) => {
    let targets = []
    if (typeof paths === 'string') {
      targets = [paths]
    } else if (paths instanceof Array) {
      targets = paths
    }

    if (targets.length > 0) {
      const files = [...targets]
      targets.forEach(path => {
        this.nodeMap.delete(path)
      })

      // delete descendant files
      ; [...this.nodeMap.values()].forEach(node => {
        targets.forEach(path => {
          const windowsFormatPath = `${path.replace(/[\\/]$/, '')}\\`
          const linuxFormatPath = `${path.replace(/[\\/]$/, '')}/`
          // delete childNodes
          if (
            node.path.startsWith(windowsFormatPath) ||
            node.path.startsWith(linuxFormatPath)
          ) {
            files.push(node.path)
            this.nodeMap.delete(node.path)
          }
        })
      })

      // trigger afterDelete hook
      this.hooks.afterDelete.callAsync(files, () => {})
    }
  }

  harmony = props => {
    const harmonyProps = cloneDeep(props)

    harmonyProps.modifiedTime = props.m_date
    Reflect.deleteProperty(harmonyProps, 'm_date')

    harmonyProps.isFile = !props.is_dir
    Reflect.deleteProperty(harmonyProps, 'is_dir')

    return harmonyProps
  }

  update = (nodes = []) => {
    if (!nodes) {
      throw new Error('nodes must be array')
    }

    const newNodes = []
    const updateNodes = []

    nodes.forEach(node => {
      const oldStruct = this.get(node.path)
      const newStruct = Object.assign({}, oldStruct || {}, {
        ...this.harmony(node),
      })

      // update old struct
      if (oldStruct) {
        // ignore same struct
        if (isEqual(oldStruct, newStruct)) {
          return
        }
        this.nodeMap.set(newStruct.path, newStruct)
        updateNodes.push({
          path: newStruct.path,
          oldProps: cloneDeep(oldStruct),
          newProps: cloneDeep(newStruct),
        })
      } else {
        // add new node
        this.nodeMap.set(newStruct.path, newStruct)
        newNodes.push(newStruct)
      }
    })

    if (newNodes.length > 0) {
      // trigger afterAdd hook
      this.hooks.afterAdd.callAsync(cloneDeep(newNodes), () => {})
    }

    if (updateNodes.length > 0) {
      // trigger afterUpdate hook
      this.hooks.afterUpdate.callAsync(updateNodes, () => {})
    }
  }

  rename = (path: string, newName: string) => {
    let oldStruct = this.get(path)
    if (oldStruct) {
      const newPath = `${path.replace(/[\\/][^\\/]+$/, '')}/${newName}`
      const newStruct = Object.assign({}, oldStruct, {
        name: newName,
        path: newPath,
      })

      // generate updateInfo
      const windowsFormatPath = `${path.replace(/[\\/]$/, '')}\\`
      const linuxFormatPath = `${path.replace(/[\\/]$/, '')}/`
      const deletedNodes = [...this.nodeMap.values()].reduce((res, node) => {
        // delete childNodes
        if (
          node.path.startsWith(windowsFormatPath) ||
          node.path.startsWith(linuxFormatPath)
        ) {
          res.push(node.path)
          this.nodeMap.delete(node.path)
        }

        return res
      }, [])

      // trigger afterDelete hook
      if (deletedNodes.length > 0) {
        this.hooks.afterDelete.callAsync(deletedNodes, () => {})
      }

      // update
      this.nodeMap.delete(path)
      this.nodeMap.set(newStruct.path, newStruct)
      // trigger afterUpdate hook
      this.hooks.afterUpdate.callAsync(
        [
          {
            path,
            oldProps: oldStruct,
            newProps: cloneDeep(newStruct),
          },
        ],
        () => {}
      )
    }
  }
}

export default new FileStore()
