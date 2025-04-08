/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { Slider } from 'antd'
import { SliderBaseProps } from 'antd/lib/slider'
import { observer } from 'mobx-react-lite'
import { CoreSelectorStyle } from './style'

interface IProps {
  steps?: number[]
  value: number
  props?: SliderBaseProps
  onChange: (value: number) => any
}

export const CoreSelector = observer(function CoreSelector(props: IProps) {
  const onChange = val => {
    const { onChange, steps } = props
    if (steps.includes(val)) {
      onChange(val)
      return
    }
    for (let i = 1; i < steps.length; i++) {
      if (steps[i] >= val && steps[i - 1] < val) {
        onChange(steps[i])
        return
      }
    }
  }

  const { props: ps, steps, value } = props

  const marks = useMemo(
    () =>
      steps.reduce((mark, num) => {
        mark[num] = num
        return mark
      }, {}),
    [steps]
  )

  const min = steps[0]
  const max = steps[steps.length - 1]

  return (
    <CoreSelectorStyle marks={marks}>
      <Slider
        {...ps}
        marks={marks}
        min={min}
        max={max}
        value={value}
        onChange={onChange}
      />
    </CoreSelectorStyle>
  )
})
