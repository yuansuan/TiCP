/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { Tooltip } from 'antd'
import { ButtonProps as AntButtonProps, ButtonType } from 'antd/es/button'
import { TooltipProps } from 'antd/es/tooltip'
import ButtonGroup from 'antd/es/button/button-group'
import Icon from '../Icon'
import Mask from '../Mask2'
import { StyledButton } from './style'
import { isElement } from '../utils'

export type ButtonProps = Omit<
  AntButtonProps,
  'type' | 'disabled' | 'onClick'
> & {
  type?: ButtonType | 'secondary' | 'cancel'
  disabled?: boolean | string
  toolTipProps?: Partial<TooltipProps>
  onClick?: (e: React.MouseEvent<HTMLElement, MouseEvent>) => any | Promise<any>
}

class YSButton extends React.Component<ButtonProps> {
  static Group: typeof ButtonGroup
}

class Button extends YSButton {
  state = {
    loading: false
  }

  onClick = async e => {
    const onClick = this.props?.onClick
    const promise = onClick && onClick(e)
    if (promise && promise.then) {
      try {
        this.setState({
          loading: true
        })
        await promise
      } finally {
        this.setState({
          loading: false
        })
      }
    }
  }

  render() {
    let props = { ...this.props }
    const {
      className,
      loading,
      icon,
      type,
      disabled,
      onClick,
      toolTipProps,
      ...rest
    } = props
    let finalLoading = loading === undefined ? this.state.loading : loading

    // custom icon
    let YSIcon = null
    if (icon) {
      if (typeof icon === 'string') {
        YSIcon = <Icon type={icon} />
      } else if (isElement(icon)) {
        YSIcon = icon
      }
    }

    // custom button type
    const classNames = className ? [className] : []
    let finalType: ButtonType
    if (type === 'secondary') {
      classNames.push('ant-btn-secondary')
    } else if (type === 'cancel') {
      classNames.push('ant-btn-cancel')
    } else {
      finalType = type
    }

    if (finalLoading) {
      classNames.push('loading')
    }

    // custom disabled
    const finalDisabled = typeof disabled === 'string' || !!disabled

    const finalButton = (
      <StyledButton
        type={finalType}
        disabled={finalDisabled}
        className={classNames.join(' ')}
        onClick={this.onClick}
        {...rest}>
        {finalLoading && (
          <Mask.Spin
            style={{
              borderRadius: 'inherit'
            }}
            spinProps={{
              size: 'small'
            }}
          />
        )}
        <div className='container'>
          {YSIcon}
          {this.props.children && <span>{this.props.children}</span>}
        </div>
      </StyledButton>
    )

    return typeof disabled === 'string' ? (
      <Tooltip title={disabled} {...toolTipProps}>
        {finalButton}
      </Tooltip>
    ) : (
      finalButton
    )
  }
}

Button.Group = ButtonGroup

export default Button
