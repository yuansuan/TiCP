/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Badge, Dropdown, Tabs } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { uploader } from '@/domain'
import { List } from './List'
import { Box } from './Box'
import styled from 'styled-components'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

const StyledLayout = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
`

const StyledPanel = styled.div`
  padding: 20px;
  width: 558px;
  max-height: 475px;
  overflow: auto;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

  .ant-tabs-ink-bar {
    width: 80px !important;
  }

  > .body {
    .tabName {
      padding: 0 12px;
    }

    .ant-tabs-bar {
      margin-bottom: 0;
    }

    .ant-list-bordered {
      border: none;
      border-radius: 0;
    }
  }

  .item {
    margin: 0;
    height: 84px;
    display: flex;
    align-items: center;
    width: 100%;
    border-bottom: 1px solid #e8e8e8;
  }

  .ant-list-item-meta-description {
    word-break: break-all;
  }

  .ant-list-bordered {
    > .ant-list-footer {
      padding: 0;
    }
  }

  .ant-list-footer {
    padding: 0;

    .footer {
      display: flex;

      > div {
        flex: 1;
        text-align: center;
        padding: 12px 0;
        cursor: pointer;

        &:hover {
          color: black;
        }

        &:not(:last-child) {
          border-right: 1px solid #e8e8e8;
        }
      }
    }
  }
`

const { TabPane } = Tabs

export const Uploader = observer(function Uploader() {
  const state = useLocalStore(() => ({
    visible: false,
    setVisible(visible) {
      this.visible = visible
    },
    get displayList() {
      return uploader.fileList.filter(file =>
        ['uploading', 'paused', 'error'].includes(file.status)
      )
    },
  }))

  useEffect(() => {
    const handler = ({ visible }) => {
      if (visible) {
        state.setVisible(visible)
      } else {
        if (state.displayList.length === 0) state.setVisible(visible)
      }
    }

    EE.on(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, handler)

    return () => {
      EE.off(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, handler)
    }
  }, [])

  return (
    <Dropdown
      visible={state.visible}
      onVisibleChange={visible => state.setVisible(visible)}
      placement='bottomRight'
      overlay={
        <StyledPanel>
          <div className='body' onClick={e => e.stopPropagation()}>
            <Tabs defaultActiveKey='upload'>
              <TabPane
                key='upload'
                tab={<span className='tabName'>正在上传</span>}>
                <List list={[...state.displayList]} />
              </TabPane>
            </Tabs>
          </div>
        </StyledPanel>
      }
      trigger={['click']}>
      <StyledLayout>
        <Box dropdownVisible={state.visible} />
        <Badge offset={[4, 0]} count={state.displayList.length} />
      </StyledLayout>
    </Dropdown>
  )
})
