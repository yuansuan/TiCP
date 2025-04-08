/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import { Checkbox, Tooltip, Divider, Spin } from 'antd'
import { Observer, observer } from 'mobx-react'
import { CaretDownOutlined, CaretRightOutlined } from '@ant-design/icons'
import { FileListStyle } from './style'
import { FileName } from './FileName'
import { FileProgress } from './FileProgress'
import { TreeTable } from '@/components'
import { Button } from '@/components'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { useStore } from '../store'
import { useStore as useJobStore } from '../../../store'
import { buryPoint } from '@/utils'
import { showTextEditor } from '@/components/TextEditor'
import { newBoxServer } from '@/server'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { EDITABLE_SIZE } from '@/constant'

interface IProps {
  fileTree: FileTree
}
const isImage = type => /gif|jpe?g|tiff|png|webp|bmp$/i.test(type)
const canReadFile = status => status === 'error' || status === 'paused'
export const FileList = observer(({ fileTree }: IProps) => {
  const store = useStore()
  const jobStore = useJobStore()
  const [myFileLoading, setMyFileLoading] = useState(false)

  useEffect(() => {
    EE.on(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, ({ taskKey }) => {
      setMyFileLoading(!!taskKey)
    })
  }, [])

  useEffect(() => {
    setTimeout(() => {
      if (jobStore.expandFlag) {
        store.mainFilesExpand()
      }
      jobStore.expandFlag = false
    }, 100)
  }, [jobStore.expandFlag])

  const columns = [
    {
      title: '文件名',
      key: 'name',
      render: (text, record, index, isExpand) => (
        <FileName
          record={record}
          isExpand={isExpand}
          createFolder={store.createFolder}
          copyOrUploadFilesFromFileManager={
            store.copyOrUploadFilesFromFileManager
          }
          myFileLoading={myFileLoading}
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
                      if (!jobStore.isTempDirPath) {
                        fileTree.tapNodes(
                          () => true,
                          node => {
                            node.isMain = false
                          }
                        )
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
              // disabled={record.path === jobStore.tempDirPath}
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
              disabled={
                !record.isFile ||
                isImage(record.name) ||
                canReadFile(record.status) ||
                record.size > EDITABLE_SIZE
              }
              onClick={() => {
                buryPoint({
                  category: '作业提交',
                  action: '编辑/查看'
                })
                const node = jobStore.fileTree.filterFirstNode(
                  item => item.id === record.id
                )
                if (node && node.name) {
                  const { path } = node
                  const tmpPath = jobStore.isTempDirPath
                    ? path
                    : path.replace(new RegExp(jobStore.tempDirPath), '')
                  const newPath = (
                    jobStore.tempDirPath +
                    '/' +
                    tmpPath
                  ).replace(/\/\//g, '/')
                  if (node.status === 'done' || node.status === 'success') {
                    showTextEditor({
                      path: newPath,
                      fileInfo: {
                        ...record,
                        cross: jobStore.isTempDirPath,
                        is_cloud: jobStore.isCloud,
                        path: newPath
                      },
                      readonly: true,
                      // readonly: !(record.isFile && record.size <= EDITABLE_SIZE),
                      boxServerUtil: newBoxServer
                    })
                  }
                }
              }}>
              {/* {record.isFile && record.size <= EDITABLE_SIZE ? '编辑' : '查看'} */}
              {'查看'}
            </Button>
          </>
        ) : (
          <span className='disabled'>/</span>
        )
    }
  ]

  return (
    <FileListStyle className='uploadedFileTree'>
      <Spin spinning={store.loading} tip='文件加载中...'>
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
      </Spin>
    </FileListStyle>
  )
})
