/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { useStore } from './store'
import { Divider } from 'antd'
import { Export } from './Export'
import { formatAmount } from '@/utils'

const StyledLayout = styled.div`
  display: flex;
  margin: 0 20px;

  > div {
    margin-left: auto;
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()

  return (
    <StyledLayout>
      <div>
        共计{store.model.page_ctx.total}个作业，消费
        {formatAmount(store.totalAmount)}
        元
        <Divider type='vertical' />
        <Export />
      </div>
    </StyledLayout>
  )
})
