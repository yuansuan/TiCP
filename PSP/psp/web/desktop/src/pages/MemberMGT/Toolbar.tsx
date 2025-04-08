/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useCallback, useState } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from './model'
import { companyServer } from '@/server'
import { currentUser, env } from '@/domain'
import debounce from 'lodash/debounce'
import { Input, message } from 'antd'
import { Button, Modal } from '@/components'
import { InviteModal } from './InviteModal'
import { InviteResultModal } from './InviteResultModal'

const StyledDiv = styled.div``

export const Toolbar = observer(function Toolbar() {
  const store = useStore()

  const [key, setKey] = useState('')

  const debouncedSetQuery = useCallback(
    debounce(function (q) {
      store.setQuery(q)
    }, 300),
    []
  )

  const invite = async () => {
    const res = await Modal.show({
      title: '邀请成员',
      width: 600,
      bodyStyle: {
        height: 400
      },
      content: ({ onCancel, onOk }) => (
        <InviteModal
          onCancel={onCancel}
          onOk={onOk}
          departments={store.departmentList}
        />
      ),
      footer: null
    })

    await Modal.show({
      title: '邀请结果',
      width: 600,
      content: <InviteResultModal list={res.invite_result} />,
      CancelButton: ({ onCancel }) => (
        <Button
          onClick={() => {
            invite()
            onCancel()
          }}>
          继续邀请
        </Button>
      )
    })
  }

  return (
    <StyledDiv>
      <div className='header'>
        <h3>当前总人数：{store.total}人</h3>
      </div>
      <div className='toolbar'>
        <div className='left'>
          <Button type='primary' onClick={invite}>
            邀请
          </Button>
        </div>

        <div className='right'>
          <Input.Search
            allowClear
            placeholder='请输入成员姓名/手机号'
            value={key}
            onChange={e => {
              setKey(e.target.value)
              debouncedSetQuery({
                key: e.target.value
              })
            }}
          />
        </div>
      </div>
    </StyledDiv>
  )
})
