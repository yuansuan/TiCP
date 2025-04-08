/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Dropdown, Menu, message } from 'antd'
import { Button } from '@/components'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { showFailure } from '../Failure'
import { v4 as uuid } from 'uuid'
import { UploadProps } from '@/components/Uploader'
import { uploader, NewBoxHttp, env } from '@/domain'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { globalSizes } from '@/domain/Box/states'

const StyledLayout = styled.div``

const UPLOAD_ID = uuid()

type Props = {
  upload?: (props?: Partial<UploadProps>) => void
}

export const Upload = observer(function Upload(props?: Props) {
  const store = useStore()
  const [refresh] = store.useRefresh()
  const { dir, server } = store
  const upload =
    props.upload ||
    // TODO: 需要修改为新版上传
    (props => {
      uploader.upload({
        action: '/filemanager/upload',
        httpAdapter: NewBoxHttp(),
        ...props,
        origin: props.origin,
        data: {
          ...props.data,
          bucket: 'common',
        }
      })
    })

  function _upload(directory = false) {
    const uniqueID = uuid()
    const dirFinal = `/${dir.path.replace(/^\//, '')}`
    upload({
      origin: UPLOAD_ID,
      by: 'chunk',
      multiple: true,
      data: {
        directory,
        _uid: uniqueID,
        dir: dirFinal
      },
      directory,
      beforeUpload: async files => {
        const resolvedFiles = []
        const rejectedFiles = []
        files.forEach(file => {
          const filePath = file.webkitRelativePath || file.name
          const fileName = filePath.split('/')[0]

          if (dir.getDuplicate({ id: undefined, name: fileName })) {
            rejectedFiles.push(file)
          } else {
            resolvedFiles.push(file)
          }
        })

        globalSizes[uniqueID] = 0
        files.forEach(file => {
          globalSizes[uniqueID] += file.size
        })

        if (rejectedFiles.length > 0) {
          if (directory) {
            const filePath =
              rejectedFiles[0].webkitRelativePath || rejectedFiles[0].name
            const topDirName = filePath.split('/')[0]

            const coverNodes = await showFailure({
              actionName: '上传',
              items: [
                {
                  isFile: false,
                  name: topDirName
                }
              ]
            })
            if (coverNodes.length > 0) {
              // remove dir
              await server.delete([`${dirFinal}/${topDirName}`])
              // should upload newFiles
              resolvedFiles.push(...rejectedFiles)
            }
          } else {
            const coverNodes = await showFailure({
              actionName: '上传',
              items: rejectedFiles.map(item => ({
                name: item.name,
                uid: item.uid,
                isFile: true
              }))
            })
            if (coverNodes.length > 0) {
              await server.delete(
                coverNodes.map(item => `${dirFinal}/${item.name}`),
                
              )
              resolvedFiles.push(...coverNodes)
            }
          }
        }

        if (resolvedFiles.length > 0) {
          message.success('文件开始上传')
          // 触发显示 dropdown
          EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: true })
        }

        return resolvedFiles.map(item => item.uid)
      },
      onChange: ({ file, origin }) => {
        if (origin !== UPLOAD_ID) {
          return
        }
        if (file.status === 'done') {
          // 有文件上传完成，check 是否要关闭 dropdown
          EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: false })
          refresh()
        }
      }
    })
  }

  return (
    <StyledLayout>
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item onClick={() => _upload()}>上传文件</Menu.Item>
            <Menu.Item onClick={() => _upload(true)}>上传目录</Menu.Item>
          </Menu>
        }>
        <Button>上传</Button>
      </Dropdown>
    </StyledLayout>
  )
})
