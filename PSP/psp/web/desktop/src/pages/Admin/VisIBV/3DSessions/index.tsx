/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'

import { observer } from 'mobx-react-lite'
import { StyledLayout } from './style'
import { useStore } from '../store'
import { SessionList } from './SessionList'
import { Export } from './Export'
import { Toolbar } from './Toolbar'

interface IProps {
  height?: number
}

const SessionListPage = observer(function ListPage(props: IProps) {
  const store = useStore()

  return (
    <StyledLayout>
      <Toolbar />
      <div className='info'>
        {/* <Export disabled={!['admin', 'finance', 'it', 'cs', 'wx_sc'].includes(currentUser.role)} /> */}
        {/* <div>共计{store.model.total}个会话</div> */}
      </div>
      <SessionList height={props.height || 400}/>
    </StyledLayout>
  )
})

export default SessionListPage
