/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Input } from 'antd'
import { useStore } from '../store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Move } from './Move'
import { Refresh } from './Refresh'
import { Delete } from './Delete'
import { Upload } from './Upload'
import { Download } from './Download'
import { TestFileDownload } from './TestFileDownload'
import { env } from '@/domain'
import { Record } from './Record'

const StyledLayout = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;

  > .left {
    display: flex;
    z-index: 99;

    > * {
      margin: 0 4px;
    }
  }

  > .right {
    z-index: 99;
    display: flex;
    align-items: center;
    /* margin-left: auto; */
    .AreaSelectContainer {
      /* min-width: 330px; */
      margin-right: 10px;
      display: flex;
      align-items: center;
      .areaWrap {
        width: 46px;
      }
      .ant-select {
        flex: 1;
      }
    }
  }
`

type Props = {}

export const Toolbar = observer(function Toolbar(props: any) {
  const store = useStore()
  const { searchKey, refresh } = store
  const state = useLocalStore(() => ({
    get emptyDisabled() {
      const { selectedKeys } = store
      if (selectedKeys.length === 0) {
        return '请选择至少一个文件'
      }

      return false
    }
  }))
  const { emptyDisabled } = state

  function search(e) {
    const { value } = e.target
    store.setSearchKey(value)
  }

  return (
    <StyledLayout>
      <div className='left'>
        {store.getWidget('upload') || <Upload />}
        {store.getWidget('download') || <Download disabled={emptyDisabled} />}
        {store.getWidget('move') || <Move disabled={emptyDisabled} />}
        {store.getWidget('delete') || <Delete disabled={emptyDisabled} />}
        {store.getWidget('refresh') || <Refresh />}
        {env.isPersonal &&
          (store.getWidget('testFile') || <TestFileDownload />)}
      </div>
      <div className='right'>
        <Input.Search
          placeholder='请输入文件名'
          maxLength={64}
          allowClear
          value={searchKey}
          onChange={search}
        />
      </div>
    </StyledLayout>
  )
})
