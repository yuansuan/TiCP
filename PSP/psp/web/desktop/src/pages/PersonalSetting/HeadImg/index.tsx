/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { UserOutlined } from '@ant-design/icons'
import { StyledItem } from '../style'
import { currentUser } from '@/domain'

export const HeadImg = observer(function HeadImg() {
  return (
    <StyledItem>
      <label>我的头像：</label>
      {(currentUser.headimg_url && (
        <img
          style={{ height: '32px', width: '32px', borderRadius: '50%' }}
          src={currentUser.headimg_url}
          alt='头像'
        />
      )) || <UserOutlined />}
    </StyledItem>
  )
})
