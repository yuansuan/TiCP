import React from 'react'
import { Tree, message, Tooltip } from 'antd'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'
import { filter } from 'rxjs/operators'
import { Subject } from 'rxjs'

import { createMobxStream, formatRegExpStr } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { BaseDirectory } from '@/utils/FileSystem'
import { Validator } from '@/utils'
import { EditableText } from '@/components'
import { Icon } from '@/components'
import TreeWrapper from './TreeWrapper'

const DirectoryTree = Tree.DirectoryTree
const { TreeNode } = Tree

interface IProps {
  newDirectory$?: Subject<any>
  selectedKeys$?: Subject<string[]>
  expandedKeys$?: Subject<string[]>
  editable?: boolean
  points: any[]
  path?: string
  onContextMenu?: (args: any, point: any) => void
  expandedKeys?: string[]
  updateExpandedKeys?: (keys: string[]) => void
  onExpand?: (options: { keys: string[]; context: any; point: any }) => void
  selectedKeys?: string[]
  updateSelectedKeys?: (keys: string[]) => void
  onSelect?: (options: { keys: string[]; context: any; point: any }) => void
  disbaledKeys?: string[]
  hasPerm?: boolean
}

@observer
export default class TreeMenu extends React.Component<IProps> {
  @observable _expandedKeys = []
  @observable _selectedKeys = []
  @observable virtualNode = null
  @action
  _updateExpandedKeys = keys => (this._expandedKeys = keys)
  @action
  _updateSelectedKeys = keys => (this._selectedKeys = keys)
  @action
  setVirtualNode = node => (this.virtualNode = node)

  @computed
  get selectedKeys() {
    return this.props.selectedKeys || this._selectedKeys
  }

  @computed
  get expandedKeys() {
    return this.props.expandedKeys || this._expandedKeys
  }

  @computed
  get currentPath() {
    return this.selectedKeys[0]
  }

  @computed
  get selectedPoint() {
    if (!this.currentPath) {
      return null
    }

    const { points } = this.props
    const point = points.find(
      item =>
        item.rootPath === this.currentPath ||
        // windows/linux compatible
        new RegExp(
          `^${formatRegExpStr(item.rootPath).replace(/[\\/]$/, '')}[\\/]`
        ).test(this.currentPath)
    )

    return point
  }

  updateExpandedKeys = keys => {
    const { expandedKeys, updateExpandedKeys } = this.props

    if (expandedKeys) {
      updateExpandedKeys && updateExpandedKeys(keys)
    } else {
      this._updateExpandedKeys(keys)
    }
  }

  updateSelectedKeys = keys => {
    const { selectedKeys, updateSelectedKeys } = this.props

    if (selectedKeys) {
      selectedKeys && updateSelectedKeys(keys)
    } else {
      this._updateSelectedKeys(keys)
    }
  }

  componentDidMount() {
    const { newDirectory$, path } = this.props

    if (path) {
      this.updateSelectedKeys([path])
    }

    // new directory
    newDirectory$ &&
      newDirectory$.pipe(untilDestroyed(this)).subscribe(() => {
        if (!this.selectedPoint) {
          return
        }

        const parentNode = this.selectedPoint.filterFirstNode(
          item => item.path === this.currentPath
        )

        // create virtual node
        this.setVirtualNode({
          parentPath: parentNode.path,
          node: new BaseDirectory({
            name: '未命名文件夹'
          })
        })

        // expand parentNode
        let targetNode = parentNode
        const expandedKeys = [...this.expandedKeys, parentNode.path]
        while (targetNode && targetNode.parent) {
          expandedKeys.push(targetNode.parent.path)
          targetNode = targetNode.parent
        }
        this.updateExpandedKeys([...new Set(expandedKeys)])
      })

    // expandedKeys stream
    const expandedKeys$ = createMobxStream(() => this.expandedKeys).pipe(
      untilDestroyed(this)
    )
    // export expandedKeys
    this.props.expandedKeys$ &&
      expandedKeys$.subscribe(this.props.expandedKeys$)

    // selectedKeys stream
    const selectedKeys$ = createMobxStream(() => this.selectedKeys).pipe(
      untilDestroyed(this)
    )
    // export selectedKeys
    this.props.selectedKeys$ &&
      selectedKeys$.subscribe(this.props.selectedKeys$)
    // auto expand selected key
    selectedKeys$
      .pipe(filter((keys: []) => keys.length > 0))
      .subscribe(keys => {
        if (!this.selectedPoint) {
          return
        }

        const selectedKey = keys[0]
        const expandedKeys = [...this.expandedKeys]
        const { rootPath } = this.selectedPoint

        selectedKey
          .replace(new RegExp(`^${formatRegExpStr(rootPath)}[\\\\/]?`), '')
          .split(/[\\/]/)
          .reduce((path, seg) => {
            expandedKeys.push(path)
            return `${path}/${seg}`
          }, rootPath)

        this.updateExpandedKeys([...new Set(expandedKeys)])
      })
  }

  private onSelect = (selectedKeys, context) => {
    const { onSelect } = this.props
    const point = this.props.points.find(item => {
      const path = context.node.props.eventKey
      return (
        item.rootPath === path ||
        // windows/linux compatible
        new RegExp(
          `^${formatRegExpStr(item.rootPath).replace(/[\\/]$/, '')}[\\/]`
        ).test(path)
      )
    })
    onSelect && onSelect({ keys: selectedKeys, context, point })

    if (!this.props.selectedKeys) {
      this.updateSelectedKeys(selectedKeys)
    }
  }

