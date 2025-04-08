/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import { Modal, ConfigProvider } from 'antd'
import { ConfigProviderProps } from 'antd/es/config-provider'
import ReactDOM, { render } from 'react-dom'
import { ThemeProvider } from 'styled-components'
import { ButtonProps } from 'antd/lib/button'
import { ModalProps } from 'antd/lib/modal'
import { isReactComponent } from '../utils'
import ModalFooter from './Footer'
import { ConfirmContent } from './ConfirmContent'
import { Toolbar } from './Toolbar'
import { GlobalStyle } from './style'
import { theme as defaultTheme } from '@/utils'

export type IModalProps = ModalProps & {
  title: string
  toolbar: any
  content:
    | string
    | React.ReactNode
    | React.ComponentType<{
        onCancel?: (data?: any) => void | Promise<any>
        onOk?: (data?: any) => void | Promise<any>
      }>
  onOk: (next?) => void | Promise<any>
  onCancel: (next?) => void | Promise<any>
  footer: React.ComponentType<{
    onCancel?: (data?: any) => void | Promise<any>
    onOk?: (data?: any) => void | Promise<any>
  }>
  bodyStyle: React.CSSProperties
  okButtonProps: ButtonProps
  centered: boolean
  cancelButtonProps: ButtonProps
  className: string
  CancelButton: React.ComponentType<{
    onCancel?: (event?: React.MouseEvent<HTMLElement, MouseEvent>) => void
    loading?: boolean
  }>
  OkButton: React.ComponentType<{
    onOk?: (event?: React.MouseEvent<HTMLElement, MouseEvent>) => void
    loading?: boolean
  }>
  showHeader: boolean
}

export const showConfirm = (props: Partial<IModalProps> = {}): Promise<any> => {
  return showModal({
    ...props,
    className: 'confirm',
    title: '',
    footer: null,
    width: 433,
    bodyStyle: {
      padding: '0px 24px'
    },
    content: ({ onCancel, onOk }) => (
      <ConfirmContent {...props} onCancel={onCancel} onOk={onOk} />
    )
  })
}

// lazy render content to harmony animation
function LazyContent({ children }) {
  const [content, setContent] = useState(
    <div style={{ position: 'relative', height: 200 }}></div>
  )

  useEffect(() => {
    setTimeout(() => {
      setContent(children)
    }, 400)
  }, [])

  return content
}

export const showModal = ({
  title = '模态弹窗',
  content = '确认要执行此操作吗？',
  keyboard = false,
  onOk = undefined,
  onCancel = undefined,
  width = 520,
  footer = undefined,
  bodyStyle,
  centered = true,
  cancelButtonProps,
  okButtonProps,
  OkButton,
  CancelButton,
  toolbar,
  showHeader = true,
  getContainer,
  ...rest
}: Partial<IModalProps> = {}): Promise<any> => {
  return new Promise((resolve, reject) => {
    const Wrapper =
      YSModal.Wrapper ||
      (({ children }) => (
        <ConfigProvider {...YSModal.configProviderProps}>
          <ThemeProvider theme={YSModal.theme || defaultTheme}>
            {children}
          </ThemeProvider>
        </ConfigProvider>
      ))

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

    // close modal
    function close() {
      let unmountResult = ReactDOM.unmountComponentAtNode(div)

      if (unmountResult && div.parentNode) {
        div.parentNode.removeChild(div)
      }
    }

    let finalContent

    // React.ComponentType
    if (isReactComponent(content)) {
      const Content = content as any
      finalContent = <Content onCancel={finalOnCancel} onOk={finalOnOk} />
    } else if (typeof content === 'string') {
      finalContent = <span style={{ wordBreak: 'break-all' }}>{content}</span>
    } else {
      // null or React.ReactNode
      finalContent = content
    }

    // generate footer
    let finalFooter
    // use null to hide footer
    if (footer === null) {
      finalFooter = null
    } else if (isReactComponent(footer)) {
      // React.ComponentType
      const Footer = footer as any
      finalFooter = (
        <Wrapper>
          <Footer onCancel={finalOnCancel} onOk={finalOnOk} />
        </Wrapper>
      )
    } else {
      finalFooter = (
        <Wrapper>
          <ModalFooter
            OkButton={OkButton}
            CancelButton={CancelButton}
            cancelButtonProps={cancelButtonProps}
            okButtonProps={okButtonProps}
            onCancel={finalOnCancel}
            onOk={finalOnOk}
          />
        </Wrapper>
      )
    }

    render(
      <Wrapper>
        <YSModal
          visible={true}
          title={<span title={title}>{title}</span>}
          width={width}
          footer={finalFooter}
          bodyStyle={{ overflow: 'auto', ...bodyStyle }}
          maskClosable={false}
          onCancel={finalOnCancel}
          onOk={finalOnOk}
          centered={centered}
          keyboard={keyboard}
          getContainer={getContainer}
          {...rest}>
          <Wrapper>
            <GlobalStyle showHeader={showHeader} />
            <LazyContent>
              <>
                {toolbar && <Toolbar>{toolbar}</Toolbar>}
                {finalContent}
              </>
            </LazyContent>
          </Wrapper>
        </YSModal>
      </Wrapper>,
      div
    )
  })
}

type YSModalType = typeof Modal & {
  theme: any
  configProviderProps: ConfigProviderProps
  Wrapper: React.ElementType
  show: typeof showModal
  showConfirm: typeof showConfirm
  Footer: typeof ModalFooter
  Toolbar: typeof Toolbar
}

const YSModal: any = Modal
YSModal.showConfirm = showConfirm
YSModal.show = showModal
YSModal.Footer = ModalFooter
YSModal.Toolbar = Toolbar

export default YSModal as YSModalType
