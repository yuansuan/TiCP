/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { ButtonProps } from 'antd/lib/button'

import Button from '../../Button'
import { StyledModalFooter } from './style'
import { isReactComponent } from '../../utils'

export interface FooterProps {
  onOk?: () => Promise<any> | void
  onCancel?: () => Promise<any> | void
  okText?: React.ReactNode
  cancelText?: React.ReactNode
  className?: string
  CancelButton?:
    | React.ComponentType<{
        onCancel?: (event?: React.MouseEvent<HTMLElement, MouseEvent>) => void
        loading?: boolean
      }>
    | React.ReactNode
  OkButton?:
    | React.ComponentType<{
        onOk?: (event?: React.MouseEvent<HTMLElement, MouseEvent>) => void
        loading?: boolean
      }>
    | React.ReactNode
  cancelButtonProps?: ButtonProps
  okButtonProps?: ButtonProps
}

export default class ModalFooter extends React.Component<FooterProps> {
  state = {
    okLoading: false,
    cancelLoading: false,
  }

  onCancel: (event: React.MouseEvent<HTMLElement, MouseEvent>) => void = () => {
    const { onCancel } = this.props

    if (onCancel) {
      const promise = onCancel()
      if (promise instanceof Promise) {
        this.setState({
          cancelLoading: true,
        })
        promise.finally(() => {
          this.setState({
            cancelLoading: false,
          })
        })
      }
    }
  }

  onOk: (event: React.MouseEvent<HTMLElement, MouseEvent>) => void = () => {
    const { onOk } = this.props

    if (onOk) {
      const promise = onOk()
      if (promise instanceof Promise) {
        this.setState({
          okLoading: true,
        })
        promise.finally(() => {
          this.setState({
            okLoading: false,
          })
        })
      }
    }
  }

  render() {
    const {
      okText,
      cancelText,
      className,
      CancelButton,
      cancelButtonProps = {},
      OkButton,
      okButtonProps = {},
    } = this.props
    const { onCancel, onOk } = this
    const { okLoading, cancelLoading } = this.state

    let finalOkButton
    // React.ComponentType
    if (isReactComponent(OkButton)) {
      const Button = OkButton as any
      finalOkButton = <Button loading={okLoading} onOk={onOk} />
    } else if (typeof OkButton === 'object') {
      // null or React.ReactNode
      finalOkButton = OkButton
    } else {
      finalOkButton = (
        <Button
          type='primary'
          loading={okLoading}
          {...okButtonProps}
          onClick={onOk}>
          {okText || '确认'}
        </Button>
      )
    }

    let finalCancelButton
    // React.ComponentType
    if (isReactComponent(CancelButton)) {
      const Button = CancelButton as any
      finalCancelButton = <Button loading={cancelLoading} onCancel={onCancel} />
    } else if (typeof CancelButton === 'object') {
      // null or React.ReactNode
      finalCancelButton = CancelButton
    } else {
      finalCancelButton = (
        <Button
          loading={cancelLoading}
          onClick={onCancel}
          {...cancelButtonProps}>
          {cancelText || '取消'}
        </Button>
      )
    }

    return (
      <StyledModalFooter className={className}>
        <div className='main'>
          {finalCancelButton}
          {finalOkButton}
        </div>
      </StyledModalFooter>
    )
  }
}
