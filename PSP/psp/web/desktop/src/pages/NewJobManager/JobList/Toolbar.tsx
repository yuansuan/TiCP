/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Dropdown, Menu, message } from 'antd'
import { Button, Modal } from '@/components'
import { useStore } from './store'
import { env } from '@/domain'
import { showDownloader } from './showDownloader'
import { newBoxServer, jobServer } from '@/server'
import { useStore as fatherStore } from '../store'
import { getUrlParams } from '@/utils/Validator'
import { SearchFiles } from './SearchFiles'
import { Http } from '@/utils'

const StyledLayout = styled.div`
  padding: 20px 20px 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  > .left {
    > * {
      margin-right: 10px;
    }
  }
`

export const Toolbar = observer(function Toolbar() {
  const params = getUrlParams()
  const store = useStore()
  const managerStore = fatherStore()

  const { model } = store

  const state = useLocalStore(() => ({
    downloadDisabled() {
      // if (!store.selectedKeys.length) return '请选择要下载的作业'
      // console.log('store.selectedKeys: ', store.selectedKeys)
      // // const invalidJobNames = store.selectedKeys
      // //   .map(key => model.list.find(item => item.id === key))
      // //   .filter(item => !!item)
      // //   .map(item => item.name)
      // // if (invalidJobNames.length > 0) {
      // //   return `作业 ${invalidJobNames.join(', ')} 未回传完成，不能下载`
      // // }
      // return false
      if (store.selectedKeys.length === 1) {
        return (
          store.selectedKeys.findIndex(
            state =>
              state === 'Running' ||
              state === 'Suspended' ||
              state === 'Pending'
          ) !== -1
        )
      } else {
        return true
      }
    },
    get cancelDisabled() {
      if (!store.selectedKeys.length) return '请选择要取消的作业'
      const invalidJobNames = store.selectedKeys
        .map(key => model.list.find(item => item.id === key))
        .filter(item => !!item)
        .map(item => item.name)
      if (invalidJobNames.length > 0) {
        return `作业 ${invalidJobNames.join(', ')} 不能取消`
      }
      return false
    },
    get deleteDisabled() {
      if (!store.selectedKeys.length) return '请选择要删除的作业'
      const invalidJobNames = store.selectedKeys
        .map(key => model.list.find(item => item.id === key))
        .filter(item => !!item)
        .map(item => item.name)
      if (invalidJobNames.length > 0) {
        return `作业 ${invalidJobNames.join(', ')} 未结束，不能删除`
      }
      return false
    },

    get canStop() {
      if (!store.selectedKeys.length) return '请选择要暂停的作业'
      const invalidJobNames = store.selectedKeys
        .map(key => model.list.find(item => item.id === key))
        .filter(item => !!item)
        .map(item => item.name)
      if (invalidJobNames.length > 0) {
        return `作业 ${invalidJobNames.join(', ')} 暂停，不能暂停`
      }
      return false
    },
    get canTerminate() {
      const filteredArray = store.model.list.filter(item =>
        store.selectedKeys.some(key => key === item.out_job_id)
      )

      if (store.selectedKeys.length === 1) {
        return (
          store.selectedKeys.findIndex(
            state =>
              state === 'Running' || state === 'Suspend' || state === 'Pending'
          ) !== -1
        )
      } else {
        return true
      }
    }
  }))

  async function downloadJobs() {
    const jobs = store.selectedKeys
      .map(key => model.list.find(item => item.id === key))
      .map(item => ({
        jobId: item.id,
        jobName: item.name
      }))
    store.setSelectedKeys([])
  }

  async function downloadJobsToSearch() {
    await Modal.show({
      title: '搜索下载至本地',
      footer: null,
      width: 800,
      bodyStyle: {
        height: 600
      },
      okText: '下载',
      content: ({ onCancel, onOk }) => {
        const OK = async selectedDownloadFile => {
          await newBoxServer.download({
            paths: selectedDownloadFile,
            base: '.'
          })
          onOk()
        }
        return <SearchFiles onOk={OK} selectedJob={store.selectedKeys} />
      }
    })
  }

  async function downloadJobsToCommon() {
    const selectedJobs = model.list
      .filter(item => store.selectedKeys.includes(item.id))
      .map(item => ({
        id: item.id,
        name: item.name
      }))
    const resolvedJobs = await showDownloader(selectedJobs)
    if (Object.keys(resolvedJobs).length > 0) {
      message.success('下载完成')
      store.setSelectedKeys([])
    }
  }
  async function terminateJob(job) {
    await Modal.showConfirm({
      title: '确认终止',
      content: '是否确认终止选中作业？'
    })
    // store.selectedKeys
    await Http.post('job/terminate', {
      out_job_id: job
    })
  }
  async function cancelJobs() {
    await Modal.showConfirm({
      title: '确认取消',
      content: '是否确认取消选中作业？'
    })
    await jobServer.cancel(store.selectedKeys)
    await store.refresh()
    store.setSelectedKeys([])
    message.success('取消成功')
  }

  async function deleteJobs() {
    await Modal.showConfirm({
      title: '确认删除',
      content: '删除作业同时会删除作业产生的文件，是否确认？'
    })
    await jobServer.delete(store.selectedKeys)
    const paths = store.selectedKeys
      .map(key => model.list.find(item => item.id === key))
      .map(item => item.id)
    await newBoxServer.delete({
      paths
    })

    await store.refresh()
    //id && (await jobSetList.fetch())
    store.setSelectedKeys([])
    message.success('删除成功')
  }

  return (
    <StyledLayout>
      <div className='left'>
        {/* {!!state.downloadDisabled ? (
          <Dropdown
            overlay={
              <Menu>
                <Menu.Item
                  // disabled={managerStore?.downloadFlag}
                  onClick={downloadJobs}>
                  下载至本地
                </Menu.Item>
                <Menu.Item
                  // disabled={managerStore?.downloadFlag}
                  onClick={downloadJobsToSearch}>
                  搜索下载至本地
                </Menu.Item>
                <Menu.Item
                  // disabled={managerStore?.downloadFlag}
                  onClick={downloadJobsToCommon}>
                  下载至我的文件
                </Menu.Item>
              </Menu>
            }
            placement='bottomLeft'>
            <Button type='primary'>下载</Button>
          </Dropdown>
        ) : (
          <Button disabled={state.downloadDisabled}>下载</Button>
        )} */}
        {/* <Button disabled={state.canTerminate} onClick={terminateJob}>
          终止
        </Button> */}
        {/* <Button onClick={cancelJobs} disabled={state.cancelDisabled}>
          取消
        </Button>
        <Button onClick={deleteJobs} disabled={state.canTerminate}>
          删除
        </Button> */}
      </div>
    </StyledLayout>
  )
})
