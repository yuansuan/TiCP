/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { env } from '@/domain'
import { useStore } from '../store'
import { Dropdown, Input, Menu, message } from 'antd'
import {
  boxFileServer,
  showDirSelector,
  showFailure
} from '@/components/NewFileMGT'
import { Button } from '@/components'
import { newBoxServer } from '@/server'

const StyledLayout = styled.div`
  display: flex;
  margin: 4px 0px 20px;

  > .right {
    margin-left: auto;
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()
  const { job, jobFile } = store
  const state = useLocalStore(() => ({
    get downloadDisabled() {
      if (store.selectedKeys.length === 0) {
        return '无选中文件'
      }

      return false
    }
  }))

  const downloadMultiLocal = async () => {
    let size
    const jobParams = {
      base: job?.work_dir,
      paths: store.selectedKeys.map(
        filename => `${job.work_dir?.replace(/\/$/, '')}/${filename}`
      ),
      is_cloud: job.isCloud
    }

    if (store.selectedKeys.length === 1) {
      const file = store.jobFile.getName(store.selectedKeys[0])
      size = file.size
    }

    const params =
      store.selectedKeys.length === 1
        ? { ...jobParams, sizes: [size], types: [true] }
        : jobParams
    await newBoxServer.download({
      ...params
    })

    store.setSelectedKeys([])
  }

  const downloadMultiCommon = async () => {
    const targetDir = await showDirSelector()
    const existNodes = []
    const dir = await boxFileServer.fetch(targetDir)
    const allDirPaths = dir.flatten().map(item => item.path)
    let dstPath = ''
    let selectedFiles = []
    const pathsObj = store.selectedKeys.reduce((o, p) => {
      const file = store.jobFile.getName(p)
      selectedFiles.push(file)
      const srcPath = `${job?.work_dir?.replace(/\/$/, '')}/${p}`
      dstPath = targetDir ? `${targetDir}/${p}` : p
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
        await boxFileServer.delete(
          coverNodes.map(item => item.path),
          false,
          false
        )

        // add to objMap for cover to mv
        coverNodes.reduce((o, node) => {
          o[
            `${job?.work_dir?.replace(/\/$/, '')}/${node.name}`
          ] = `${targetDir}/${node.name}`
          return o
        }, pathsObj)
      }
    }
    const srcDirPaths = selectedFiles
      .filter(file => file?.is_dir)
      .map(item => item?.path)
    const srcFilePaths = selectedFiles
      .filter(file => !file?.is_dir)
      .map(item => item?.path)
    if (Object.keys(pathsObj).length > 0) {
      await newBoxServer.linkToCommon({
        current_path: job?.work_dir,
        dest_dir_path: targetDir,
        src_dir_paths: srcDirPaths,
        src_file_paths: srcFilePaths,
        user_name: job.user_name
      })

      message.success('开始下载，请稍候刷新文件管理')
      store.selectedKeys = []
    }
  }

  const downloadDropdownMenu = (
    <Menu>
      <Menu.Item onClick={downloadMultiLocal}>下载至本地</Menu.Item>
      {/* {job.isCloud && (
        <Menu.Item onClick={downloadMultiCommon}>下载至我的文件</Menu.Item>
      )} */}
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
