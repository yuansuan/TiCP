/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Input } from 'antd'
import { InputProps } from 'antd/lib/input'

type IProps = InputProps & {
  validator?: RegExp | ((value: string) => boolean)
}

export default class ValidInput extends React.Component<IProps> {
  state = {
    value: '',
  }

  // fix chinese input issue
  isOnComposition = null
  emittedInput = null
  beforeCompositionValue = ''
  prevValue = ''

  inputRef = null

  private validate = value => {
    const { validator } = this.props

    if (validator) {
      if (validator instanceof RegExp) {
        return validator.test(value)
      } else {
        return validator(value)
      }
    }

    return true
  }

  private onChange = e => {
    let { value } = e.target
    const { onChange } = this.props

    if (!this.isOnComposition) {
      this.emittedInput = true

      // change valid value
      if (this.validate(value)) {
        if (onChange) {
          onChange(e)
        } else {
          this.setState({
            value,
          })
        }
      }
    } else {
      this.emittedInput = false

      // ignore onComposition validation
      if (onChange) {
        onChange(e)
      } else {
        this.setState({
          value,
        })
      }
    }

    this.prevValue = value
  }

  private onCompositionStart = e => {
    this.beforeCompositionValue = e.target.value

    this.isOnComposition = true
    this.emittedInput = false
  }

  private onCompositionEnd = e => {
    const { onChange } = this.props

    this.isOnComposition = false
    if (!this.emittedInput) {
      if (!this.validate(e.target.value)) {
        e.target.value = this.beforeCompositionValue
      }

      if (onChange) {
        onChange(e)
      } else {
        this.setState({
          value: e.target.value,
        })
      }
    }
  }

  render() {
    const { validator, value, ...rest } = this.props

    return (
      <Input
        ref={ref => (this.inputRef = ref)}
        {...rest}
        value={value || this.state.value}
        onChange={this.onChange}
        onCompositionStart={this.onCompositionStart}
        onCompositionEnd={this.onCompositionEnd}
      />
    )
  }
}
