/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useCallback } from 'react'
import styled from 'styled-components'
import { Input } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from './store'
import { Button } from '@/components'
import debounce from 'lodash/debounce'

const StyledLayout = styled.div`
  display: flex;
  background-color: #ffffff;
  justify-content: space-between;
  align-items: center;
  .left {
    > .top {
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      .name {
        width: 100px;
      }
    }
    .bottom {
      margin: 16px 0;
      color: #666666;
    }
  }

  > .right {
    .bottom {
      margin-top: 10px;
      color: #666666;
    }
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()
  const state = useLocalStore(() => ({
    key: '',
    setKey(key) {
      this.key = key
    },
  }))

  const debounceKeyChange = useCallback(
    debounce(key => {
      store.setQueryKey(key)
    }, 300),
    []
  )

  return (
    <StyledLayout>
      <div className='left'>
        <div className='top'>
          <div className='name'>软件名称：</div>
          <Input
            value={state.key}
            onChange={e => {
              state.setKey(e.target.value)
              debounceKeyChange(state.key)
            }}
            placeholder='请输入关键字'
          />
        </div>
        <div className='bottom'>
          以下软件为自带许可证的软件，需配置好许可证后使用。
        </div>
      </div>
      <div className='right'>
        <Button
          onClick={e => {
            store.setQueryKey('')
            state.setKey('')
          }}>
          重置
        </Button>
        <div className='bottom'>共计{store.model.list.length}个软件</div>
      </div>
    </StyledLayout>
  )
})
