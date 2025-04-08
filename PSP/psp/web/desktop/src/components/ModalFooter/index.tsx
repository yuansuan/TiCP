import * as React from 'react'
import { Button } from '@/components'

import { StyledModalFooter } from './style'

interface IProps {
  onOk?: () => void
  onCancel?: () => void
  okText?: string
  cancelText?: string
  className?: string
  cancelButtonProps?: any
  okButtonProps?: any
}

export default class ModalFooter extends React.Component<IProps> {
  onCancel = () => {
    const { onCancel } = this.props

    onCancel && onCancel()
  }

  onOk = () => {
    const { onOk } = this.props

    onOk && onOk()
  }

  render() {
    const {
      okText,
      cancelText,
      className,
      cancelButtonProps = {},
      okButtonProps = {},
    } = this.props
    const { onCancel, onOk } = this

    return (
      <StyledModalFooter className={className}>
        <div className='main'>
          <Button type='primary' {...okButtonProps} onClick={onOk}>
            {okText || '确认'}
          </Button>
          <Button onClick={onCancel} {...cancelButtonProps}>
            {cancelText || '取消'}
          </Button>
        </div>
      </StyledModalFooter>
    )
  }
}
