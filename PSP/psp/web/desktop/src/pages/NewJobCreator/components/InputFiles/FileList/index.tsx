/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Checkbox, Tooltip, Divider } from 'antd'
import { Observer, observer } from 'mobx-react'
import { CaretDownOutlined, CaretRightOutlined } from '@ant-design/icons'
import { FileListStyle } from './style'
import { FileName } from './FileName'
import { FileProgress } from './FileProgress'
import { TreeTable } from '@/components'
import { Button } from '@/components'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { useStore } from '../store'
import { useStore as JobStore } from '../../../store'
import { buryPoint } from '@/utils'
import { showTextEditor } from '@/components/TextEditor'
import { newBoxServer } from '@/server'
import { EDITABLE_SIZE } from '@/constant'

interface IProps {
  fileTree: FileTree
}
const isImage = type => /gif|jpe?g|tiff|png|webp|bmp$/i.test(type)

export const FileList = observer(({ fileTree }: IProps) => {
  const store = useStore()
  function bfsSearch(
    list: FileTree[],
    pickFunc: (item: FileTree) => boolean
  ): Array<FileTree> {
    const queue = []
    queue.push(...list)
    const results = []
    while (queue.length > 0) {
      const item = queue.shift()
      if (pickFunc.apply(item, [item])) {
        results.push(item)
      } else {
        if (Array.isArray(item.children)) {
          queue.push(...item.children)
        }
      }
    }

    return results
  }
  const columns = [
    {
      title: '文件名',
      key: 'name',
      render: (text, record, index, isExpand) => (
        <FileName
          record={record}
          isExpand={isExpand}
          createFolder={store.createFolder}
          copyOrUploadFilesFromFileManager={store.copyOrUploadFilesFromFileManager}
          upload={store.upload}
        />
      )
    },
    {
      title: '文件大小',
      key: 'displaySize',
      width: 140,
      render: (text, record) => {
        if (!record.parent) {
          return <span className='disabled'>/</span>
        }
        if (record.isFile) {
          return text
        }
        return '--'
      }
    },
    {
      title: '状态',
      key: 'status',
      width: 120,
      render: (text, record) => {
        if (!record.parent) {
          return <span className='disabled'>/</span>
        }
        if (record.isFile) {
          return <FileProgress file={record} />
        }
        return '--'
      }
    },
    {
      title: '主文件',
      key: 'mainFile',
      width: 160,
      headerRender: text => {
        return (
          <div className='mainFile'>
            <span style={{ color: '#F5222D', marginRight: '2px' }}>*</span>
            {text}
          </div>
        )
      },
      render: (text, record) => {
        if (!record.parent) {
          return <span className='disabled'>/</span>
        }
        return (
          <Observer>
            {() => (
              <div className='mainFileChecker'>
                <Tooltip
                  title={
                    record.name.includes(' ') ? '主文件名不能包含空格' : ''
                  }>
                  <Checkbox
                    disabled={!record.isFile || record.name.includes(' ')}
                    value={record.id}
                    checked={record.isMain}
                    onChange={e => {
                      if (e.target.checked) {
                        buryPoint({
                          category: '作业提交',
                          action: '设置主文件'
                        })
                      }
                      record.isMain = e.target.checked
                    }}>
                    {record.isFile ? '设置为主文件' : '不可设置'}
                  </Checkbox>
                </Tooltip>
              </div>
            )}
          </Observer>
        )
      }
    },
    {
      title: '操作',
      key: 'options',
      width: 180,
      render: (text, record) =>
        record.parent ? (
          <>
            <Button
              type='link'
              onClick={() => {
                buryPoint({
                  category: '作业提交',
                  action: '删除'
                })
                store.deleteNode(record.id)
              }}>
              删除
            </Button>
            <Divider type='vertical' />
            <Button
              type='link'
              disabled={!record.isFile || isImage(record.name)}
              onClick={() => {
                buryPoint({
                  category: '作业提交',
                  action: '编辑/查看'
                })
                showTextEditor({
                  path: record.path,
                  fileInfo: {
                    ...record
                  },
                  readonly: !(record.isFile && record.size <= EDITABLE_SIZE),
                  boxServerUtil: newBoxServer
                })
              }}>
              {record.isFile && record.size <= EDITABLE_SIZE ? '编辑' : '查看'}
            </Button>
          </>
        ) : (
          <span className='disabled'>/</span>
        )
    }
  ]

  return (
    <FileListStyle className='uploadedFileTree'>
      <TreeTable
        dataSource={[fileTree]}
        columns={columns}
        expandedKeys={store.expandedKeys}
        onExpand={store.onExpand}
        rowKey='id'
        expandIcon={isExpand =>
          isExpand ? <CaretDownOutlined /> : <CaretRightOutlined />
        }
        onDragEnd={store.onDrop}
      />
    </FileListStyle>
  )
})
