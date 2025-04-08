// Copyright (C) 2023 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import List from './List'
import { Wrapper } from './style'

interface IProps {
  entry?: 'desktop' | 'admin'
  refresh?: boolean
}

@observer
export default class ProjectMG extends React.Component<IProps> {
  @observable isAdmin = this.props.entry ? this.props.entry === 'admin' : true

  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Wrapper>
        <List
          isRefresh={this.props.refresh}
          isAdmin={this.isAdmin}
        />
      </Wrapper>
    )
  }
}
