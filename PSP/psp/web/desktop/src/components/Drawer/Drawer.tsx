/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import ReactDOM, { render } from 'react-dom'
import { isReactComponent } from '../utils'
import { Drawer } from 'antd'
import { DrawerProps } from 'antd/lib/drawer'

const StyledDiv = styled.div``

type Props = DrawerProps & {
  content:
    | string
    | React.ReactNode
    | React.ComponentType<{
        onCancel?: (data?: any) => void | Promise<any>
        onOk?: (data?: any) => void | Promise<any>
      }>
  onOk: (next?) => void | Promise<any>
  onCancel: (next?) => void | Promise<any>
}

function DrawerComponent({
  content,
  onOk,
  onCancel,
  maskClosable = false,
  ...rest
}: Partial<Props>) {
  const Wrapper = YSDrawer.Wrapper || (({ children }) => children)

  let finalContent
  if (isReactComponent(content)) {
    const Content = content as any
    finalContent = <Content onCancel={onCancel} onOk={onOk} />
  } else if (typeof content === 'string') {
    finalContent = <span style={{ wordBreak: 'break-all' }}>{content}</span>
  } else {
    finalContent = content
  }

  return (
    <StyledDiv>
      <Wrapper>
        <Drawer
          {...rest}
          onClose={onCancel}
          visible={true}
          maskClosable={maskClosable}>
          <Wrapper>{finalContent}</Wrapper>
        </Drawer>
      </Wrapper>
    </StyledDiv>
  )
}

export function showDrawer({
  onCancel = undefined,
  onOk = undefined,
  getContainer,
  ...rest
}: Partial<Props>): Promise<any> {
  return new Promise((resolve, reject) => {
    // create mount node
    let mountNode = document.body
    if (getContainer) {
      mountNode =
        typeof getContainer === 'string'
          ? document.querySelector(getContainer)
          : typeof getContainer === 'function'
          ? getContainer()
          : getContainer
    }
    const div = document.createElement('div')
    mountNode.appendChild(div)

    function close() {
      let unmountResult = ReactDOM.unmountComponentAtNode(div)

      if (unmountResult && div.parentNode) {
        div.parentNode.removeChild(div)
      }
    }

    const finalOnCancel = (data?) => {
      const callback = (data?) => {
        // fix warning: Can't perform a React state update on an unmounted component
        Promise.resolve().then(() => {
          reject(data)
          close()
        })
      }

      // proxy onCancel
      if (onCancel) {
        // support async function
        const promise = onCancel(callback)
        if (promise && promise.then) {
          promise.then(() => {
            callback()
          })
        }

        return promise
      } else {
        callback(data)
      }

      return null
    }

    const finalOnOk = (data?) => {
      const callback = (data?) => {
        // fix warning: Can't perform a React state update on an unmounted component
        Promise.resolve().then(() => {
          resolve(data)
          close()
        })
      }

      // proxy onOk
      if (onOk) {
        // support async function
        const promise = onOk(callback)
        if (promise && promise.then) {
          promise.then(() => {
            callback()
          })
        }

        return promise
      } else {
        callback(data)
      }

      return null
    }

    render(
      <DrawerComponent
        getContainer={getContainer}
        onOk={finalOnOk}
        onCancel={finalOnCancel}
        {...rest}
      />,
      div
    )
  })
}

type YSDrawerType = Props & {
  show: typeof showDrawer
  Wrapper: React.ElementType
}

const YSDrawer: any = Drawer
YSDrawer.show = showDrawer

export default YSDrawer as YSDrawerType
