/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

export default class Page404 extends React.Component {
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
          <img src={require('@/assets/images/404.svg')} />
        </div>
        <div>
          {/* 页面不存在，返回 <Link to='/'>首页</Link> 。 */}
          页面不存在 。
        </div>
      </div>
    )
  }
}
