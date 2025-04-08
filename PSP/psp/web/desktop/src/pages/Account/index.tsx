/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import AccountDetail from './AccountDetail/index'
import AccountList from './AccountList/index'
import { StyledLayout } from './style'

const Account = () => {
  return (
    <StyledLayout>
      <AccountDetail />
      <AccountList />
    </StyledLayout>
  )
}

export default Account
