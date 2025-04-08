/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

interface IProps {
  description?: string
}

export default class Page500 extends React.Component<IProps> {
  render() {
    return (
      <div
        style={{
          display: 'flex',
          height: '100vh',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
        }}>
        <div>
          <img src={require('@/assets/images/500.svg')} />
        </div>
        <div>{this.props.description}</div>
      </div>
    )
  }
}
