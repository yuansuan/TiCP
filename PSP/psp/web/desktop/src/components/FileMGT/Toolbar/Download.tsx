/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { useStore } from '../store'
import styled from 'styled-components'
import { useStore as useJobStore } from '@/pages/NewJobDetail/store'
import {
  boxFileServer,
  showDirSelector,
  showFailure
} from '@/components/NewFileMGT'
import { newBoxServer } from '@/server'
import { Dropdown, Menu, message, Tooltip } from 'antd'
import { QuestionCircleOutlined } from '@ant-design/icons'

type Props = {
  disabled: boolean | string
  userName: string
}

const StyledLayout = styled.div`
  display: flex;
  margin: 4px 0px 20px;
`

export const Download = observer(function Download({
  disabled,
  userName
}: Props) {
  const store = useStore()
  const jobStore = useJobStore()
  const { job } = jobStore
  const { selectedKeys } = store

  const downloadMultiLocal = async () => {
    const nodes = store.dir.filterNodes(item => selectedKeys.includes(item.id))
    const types = nodes.map(item => item.isFile)
    const sizes = nodes.map(item => item.size)
    const jobParams = {
      base: job?.work_dir,
      paths: nodes.map(
        file => `${job.work_dir?.replace(/\/$/, '')}/${file.name}`
      ),
      is_cloud: false // 作业回传到本地才可以下载，所以这里就是false
    }

    const params =
      selectedKeys.length === 1 ? { ...jobParams, sizes, types } : jobParams
    await newBoxServer.download({
      ...params,
      user_name: userName
    })

    store.setSelectedKeys([])
  }

  const downloadMultiCommon = async () => {
    const targetDir = await showDirSelector()
    const existNodes = []
    const serverDir = await boxFileServer.fetch(targetDir)
    const allDirPaths = serverDir.flatten().map(item => item.path)
    let dstPath = ''

    const nodes = store.dir.filterNodes(item => selectedKeys.includes(item.id))

    const pathsObj = nodes.reduce((o, p) => {
      const srcPath = `${job?.work_dir?.replace(/\/$/, '')}/${p.name}`
      dstPath = targetDir ? `${targetDir}/${p.name}` : p.name
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
    const srcDirPaths = nodes
      .filter(file => !file?.isFile)
      .map(item => item?.path)
    const srcFilePaths = nodes
      .filter(file => file?.isFile)
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
      <Menu.Item onClick={downloadMultiLocal} key={'1'}>
        下载至本地
      </Menu.Item>
      {/* {job.isCloud && (
        <Menu.Item onClick={downloadMultiCommon} key={'2'}>
          下载至我的文件
        </Menu.Item>
      )} */}
    </Menu>
  )
  return (
    <StyledLayout>
      {!disabled && (
        <Dropdown overlay={downloadDropdownMenu} placement='bottomLeft'>
          <Tooltip
            title={'为提高下载速度，建议数据回传完成之后再下载文件。'}
            visible={job?.data_state !== 'Downloaded' && job?.isCloud}>
            <Button type='primary'>
              {job?.data_state !== 'Downloaded' && job?.isCloud && (
                <QuestionCircleOutlined />
              )}{' '}
              下载
            </Button>
          </Tooltip>
        </Dropdown>
      )}
      {disabled && (
        <Button type='primary' disabled={disabled}>
          下载
        </Button>
      )}
    </StyledLayout>
  )
})
