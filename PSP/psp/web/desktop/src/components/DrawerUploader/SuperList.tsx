/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer, Observer } from 'mobx-react-lite'
import styled from 'styled-components'
import { Http } from '@/utils'
import { UploadProgress } from '@/components/UploadProgress'
import { Button, Image } from '@/components'
import { UploaderFile } from '@ys/components/dist/Uploader'
import { FixedSizeList } from 'react-window'
import AutoSizer from 'react-virtualized-auto-sizer'

const Wrapper = styled.div`
  height: 350px;

  .empty-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 80%;
    margin: auto;
    padding: 20px 0;

    & > img {
      width: 150px;
    }
  }

  .ant-list-bordered .ant-list-item {
    padding: 0;
  }

  .operators .ant-btn {
    width: 64px;
  }
`

type Props = { list: UploaderFile[] }

export const SuperList = observer(function List({ list = [] }: Props) {
  const resumeTask = async rowData => {
    const { task_key, upload_id } = rowData
    await Http.put('/storage/hpcUpload/resumeTask', {
      task_key,
      upload_id
    })
  }
  const cancelTask = async rowData => {
    const { task_key, upload_id } = rowData
    await Http.put('/storage/hpcUpload/cancelTask', {
      task_key,
      upload_id
    })
  }
  const getAction = file => {
    const actions = {
      1: {
        label: '重试',
        onClick: () => resumeTask(file)
      },
      2: {
        label: '取消',
        onClick: () => cancelTask(file)
      },
      3: {
        label: '取消',
        onClick: () => cancelTask(file)
      }
    }

    const actionData = actions[file.state]

    if (actionData) {
      const { label, onClick } = actionData
      return [
        <Button size='small' key={0} type='primary' onClick={onClick}>
          {label}
        </Button>
      ]
    }

    return []
  }

  const renderListRow = ({ index, style }) => {
    const file = list[index]
    return (
      <Observer key={file.uid}>
        {() => (
          <div
            onClick={e => e.stopPropagation()}
            className='item server-file-uploading-to-supercomputing'
            style={style}>
            <UploadProgress
              context={file}
              direction={['local', 'box']}
              operators={[
                ...getAction(file)
              ]}
            />
          </div>
        )}
      </Observer>
    )
  }

  return (
    <Wrapper>
      {list.length === 0 && (
        <div className='empty-wrapper'>
          <Image.Empty />
        </div>
      )}
      {list.length !== 0 && (
        <AutoSizer>
          {({ height, width }) => (
            <FixedSizeList
              className='List'
              height={height}
              itemCount={list.length}
              itemSize={84}
              width={width}>
              {renderListRow}
            </FixedSizeList>
          )}
        </AutoSizer>
      )}
    </Wrapper>
  )
})
