/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Tabs, Drawer } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { uploader } from '@/domain'
import { List } from './List'
import { Http } from '@/utils'
import { SuperList } from './SuperList'
import styled from 'styled-components'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

const StyledPanel = styled.div`
  padding: 20px;
  width: 558px;
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
let timer = null
let newTaskKey = ''

EE.on(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, ({ taskKey }) => {
  newTaskKey = taskKey
})
export const Uploader = observer(function Uploader() {
  const state = useLocalStore(() => ({
    visible: false,
    setVisible(visible) {
      this.visible = visible
    },
    previousNewDataLength: 0,
    setPreviousNewDataLength(current) {
      this.previousNewDataLength = current
    },

    showSuperList: false,
    setShowSuperList(bool) {
      this.showSuperList = bool
    },
    get displayList() {
      return uploader.fileList.filter(file =>
        ['uploading', 'paused', 'error'].includes(file.status)
      )
    },
    superFileLists: [],
    setSuperFiles(fileList) {
      const statusMap = {
        1: 'error',
        2: 'uploading',
        3: 'paused',
        4: 'done',
        5: 'removed'
      }
      this.superFileLists = fileList.map(item => {
        return {
          ...item,
          name: item.file_name,
          loaded: item.current_size,
          percent: (item.current_size / item.total_size) * 100,
          status: statusMap[item.state],
          state: item.state,
          size: item.total_size
        }
      })
    }
  }))

  function getUploaderFileList(taskKey) {
    return Http.get(`/storage/hpcUpload/fileTaskList`, {
      params: {
        taskKey: taskKey
      }
    })
  }

  useEffect(() => {
    const fetchDataAndProcess = async taskKey => {
      if (!taskKey) return
      const result = await getUploaderFileList(taskKey)

      if (result) {
        const newData = result.data.map(item => ({
          ...item,
          task_key: taskKey
        }))
        state.setSuperFiles(newData)

        if (newData.length < state.previousNewDataLength) {
          EE.emit(EE_CUSTOM_EVENT.SERVER_FILE_TO_SUPERCOMPUTING, {
            file_status: 'success'
          })
        }

        taskKey && state.setShowSuperList(true)
        state.setPreviousNewDataLength(newData.length)

        if (newData.length === 0) {
          clearInterval(timer)
          EE.emit(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, { taskKey: '' })
          state.setVisible(false)
          state.setShowSuperList(false)
        }
      }
    }

    const handleUploadDropdownToggle = ({ visible }) => {
      state.setVisible(visible)

      if (visible && newTaskKey) {
        timer = setInterval(() => fetchDataAndProcess(newTaskKey), 1000)
      } else {
        state.setShowSuperList(false)
      }
    }
    EE.on(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, handleUploadDropdownToggle)

    return () => {
      clearInterval(timer)
      EE.off(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, handleUploadDropdownToggle)
    }
  }, [])

  const onClose = () => {
    state.setVisible(false)
  }
  return (
    <Drawer
      title='上传窗口'
      placement='right'
      onClose={onClose}
      width='600'
      maskClosable
      visible={state.visible as boolean}>
      <StyledPanel>
        <div className='body' onClick={e => e.stopPropagation()}>
          <Tabs defaultActiveKey='upload'>
            <TabPane
              key='upload'
              tab={<span className='tabName'>正在上传</span>}>
              {state.showSuperList ? (
                <SuperList list={[...state.superFileLists]} />
              ) : (
                <List list={[...state.displayList]} />
              )}
            </TabPane>
          </Tabs>
        </div>
      </StyledPanel>
    </Drawer>
  )
})
