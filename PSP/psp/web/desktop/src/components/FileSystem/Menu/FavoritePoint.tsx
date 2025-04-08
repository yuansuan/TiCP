import * as React from 'react'
import { Tree, message, Tooltip } from 'antd'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { Subject } from 'rxjs'

import { Http } from '@/utils'
import { FavoritePoint, PathHistory } from '@/domain/FileSystem'
import { EditableText } from '@/components'
import { Icon } from '@/components'
import TreeWrapper from '../TreeWrapper'

const DirectoryTree = Tree.DirectoryTree
const { TreeNode } = Tree

interface IProps {
  selectedKeys: string[]
  favoritePoint: FavoritePoint
  history: PathHistory
}

@observer
export default class Favorite extends React.Component<IProps> {
  @observable expandedKeys = []
  @action
  updateExpandedKeys = keys => (this.expandedKeys = keys)

  rename$ = new Subject()

  private onSelect = (selectedKeys, info) => {
    const { favoritePoint, history } = this.props

    // onSelect favoritePoint
    if (selectedKeys[0] === favoritePoint.favoriteId) {
      favoritePoint.fetch().then(data => {
        if (!data || data.length === 0) {
          message.warn('there are no favorites')
        } else {
          this.updateExpandedKeys([selectedKeys[0]])
        }
      })
    } else {
      const favoriteId = info.node.props.eventKey
      const node = favoritePoint.filterFirstNode(
        item => item.favoriteId === favoriteId
      )

      history.push({ source: favoritePoint, path: node.path })
    }
  }

  private onExpand = expandedKeys => {
    this.updateExpandedKeys(expandedKeys)
  }

  private onRename = ({ name, node }) => {
    if (!name || name === node.name) {
      return
    }
    Http.put('/file/favorite', {
      name,
      id: node.originId,
    }).then(() => {
      message.success('重命名成功')
      const targetNode = this.props.favoritePoint.filterFirstNode(
        item => item.id === node.id
      )
      if (targetNode) {
        targetNode.name = name
      }
    })
  }

  private renderTreeMenu = root => {
    let files = root && root.children

    return files.map(node => (
      <TreeNode
        title={
          <EditableText
            EditIcon={
              <Tooltip title='重命名'>
                <Icon type='edit-filled' />
              </Tooltip>
            }
            defaultValue={node.name}
            defaultShowEdit={false}
            onConfirm={name => this.onRename({ name, node })}
          />
        }
        key={node.favoriteId}>
        {this.renderTreeMenu(node)}
      </TreeNode>
    ))
  }

  render() {
    const { selectedKeys, favoritePoint } = this.props

    return (
      <TreeWrapper>
        <DirectoryTree
          defaultExpandAll
          selectedKeys={selectedKeys}
          expandedKeys={this.expandedKeys}
          onSelect={this.onSelect}
          onExpand={this.onExpand}
          expandAction='doubleClick'>
          {favoritePoint ? (
            <TreeNode title={favoritePoint.name} key={favoritePoint.favoriteId}>
              {this.renderTreeMenu(favoritePoint)}
            </TreeNode>
          ) : null}
        </DirectoryTree>
      </TreeWrapper>
    )
  }
}
