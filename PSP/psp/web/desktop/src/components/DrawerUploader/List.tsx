/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer, Observer } from 'mobx-react-lite'
import styled from 'styled-components'
import { uploader } from '@/domain'
import { UploadProgress } from '@/components/UploadProgress'
import { Button, Image } from '@/components'
import { UploaderFile } from '@/components/Uploader'
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

export const List = observer(function List({ list }: Props) {
  const getAction = file => {
    switch (file.status) {
      case 'error':
        return [
          <Button
            size='small'
            key={0}
            type='primary'
            onClick={() => uploader.retry(file.uid)}>
            重试
          </Button>
        ]
      case 'paused':
        return [
          <Button
            size='small'
            key={0}
            type='primary'
            onClick={() => uploader.resume(file.uid)}>
            开始
          </Button>
        ]
      case 'uploading':
        return [
          <Button
            size='small'
            key={0}
            type='primary'
            onClick={() => uploader.pause(file.uid)}>
            暂停
          </Button>
        ]
      default:
        return []
    }
  }

  const renderListRow = ({ index, style }) => {
    const file = list[index]
    return (<Observer key={file.uid}>
      {() => (
        <div onClick={e => e.stopPropagation()} className='item' style={style}>
          <UploadProgress
            context={file}
            direction={['local', 'box']}
            operators={[
              ...getAction(file),
              <Button
                size='small'
                key={1}
                onClick={() => uploader.remove(file.uid)}>
                取消
              </Button>
            ]}
          />
        </div>
      )}
    </Observer>)
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
                className="List"
                height={height}
                itemCount={list.length}
                itemSize={84}
                width={width}
              >
              {renderListRow}
            </FixedSizeList>
          )}
        </AutoSizer>
        )}
    </Wrapper>
  )
})
