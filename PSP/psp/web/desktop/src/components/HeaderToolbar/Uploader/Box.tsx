/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Icon } from '@/components'
import { formatByte } from '@/utils/Validator'
import { env } from '@/domain'
import { buryPoint } from '@/utils'
import { Tooltip } from 'antd'
import { notification } from '@/components'

const StyledLayout = styled.div`
  display: flex;
  align-items: center;

  > .space {
    overflow: hidden;
    white-space: nowrap;
    margin-left: 4px;
  }
`

const statusColorMap = {
  notfound: '#F5222D',
  success: '#9B9B9B',
  error: '#F5222D',
  warning: '#faad14'
}

type BoxData = {
  disk: {
    free: number
    used: number
    total: number
  }
  system: {
    load: {
      1: number
      5: number
      15: number
    }
  }
  time: number
}

type Props = {
  dropdownVisible?: boolean
}

export const Box = observer(function Box({ dropdownVisible }: Props) {
  const state = useLocalStore(() => ({
    diskTotal: 0,
    diskUsed: 0,
    tipVisible: false,
    hovered: false,
    setHovered(flag) {
      this.hovered = flag
    },
    setTipVisible(visible) {
      this.tipVisible = visible
    },
    get displayDiskTotal() {
      return formatByte(this.diskTotal)
    },
    get displayDiskUsed() {
      return formatByte(this.diskUsed)
    },
    
  }))

  // 展开下拉菜单，隐藏 tip
  useEffect(() => {
    if (dropdownVisible) {
      state.setTipVisible(false)
    }
  }, [dropdownVisible])



  return (
    <Tooltip
      onVisibleChange={visible => state.setTipVisible(visible)}
      visible={state.tipVisible}
      title='存储空间'>
      <StyledLayout
        onMouseLeave={() => {
          state.setHovered(false)
        }}
        onMouseEnter={() => {
          state.setHovered(true)
        }}
        onClick={() => {
          buryPoint({
            category: '导航栏',
            action: '存储空间'
          })
        }}>
        <Icon
          type={
            dropdownVisible || state.hovered
              ? 'storage_active'
              : 'storage_default'
          }
          
        />
        <span className='space'>
          {state.displayDiskUsed} / {state.displayDiskTotal}
        </span>
      </StyledLayout>
    </Tooltip>
  )
})
