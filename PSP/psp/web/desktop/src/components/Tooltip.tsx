/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import Trigger, { TriggerProps } from 'rc-trigger'
import 'rc-trigger/assets/index.less'

const StyledLayout = styled.div`
  background: #ffffff;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.15);
  padding: 16px;
  display: inline-block;
`

type Props = {
  title: string
  popup?: React.ReactNode | (() => React.ReactNode)
} & Omit<TriggerProps, 'popup'>

const builtinPlacements = {
  left: {
    points: ['cr', 'cl'],
  },
  right: {
    points: ['cl', 'cr'],
  },
  top: {
    points: ['bc', 'tc'],
  },
  bottom: {
    points: ['tc', 'bc'],
  },
  topLeft: {
    points: ['bl', 'tl'],
  },
  topRight: {
    points: ['br', 'tr'],
  },
  bottomRight: {
    points: ['tr', 'br'],
  },
  bottomLeft: {
    points: ['tl', 'bl'],
  },
}

export const Tooltip = observer(function Tooltip({
  children,
  title,
  ...props
}: Props) {
  return (
    <Trigger
      action={['hover']}
      popup={<StyledLayout>{title}</StyledLayout>}
      popupPlacement='bottom'
      builtinPlacements={builtinPlacements}
      {...props}>
      {children}
    </Trigger>
  )
})
