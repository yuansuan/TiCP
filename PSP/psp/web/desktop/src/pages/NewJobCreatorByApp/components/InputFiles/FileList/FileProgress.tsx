/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Progress } from 'antd'
import { observer } from 'mobx-react-lite'
import { FileProgressStyle } from './style'
import { JobFile } from '@/domain/JobBuilder/JobFile'

import { Status } from '@/components'

interface IProps {
  file: JobFile
}

export const FileProgress = observer((props: IProps) => {
  const { file } = props

  const statusMapping = {
    error: {
      status: 'error',
      text: '失败'
    },
    success: {
      status: 'success',
      text: '上传成功'
    },
    done: {
      status: 'success',
      text: '上传成功'
    },
    paused: {
      status: 'primary',
      text: '已暂停'
    }
  }
  return (
    <FileProgressStyle>
      {file.status === 'uploading' ? (
        <Progress
          status='active'
          percent={parseInt(file.percent.toString())}
          showInfo={true}
        />
      ) : (
        <>
          <div
            className='status-icon'
            style={{
              background: statusMapping[file.status || 'done'].color,
              border: `2px solid ${statusMapping[file.status || 'done'].borderColorBase}`
            }}
          />
          <Status
            type={statusMapping[file.status || 'done'].status}
            text={statusMapping[file.status || 'done'].text}
          />
        </>
      )}
    </FileProgressStyle>
  )
})