  private onExpand = (expandedKeys, context) => {
    const { onExpand } = this.props
    const point = this.props.points.find(item => {
      const path = context.node.props.eventKey
      return (
        item.rootPath === path ||
        // windows/linux compatible
        new RegExp(
          `^${formatRegExpStr(item.rootPath).replace(/[\\/]$/, '')}[\\/]`
        ).test(path)
      )
    })

    // fetch files by path
    const { eventKey, expanded } = context.node.props
    if (!expanded) {
      point.service.fetch(eventKey)
    }

    onExpand && onExpand({ keys: expandedKeys, context, point })

    if (!this.props.expandedKeys) {
      this.updateExpandedKeys(expandedKeys)
    }
  }

  private beforeRename = ({ value, node }) => {
    const { error } = Validator.filename(value)
    if (error) {
      const prefix = node.path ? '重命名失败：' : '新建失败：'
      message.error(`${prefix}${error.message}`)
      return false
    }

    return true
  }

  private onRename = ({ name, node }) => {
    let point = null

    // new node
    if (!node.path) {
      point = this.selectedPoint
    } else {
      point = this.props.points.find(
        item =>
          item.path === node.path ||
          // windows/linux compatible
          new RegExp(
            `^${formatRegExpStr(item.path).replace(/[\\/]$/, '')}[\\/]`
          ).test(node.path)
      )
    }

    // new node with empty name
    // delete virtual item
    if (!node.path && !name) {
      this.setVirtualNode(null)
      return
    }

    // can't rename empty name
    // must rename different name
    if (node.path && (!name || name === node.name)) {
      return
    }

    const targetNode = point.filterFirstNode(item => item.path === node.path)
    if (node.path) {
      point.service
        .rename({
          newName: name,
          path: targetNode.path
        })
        .then(() => {
          message.success('重命名成功')
        })
    } else {
      point.service
        .createDir({
          path: `${this.virtualNode.parentPath}/${name}`
        })
        .then(() => {
          message.success('新建成功')
        })
        .finally(() => {
          // delete virtual item
          this.setVirtualNode(null)
        })
    }
  }

  private onCancelRename = ({ name, node }) => {
    if (!node.path) {
      // delete virtual item
      this.setVirtualNode(null)
    }
  }

  getIcon = node => {
    const expanded = this.expandedKeys.includes(node.path)

    if (node.is_sym_link) {
      if (node.isFile) {
        return <Icon type='link-file' />
      } else {
        if (expanded) {
          return <Icon type='link-folder-open' />
        } else {
          return <Icon type='link-folder' />
        }
      }
    } else {
      if (node.isFile) {
        return <Icon type='file' />
      } else {
        if (expanded) {
          return <Icon type='folder-open' />
        } else {
          return <Icon type='folder' />
        }
      }
    }
  }

  renderTreeMenu = root => {
    const { disbaledKeys = [], editable = true } = this.props
    let dirs = root && root.children.filter(item => !item.isFile)

    // hack: show tree switcher by default
    if (!dirs || dirs.length === 0) {
      dirs = [
        {
          _isMock: true
        }
      ]
    }

    // mount virtualNode
    const { virtualNode } = this
    if (virtualNode && virtualNode.parentPath === root.path) {
      dirs.unshift(virtualNode.node)
    }

    return dirs.map(node => {
      const disabled = disbaledKeys.includes(node.path)
      // hack: show tree switcher by default
      return node._isMock ? (
        <TreeNode key='mock' style={{ display: 'none' }} />
      ) : (
        <TreeNode
          title={
            // editable shouldn't disable new directory
            this.props.hasPerm && (editable || !node.path) ? (
              <EditableText
                Text={value => (
                  <span title={value}>{value.replace(/ /g, '\u00a0')}</span>
                )}
                style={{ display: 'inline-block' }}
                EditIcon={
                  <Tooltip title='重命名'>
                    <Icon type='edit-filled' />
                  </Tooltip>
                }
                defaultEditing={!node.path}
                defaultValue={node.name}
                defaultShowEdit={false}
                beforeConfirm={value => this.beforeRename({ value, node })}
                onConfirm={name => this.onRename({ name, node })}
                onCancel={name => this.onCancelRename({ name, node })}
              />
            ) : (
              <span>{node.name}</span>
            )
          }
          icon={this.getIcon(node)}
          disabled={disabled}
          key={node.path || 'virtual_node'}>
          {!disabled && this.renderTreeMenu(node)}
        </TreeNode>
      )
    })
  }

  render() {
    const { points, disbaledKeys = [] } = this.props
    const { selectedKeys, expandedKeys, onSelect, onExpand } = this

    return (
      <TreeWrapper>
        <DirectoryTree
          selectedKeys={selectedKeys}
          expandedKeys={expandedKeys}
          onSelect={onSelect}
          onExpand={onExpand}
          expandAction={false}>
          {points.map(point => {
            const expanded = expandedKeys.includes(point.path)
            return (
              <TreeNode
                title={point.name}
                key={point.path}
                icon={
                  expanded ? (
                    <Icon type='folder-open' />
                  ) : (
                    <Icon type='folder' />
                  )
                }
                disabled={disbaledKeys.includes(point.path)}>
                {this.renderTreeMenu(point)}
              </TreeNode>
            )
          })}
        </DirectoryTree>
      </TreeWrapper>
    )
  }
}
