/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect, useCallback } from 'react'
import styled from 'styled-components'
import ReactDOM from 'react-dom'
import { SearchSelect, Modal } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from '../store'
import { appList, uploader } from '@/domain'
import { Http, getComputeType } from '@/utils'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

const Wrapper = styled.div`
  display: flex;
  margin-right: 10px;
  > .ant-select {
    width: 300px;
    margin-right: 10px;
  }
`

interface ISoftwareProps {
  name: string
  icon: string
  type: string
  versions: [string, string, string][]
  action: any
}
let newTaskKey = ''

EE.on(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, ({ taskKey }) => {
  newTaskKey = taskKey
})

export const Software = observer(function Software(props: ISoftwareProps) {
  const store = useStore()
  const ref = useRef(null)
  const state = useLocalStore(() => ({
    visible: false,
    get selected() {
      return props.versions.map(([key]) => key).includes(store.currentAppId)
    },
    get currentAppQueue() {
      return store.data.currentApp?.queues
    },
    // 检查是否有服务器文件上传到超算
    get hasUploadingList() {
      return (
        document.querySelector('.server-file-uploading-to-supercomputing') !==
        null
      )
    }
  }))

  useEffect(() => {
    if (props.action === props.type) {
      if (props.versions.length === 1) {
        selectApp(props.versions[0][0], false)
      } else if (ref.current) {
        // expand select options
        const clickEvent = document.createEvent('MouseEvents')
        clickEvent.initEvent('mousedown', true, true)
        const select: any = ReactDOM.findDOMNode(ref.current)
        select.querySelector('.ant-select-selector').dispatchEvent(clickEvent)
        selectApp(props.versions[0][0], false)
      }
    }

    return () => {
      newTaskKey = ''
    }
  }, [props.action])

  const abortAllTask = () => {
    return Http.delete('/storage/hpcUpload/abortAllTask', {
      params: {
        taskKey: newTaskKey
      }
    })
  }
  const selectApp = useCallback(
    (id, needCheckAppVersion = true) => {
      const currentApp = appList.list.find(item => item.id === id)

      if (
        props.versions?.length >= 2 &&
        store.data.currentApp?.compute_type !== currentApp?.compute_type &&
        needCheckAppVersion &&
        (store.currentJobHasParams || state.hasUploadingList)
      ) {
        Modal.confirm({
          title: '混合云切换',
          content: `当前数据会全部丢失，确认是否切换？`,
          okText: '确认',
          cancelText: '取消',
          onOk: () => {
            store.updateData({
              currentApp: currentApp
            })
            // 取消本地上传到服务器
            if (uploader.fileList.length > 0) {
              uploader.fileList.forEach(file => uploader.remove(file?.uid))
            }
            // 取消服务器到超算的文件
            newTaskKey && abortAllTask()
          }
        })
      } else {
        store.updateData({
          currentApp: currentApp
        })
      }
    },
    [store.data.currentApp?.compute_type]
  )
  return (
    <Wrapper>
      <SearchSelect
        ref={ref}
        onClick={e => e.stopPropagation()}
        getPopupContainer={triggerNode => triggerNode.parentElement}
        size='middle'
        options={props.versions.map(([key, name, compute_type]) => ({
          key,
          name: `${getComputeType(compute_type)}-${name}`
        }))}
        {...(state.selected && {
          value: `${getComputeType(store.data.currentApp?.compute_type)}-${
            store.data.currentApp?.version
          }`
        })}
        placeholder='请选择版本'
        onChange={selectApp}
        onDeselect={selectApp}
      />
    </Wrapper>
  )
})
