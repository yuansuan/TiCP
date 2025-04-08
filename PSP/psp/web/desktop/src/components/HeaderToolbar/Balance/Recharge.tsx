/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { env } from '@/domain'
import { Button } from '@/components'
import { buryPoint, history } from '@/utils'

export const Recharge = observer(function Recharge() {
  return (
    env.isPersonal && (
      <div style={{ padding: '0 10px', marginLeft: 8 }}>
        <Button
          type='secondary'
          size='small'
          onClick={() => {
            buryPoint({
              category: '导航栏',
              action: '充值',
            })
            history.push('/recharge')
          }}>
          充值
        </Button>
      </div>
    )
  )
})
