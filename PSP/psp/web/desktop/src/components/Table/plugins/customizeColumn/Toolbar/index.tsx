/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Popover } from 'antd'
import { ColumnManager } from './ColumnManager'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'

export const StyledToolbar = styled.div`
  position: absolute;
  top: 0;
  right: 40px;
  cursor: pointer;

  .anticon {
    &:hover {
      color: #63a9ff;
    }

    &.active {
      color: #63a9ff;
      transform: rotate(-90deg);
    }
  }
`

interface IProps {
  Icon: any
  columns: any[]
  config: Array<{ key: string; active: boolean }>
  setConfig: (config: any) => void
}

export const Toolbar = observer(function Toolbar(props: IProps) {
  const { Icon } = props
  const state = useLocalStore(() => ({
    visible: false,
    setVisible(visible) {
      this.visible = visible
    },
  }))

  return (
    <StyledToolbar>
      <Popover
        visible={state.visible}
        onVisibleChange={visible => state.setVisible(visible)}
        placement='bottom'
        content={
          <ColumnManager
            Icon={Icon}
            config={props.config}
            setConfig={props.setConfig}
            columns={props.columns}
          />
        }
        trigger='click'>
        <Icon
          className={state.visible ? 'active' : ''}
          style={{ verticalAlign: 'middle' }}
          type='table_column'
        />
      </Popover>
    </StyledToolbar>
  )
})
