/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal } from '@/components'
import { Editor } from './Editor'
import { StyledItem } from '../style'

export const Password = () => {
  const edit = async () => {
    await Modal.show({
      title: '修改密码',
      content: ({ onOk, onCancel }) => (
        <Editor onOk={onOk} onCancel={onCancel} />
      ),
      footer: null
    })
  }

  return (
    <StyledItem>
      <label>密码：</label>
      <span className='text psd' onClick={edit}>
        修改密码
      </span>
    </StyledItem>
  )
}
