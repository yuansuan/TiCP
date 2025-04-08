/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { account, env } from '@/domain'
import { useStore } from '../store'
import { Dropdown, Input, Menu, message } from 'antd'
import {
  boxFileServer,
  showDirSelector,
  showFailure
} from '@/components/FileMGT'
import { Button } from '@/components'
import { boxServer } from '@/server'

const StyledLayout = styled.div`
  display: flex;
  margin: 4px 0px 20px;

  > .right {
    margin-left: auto;
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()
  const { job } = store
  const state = useLocalStore(() => ({
    get downloadDisabled() {
      if (store.selectedKeys.length === 0) {
        return '无选中文件'
      } else if (job?.useRemote && store.selectedKeys.length > 1) {
        return '作业文件回传中，请不要批量下载文件'
      }

      return false
    }
  }))

  const downloadMultiLocal = async () => {
    let size
    const jobParams = job.useRemote
      ? {
          paths: store.selectedKeys,
          sync_id: job.runtime?.download_task_id
        }
      : {
          base: job.id,
          paths: store.selectedKeys.map(filename => `${job.id}/${filename}`),
          bucket: 'result'
        }

    if (store.selectedKeys.length === 1) {
      const file = store.jobFile.getName(store.selectedKeys[0])
      size = file.size
    }

    const params =
      store.selectedKeys.length === 1
        ? { ...jobParams, sizes: [size], types: [true] }
        : jobParams
    await boxServer.download(params)

    store.setSelectedKeys([])
  }

  const downloadMultiCommon = async () => {
    const targetDir = await showDirSelector()
    const existNodes = []
    const dir = await boxFileServer.fetch(targetDir)
    const allDirPaths = dir.flatten().map(item => item.path)
    const pathsObj = store.selectedKeys.reduce((o, p) => {
      const srcPath = `${job.id}/${p}`
      const dstPath = targetDir ? `${targetDir}/${p}` : p
      if (allDirPaths.includes(dstPath)) {
        existNodes.push({
          path: dstPath,
          name: dstPath.split('/').slice(-1)[0],
          isFile: true
        })
      } else {
        o[srcPath] = dstPath
      }
      return o
    }, {})

    if (existNodes.length > 0) {
      const coverNodes = await showFailure({
        actionName: '下载',
        items: existNodes
      })
      if (coverNodes.length > 0) {
        // del dest first
        await boxFileServer.delete(coverNodes.map(item => item.path))

        // add to objMap for cover to mv
        coverNodes.reduce((o, node) => {
          o[`${job.id}/${node.name}`] = `${targetDir}/${node.name}`
          return o
        }, pathsObj)
      }
    }
    if (Object.keys(pathsObj).length > 0) {
      message.success('下载完成')
      store.selectedKeys = []
    }
  }

  const downloadDropdownMenu = (
    <Menu>
      <Menu.Item onClick={downloadMultiLocal}>下载至本地</Menu.Item>
      {!job?.useRemote && (
        <Menu.Item onClick={downloadMultiCommon}>下载至我的文件</Menu.Item>
      )}
    </Menu>
  )

  return (
    <StyledLayout>
      {!state.downloadDisabled && (
        <Dropdown overlay={downloadDropdownMenu} placement='bottomLeft'>
          <Button type='primary'>下载</Button>
        </Dropdown>
      )}
      {state.downloadDisabled && (
        <Button type='primary' disabled={state.downloadDisabled}>
          下载
        </Button>
      )}
      <div className='right'>
        <Input.Search
          allowClear
          placeholder='输入文件名称搜索'
          value={store.searchKey}
          onChange={({ target: { value } }) => store.setSearchKey(value)}
        />
      </div>
    </StyledLayout>
  )
})
