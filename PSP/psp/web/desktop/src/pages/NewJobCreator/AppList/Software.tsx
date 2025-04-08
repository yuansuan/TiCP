/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect } from 'react'
import styled from 'styled-components'
import { SearchSelect } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { buryPoint } from '@/utils'
import { useStore } from '../store'
import { appList } from '@/domain'

const Wrapper = styled.div`
  width: 180px;
  > .ant-select {
    width: 100%;
  }
`

interface ISoftwareProps {
  name: string
  icon: string
  type: string
  versions: [string, string,string][]
}

export const Software = observer(function Software(props: ISoftwareProps) {
  const store = useStore()
  const ref = useRef(null)
  const state = useLocalStore(() => ({
    get selected() {
      return props.versions.map(([key]) => key).includes(store.currentAppId)
    }
  }))

  useEffect(() => {
    if (props.versions.length === 1) {
      selectApp(props.versions[0][0])
    }
  }, [])
  function selectApp(id) {
    buryPoint({
      category: '作业提交',
      action: '软件选择'
    })
    store.updateData({
      currentApp: appList.list.find(item => item.id === id)
    })
  }

  return (
    <Wrapper>
      <SearchSelect
        ref={ref}
        onClick={e => e.stopPropagation()}
        getPopupContainer={triggerNode => triggerNode.parentElement}
        size='middle'
        options={props.versions.map(([key, name]) => ({
          key,
          name
        }))}
        {...(state.selected && {
          value: store.data.currentApp.version
        })}
        placeholder='请选择版本'
        onChange={selectApp}
      />
    </Wrapper>
  )
})
