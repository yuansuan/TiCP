// Copyright (C) 2019 LambdaCal Inc.

import { observer } from 'mobx-react'
import { observable } from 'mobx'
import * as React from 'react'
import { nodeManager } from '@/domain/NodeMG'
import { history } from '@/utils'
import { NodeTable } from '../common'
import { Wrapper } from './style'

@observer
export default class NodeList extends React.Component<any> {
  @observable loading = true
  componentDidMount() {
    nodeManager.getNodeList().finally(() => {
      this.loading = false
    })
  }

  onRowClick = ({ rowData }) => {
    history.push(`/node/${rowData.node_name}`)
  }

  render() {
    return (
      <Wrapper>
        <NodeTable
          data={nodeManager.nodeList}
          loading={this.loading}
          // onRowClick={this.onRowClick}
        />
      </Wrapper>
    )
  }
}
