/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useCallback, useState } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { useStore } from './model'
import debounce from 'lodash/debounce'
import { Input, message } from 'antd'
import { Button, Modal } from '@/components'
import { EditingModal } from './DepartmentAction'

const StyledDiv = styled.div``

export const Toolbar = observer(function Toolbar() {
  const [key, setKey] = useState(null)
  const store = useStore()

  const debouncedSetQuery = useCallback(
    debounce(function (q) {
      store.list.setName(q.key)
    }, 300),
    []
  )

  const add = async () => {
    if (store.list.total >= 1000) {
      message.error('一个企业最多支持1000个部门')
      return
    }

    await Modal.show({
      title: '新增部门',
      content: ({ onCancel, onOk }) => (
        <EditingModal
          isAdding={true}
          onCancel={onCancel}
          onOk={onOk}
          refresh={store.fetch}
        />
      ),
      footer: null,
    })
  }

  return (
    <StyledDiv>
      <div className='header'>
        <h3>当前部门总数：{store.list.total} 个</h3>
      </div>
      <div className='toolbar'>
        <div className='left'>
          <Button type='primary' onClick={add}>
            新增部门
          </Button>
        </div>

        <div className='right'>
          <Input.Search
            allowClear
            placeholder='请输入部门名称'
            value={key}
            onChange={e => {
              setKey(e.target.value)
              debouncedSetQuery({
                key: e.target.value,
              })
            }}
          />
        </div>
      </div>
    </StyledDiv>
  )
})
