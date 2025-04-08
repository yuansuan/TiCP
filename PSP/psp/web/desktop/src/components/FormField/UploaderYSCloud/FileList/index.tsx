/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Checkbox, Button } from 'antd'
import { Observer, observer } from 'mobx-react'
import {
  CaretDownOutlined,
  CaretRightOutlined,
  FolderOpenOutlined,
  FolderOutlined
} from '@ant-design/icons'
import { FileListStyle } from './style'
import { TreeTable } from '@/components'
import { useStore, useModel, Context } from './store'
import { SpinMask } from '@ys/components/dist/Mask/SpinMask'
import { Props } from './store'

const canExpandKey = 'isDir'

const List = observer(() => {
  const isDisableDelBtn = status =>
    status !== 'done' &&
    status !== 'aborted' &&
    status !== 'paused' &&
    status !== 'error'

  const store = useStore()

  const columns = [
    {
      title: '文件名',
      key: 'name',
      width: 340,
      headerRender: text => {
        return <div style={{ marginLeft: 22 }}>{text}</div>
      },
      render: (text, record, index, isExpand) => {
        return (
          <>
            {isExpand ? (
              <FolderOpenOutlined />
            ) : record[canExpandKey] ? (
              <FolderOutlined />
            ) : null}
            <span
              style={{
                marginLeft: 4
              }}>
              {record.name}
            </span>
          </>
        )
      }
    },
    {
      title: '来源',
      key: 'displayFrom',
      width: 90
    },
    {
      title: '文件大小',
      key: 'displaySize',
      width: 90
    },
    {
      title: '状态',
      key: 'displayStatus',
      width: 85
    },
    {
      title: '主文件',
      key: 'mainFile',
      width: 135,
      headerRender: text => {
        return (
          <div className='mainFile'>
            <span style={{ color: '#F5222D', marginRight: '2px' }}>*</span>
            {text}
          </div>
        )
      },
      render: (text, record) => {
        return (
          <Observer>
            {() => (
              <div className='mainFileChecker'>
                <Checkbox
                  disabled={record.isDir}
                  value={record.path}
                  checked={record.isMain}
                  onChange={e => {
                    record.update({
                      isMain: e.target.checked
                    })
                    store.setMain(record.path, e.target.checked)
                  }}>
                  {!record.isDir ? '设置主文件' : '不可设置'}
                </Checkbox>
              </div>
            )}
          </Observer>
        )
      }
    },
    {
      title: '操作',
      key: 'options',
      width: 120,
      render: (text, record) =>
        record.from === 'local' ? (
          <Button
            style={{ padding: 0 }}
            type='link'
            disabled={isDisableDelBtn(record.status)}
            onClick={() => {
              store.deleteNode(record.path, record.status === 'done')
            }}>
            删除
          </Button>
        ) : (
          <span className='disabled'>/</span>
        )
    }
  ]

  return (
    <FileListStyle className='uploadedFileTree'>
      <TreeTable
        dataSource={store.fileList}
        columns={columns}
        expandedKeys={store.expandedKeys}
        onExpand={store.onExpand}
        rowKey='path'
        canExpandKey={canExpandKey}
        expandIcon={isExpand =>
          isExpand ? <CaretDownOutlined /> : <CaretRightOutlined />
        }
        draggable={false}
      />
      {store.loading && (
        <div className='mask-wrapper'>
          <SpinMask />
        </div>
      )}
    </FileListStyle>
  )
})

export const FileList = observer(function FileList(props: Props) {
  const model = useModel(props)
  return (
    <Context.Provider value={model}>
      <List />
    </Context.Provider>
  )
})
