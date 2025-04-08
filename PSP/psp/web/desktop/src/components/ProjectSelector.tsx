/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Dropdown, message } from 'antd'
import { CaretDownFilled } from '@ant-design/icons'
import { projectList, env } from '@/domain'
import { history } from '@/utils'

const StyledLayout = styled.div`
  margin: 10px;
  padding: 8px 10px;
  background: rgba(0, 52, 180, 0.2);
  color: white;
  display: flex;

  > .text {
    display: inline-block;
    max-width: 170px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  > .icon {
    margin-left: auto;
  }
`

const StyledItem = styled.div`
  width: 170px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 8px 16px;
  cursor: pointer;

  &:hover {
    background: #f0f2f5;
    color: ${({ theme }) => theme.linkColor};
  }
`

type OverlayProps = {
  onClick: (id: string) => void
}

const Overlay = observer(function Overlay({ onClick }: OverlayProps) {
  return (
    <ul>
      {projectList.list.map(item => (
        <li
          key={item.id}
          style={{ padding: 0 }}
          onClick={() => onClick(item.id)}>
          <StyledItem title={item.name}>{item.name}</StyledItem>
        </li>
      ))}
    </ul>
  )
})

export const ProjectSelector = observer(function ProjectSelector() {
  const state = useLocalStore(() => ({
    visible: false,
    setVisible(visible) {
      this.visible = visible
    }
  }))

  return (
    <Dropdown
      visible={state.visible}
      onVisibleChange={visible => state.setVisible(visible)}
      placement='bottomCenter'
      overlayStyle={{
        maxHeight: 300,
        overflowY: 'auto',
        backgroundColor: 'white'
      }}
      overlay={
        <Overlay
          onClick={async id => {
            state.setVisible(false)
            await env.changeProject(id)
            history.push('/dashboard')
            message.success('工作空间切换成功')
          }}
        />
      }>
      <StyledLayout>
        <span className='text'>{env.project?.name}</span>
        <span className='icon'>
          <CaretDownFilled />
        </span>
      </StyledLayout>
    </Dropdown>
  )
})
