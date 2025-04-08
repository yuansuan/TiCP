/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

export default class Page403 extends React.Component {
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
          <img src={require('@/assets/images/403.svg')} />
        </div>
        <div>您暂时没有权限访问。</div>
      </div>
    )
  }
}
