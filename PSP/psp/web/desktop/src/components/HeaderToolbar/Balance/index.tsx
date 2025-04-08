/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { env, account } from '@/domain'
import { Icon } from '@/components'
import { Tooltip } from 'antd'
import { Hover } from '@/components'
import { Http } from '@/utils'

const StyledLayout = styled.div`
  display: flex;
  align-items: center;

  > .balance {
    overflow: hidden;
    display: flex;
    align-items: center;
    padding: 0 10px;

    > .text {
      margin-left: 6px;
    }
  }

  > .recharge {
  }
`

export const Balance = observer(function Balance() {
  return (
    <StyledLayout>
      <div className='balance'>
        账户余额：
        <span className='text'>¥{account.account_balance}</span>
        {'\u00A0'} | {'\u00A0'} 授信额度：
        <span className='text'>¥{account.credit_quota_amount}</span>
      </div>
    </StyledLayout>
  )
})
