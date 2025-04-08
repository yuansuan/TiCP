/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { v4 as uuid } from 'uuid'
import { observable, action } from 'mobx'

export interface ITree {
  id: string
  isTree: boolean
  parent: Tree | null
  children: any[]
}

interface IOptions {
  depthFirst?: boolean
  self?: boolean
}

type MiniItem = { id: string; name: string; [key: string]: any }

type Filter = (item: { id: string; [key: string]: any }) => boolean

const defaultOptions = { depthFirst: true, self: true }

/**
 * 树形结构基类
 * @class Tree
 */
export class Tree implements ITree {
  readonly id = uuid()
  readonly isTree = true
  @observable parent = null
  @observable _children = []

  get children() {
    return this._children
  }

  set children(children) {
    this._children = children
  }

  /**
   * 清空
   * @method clear
   */
  @action
  clear() {
    this.children = []
  }

  /**
   * 根据输入节点，返回与输入节点重复的节点
   * 返回所有 name 重复的元素
   * @method getDuplicate
   */
  getDuplicate(item: MiniItem) {
    let duplicateNode = null

    this.children.every(child => {
      // check duplicate name
      if (child.name === item.name && item.id !== child.id) {
        duplicateNode = child
        return false
      }

      return true
    })

    return duplicateNode
  }

  /**
   * 在当前目录指定位置插入子文件/目录
   * @method insert
   */
  @action
  insert(item: MiniItem, index: number) {
    // 关联父级目录
    item.parent = this

    this.children.splice(index, 0, item)

    // check duplicate node
    const duplicateNode = this.getDuplicate(item)
    // remove duplicate node
    if (duplicateNode) {
      this.removeFirstNode(item => item.id === duplicateNode.id)
    }
  }

  /**
   * 在当前目录最后位置插入子文件/目录
   * @method push
   */
  @action
  push(item: MiniItem) {
    this.insert(item, this.children.length)
  }

  /**
   * 在当前目录开始位置插入子文件/目录
   * @method unshift
   */
  @action
  unshift(item: MiniItem) {
    this.insert(item, 0)
  }

  /**
   * 替换指定的文件/目录
   * @param {function} filter
   * @method replace
   */
  @action
  replace(filter: Filter, item: MiniItem) {
    const index = this.children.findIndex(filter)
    if (index < 0) {
      return
    }

    item.parent = this
    this.children.splice(index, 1, item)

    // check duplicate node
    const duplicateNode = this.getDuplicate(item)
    // remove duplicate node
    if (duplicateNode) {
      this.removeFirstNode(item => item.id === duplicateNode.id)
    }
  }

  /**
   * 根据指定规则过滤第一个节点
   * @method filterFirstNode
   * @param {function} filter
   * @param {Object} options
   */
  filterFirstNode(filter: Filter, options: IOptions = defaultOptions) {
    return this.tapFirstNode(filter, () => {}, options)
  }

  /**
   * 根据指定规则过滤所有子节点
   * @method filterNodes
   */
  filterNodes(filter: Filter, options: IOptions = defaultOptions) {
    return this.tapNodes(filter, () => {}, options)
  }

  /**
   * 移除指定规则的第一个节点
   * @method removeFirstNode
   */
  @action
  removeFirstNode(filter: Filter, options: IOptions = defaultOptions) {
    return this.tapFirstNode(
      filter,
      node => {
        if (node.parent) {
          const { children } = node.parent
          const index = children.findIndex(item => item.id === node.id)
          children.splice(index, 1)
        }
      },
      options
    )
  }

  /**
   * 移除指定规则的所有节点
   * @method removeNodes
   */
  @action
  removeNodes(filter: Filter, options: IOptions = defaultOptions) {
    return this.tapNodes(
      filter,
      node => {
        if (node.parent) {
          const { children } = node.parent
          const index = children.findIndex(item => item.id === node.id)
          children.splice(index, 1)
        }
      },
      options
    )
  }

  /**
   * 根据指定规则筛选出第一个节点，并进行指定操作
   * @method tapFirstNode
   */
  @action
  tapFirstNode(filter: Filter, operate, options: IOptions = defaultOptions) {
    const { depthFirst, self } = options
    let tappedNode = undefined

    // 遍历自身
    if (self) {
      if (filter(this)) {
        operate(this)
        return this
      }
    }

    // 深度优先遍历
    if (depthFirst) {
      this.children.every(item => {
        const flag = filter(item)
        // 节点筛选成功
        if (flag) {
          operate(item)
          tappedNode = item
          return false
        } else if (item.isTree) {
          // 遍历文件夹
          tappedNode = item.tapFirstNode(filter, operate, {
            self: false
          })
          // 查找到指定节点，则停止遍历
          if (tappedNode) {
            return false
          }
        }

        // 继续遍历
        return true
      })
    } else {
      // 广度优先遍历
      let queue = [...this.children]
      while (queue.length > 0) {
        const node = queue.pop()
        if (filter(node)) {
          operate(node)
          tappedNode = node
          break
        }
        if (node.isTree) {
          queue = queue.concat(node.children)
        }
      }
    }

    return tappedNode
  }

  /**
   * 根据指定规则筛选出所有节点，并进行指定操作
   * @method tapNodes
   * @param {function} filter
   * @param {function} operate
   * @param {Object} options
   */
  @action
  tapNodes(filter: Filter, operate, options: IOptions = defaultOptions) {
    const { depthFirst, self } = options
    let tappedNodes = []

    // 遍历自身
    if (self) {
      if (filter(this)) {
        tappedNodes.push(this)
        operate(this)
      }
    }

    // 深度优先遍历
    if (depthFirst) {
      this.children.forEach(item => {
        const flag = filter(item)
        // 节点筛选成功
        if (flag) {
          tappedNodes.push(item)
          operate(item)
        }

        // 遍历文件夹
        if (item.isTree) {
          tappedNodes = [
            ...tappedNodes,
            ...item.tapNodes(filter, operate, { self: false })
          ]
        }
      })
    } else {
      // 广度优先遍历
      let queue = [...this.children]
      while (queue.length > 0) {
        const node = queue.pop()
        if (filter(node)) {
          // you can return a flag in operate to stop the loop
          tappedNodes.push(node)
          if (operate(node)) {
            break
          }
        }
        if (node.isTree) {
          queue = queue.concat(node.children)
        }
      }
    }

    return tappedNodes
  }

  /**
   * 返回树的扁平结构
   */
  flatten() {
    let nodes = []
    this.children.forEach(node => {
      nodes = [...nodes, node]
      if (node.isTree) {
        nodes = [...nodes, ...node.flatten()]
      }
    })
    return nodes
  }
}
