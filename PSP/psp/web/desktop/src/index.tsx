/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import React, { Suspense } from 'react'
import ReactDOM from 'react-dom'
import { Spin } from 'antd'
import zhCN from 'antd/lib/locale-provider/zh_CN'
import 'moment/locale/zh-cn'
import cookie from 'cookie'
import AppDesktop from '@/AppDesktop'
import Page500 from '@/pages/500'
import { init, beforeLogin } from '@/domain'
import { StyledRoot } from './style'
import './indexDesktop.css'
import { Modal, Drawer } from '@/components'
import { theme } from '@/constant'
import { Provider as YSProvider } from '@/components'
import { Provider } from 'react-redux'
import store from './reducers'
import { needLogin } from '@/utils'
import Login from '@/pages/Login'
Modal.theme = theme
Modal.configProviderProps = { locale: zhCN }
Drawer.Wrapper = YSProvider
const HACKRLOAD = 'HACKRLOAD'

// ignore trivial rejection
window.addEventListener('unhandledrejection', event => {
  if (!event.reason) {
    event.preventDefault()
  }
})

window.addEventListener(
  'keydown',
  function (e) {
    //可以判断是不是mac，如果是mac,ctrl变为花键
    //event.preventDefault() 方法阻止元素发生默认的行为。
    if (
      [83].includes(e.keyCode) &&
      (navigator.platform.match('Mac') ? e.metaKey : e.ctrlKey)
    ) {
      e.preventDefault()
      // Process event...
    }
  },
  false
)
const render = (Child: React.ReactNode) =>
  ReactDOM.render(
    <Suspense
      fallback={
        <div id='sus-fallback'>
          <h1>Loading</h1>
        </div>
      }>
      <Provider store={store}>{Child}</Provider>
    </Suspense>,
    document.querySelector('#root')
  )

render(
  <div
    style={{
      position: 'fixed',
      width: '100%',
      height: '100%',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    }}>
    <Spin />
  </div>
)
const cookies = cookie.parse(document.cookie)

const isLoginPage = !cookies['refresh_token'] || !cookies['access_token']

if (isLoginPage || needLogin) {
  beforeLogin()
    .then()
    .finally(() => {
      render(
        <StyledRoot id='styledRoot'>
          <YSProvider>
            <Login />
          </YSProvider>
        </StyledRoot>
      )
    })
} else {
  init()
    .then(() => {
      if (
        !window.localStorage.getItem(HACKRLOAD) ||
        !window.localStorage.getItem('GlobalConfig')
      ) {
        render(
          <div
            style={{
              position: 'fixed',
              width: '100%',
              height: '100%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
            桌面加载中 <Spin />
          </div>
        )
        window.location.reload()
        window.localStorage.setItem(HACKRLOAD, HACKRLOAD)
      }

      window.localStorage.getItem('SystemPerm') &&
        render(
          <StyledRoot id='styledRoot'>
            <YSProvider>
              <AppDesktop />
            </YSProvider>
          </StyledRoot>
        )
    })
    .catch(error => {
      const { status } = error.response || {}
      if (status !== 401) {
        render(<Page500 description={error.message} />)
      }
    })
}
