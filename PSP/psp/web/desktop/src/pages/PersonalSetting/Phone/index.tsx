/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { Modal, Icon } from '@/components'
import { currentUser } from '@/domain'
import { Editor } from './Editor'
import { StyledItem } from '../style'
import { checkCaptcha } from '@/components'

export const Phone = observer(function Phone() {
  async function edit() {
    const token = await checkCaptcha()
    await Modal.show({
      title: '修改手机',
      content: ({ onOk, onCancel }) => (
        <Editor token={token} onCancel={onCancel} onOk={onOk} />
      ),
      footer: null,
    })
  }

  return (
    <StyledItem>
      <label>我的账号：</label>
      <span className='text'>{currentUser.mobile || '--'}</span>
      {/* <div className='right'>
        <Icon className='edit' type='rename' onClick={edit} />
      </div> */}
    </StyledItem>
  )
})
