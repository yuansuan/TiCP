/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Select } from 'antd'
import { observer } from 'mobx-react-lite'
import { useStore } from '../store'

const Wrapper = styled.div`
  width: 300px;
  > .ant-select {
    width: 100%;
  }
`

export const ProjectSelector = observer(function ProjectSelector() {
  const store = useStore()

  useEffect(() => {
    store.fetchProjectList(true)
  }, [])

  return (
    <Wrapper>
      <Select
        placeholder='请选择提交作业所属项目'
        value={store.projectId}
        onDropdownVisibleChange={open => {
          if (open) store.fetchProjectList(false)
        }}
        onChange={store.setProjectId}>
        {store.projectList.map(item => {
          return (
            <Select.Option key={item.id} value={item.id}>
              项目: {item.name}
            </Select.Option>
          )
        })}
      </Select>
    </Wrapper>
  )
})
