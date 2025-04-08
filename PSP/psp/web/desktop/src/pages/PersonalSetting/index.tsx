/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Empty, message } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { StyledItem, StyledLayout } from './style'
import { Phone } from './Phone'
import { Password } from './Password'
import { Name } from './Name'
import { HeadImg } from './HeadImg'
import { currentUser, env } from '@/domain'
const PersonalSetting = observer(function PersonalSetting() {
  const store = useLocalStore(() => {
    return {
      role: '--',
      setRole(role) {
        this.role = role
      }
    }
  })

  return (
    <StyledLayout>
      {currentUser.id ? (
        <>
          {/* <HeadImg /> */}
          {/* <Phone /> */}
          <Name />
          <Password />
          <StyledItem>
            {/* <label>角色：</label>
              <span className='text'>{store.role || '--'}</span> */}
          </StyledItem>
        </>
      ) : (
        <Empty description='用户未登录' />
      )}
    </StyledLayout>
  )
})

export default PersonalSetting
