/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Switch, Tooltip } from 'antd'
import { observer, useLocalStore } from 'mobx-react'
import React, { useEffect } from 'react'

import { Icon, Modal, Table } from '@/components'
import { formatFileSize } from '@/utils/formatter'
import { statusMap } from '@/domain/Uploader/Task'
import { IHarmonyFile } from '../index'
import LocalUploader from '../LocalUploader'
import ServerUploader from '../ServerUploader'
import { Wrapper } from './style'

interface IProps {
  files: IHarmonyFile[]
  deleteAction: (path: string, master?: string) => void
  setMain: (path: string, checked: boolean) => void
  model: any
  beforeUploadLocalFile: (params: any) => Promise<any>
  uploadLocalFile: (params: any, isDir: boolean) => void
  uploadServerFile: (files: any[], master?: string) => void
}

const FileTree = observer(function FileTree({
  files,
  deleteAction,
  setMain,
  model,
  uploadLocalFile,
  beforeUploadLocalFile,
  uploadServerFile
}: IProps) {
  const state = useLocalStore(() => ({
    expandedKeys: [],
    updateExpandedKeys(keys) {
      state.expandedKeys = keys
    },
    files: files,
    setFiles(arr) {
      this.files = arr
    },
    get finalFiles() {
      let res = []

      // generator master/slave fileTree
      this.files.forEach(item => {
        res.push(harmony(item))

        // flatten slave files
        const expanded = this.expandedKeys.includes(item.path)
        if (expanded) {
          if (item.slaveFiles) {
            res = [...res, ...item.slaveFiles]
          }
        }
      })

      return res
    },
    get mainFiles() {
      // just filter main files in the first layer
      return this.files.filter(item => item.isMain).map(item => item.path)
    }
  }))

  const harmony = (file: IHarmonyFile) => {
    return {
      ...file,
      canSetMain: file.status === 'done' && !file.isDir && !file.master
    }
  }

  useEffect(() => {
    state.setFiles(files)
  }, [files])

  useEffect(() => {
    // expand all mainFiles by default
    state.updateExpandedKeys(state.mainFiles)
  }, [])

  useEffect(() => {
    // harmony expandedKeys when mainFiles change
    const expandedKeys = state.expandedKeys.filter(item =>
      state.mainFiles.includes(item)
    )
    state.updateExpandedKeys(expandedKeys)
  }, [state.mainFiles])

  const deleteFile = (path, master) => {
    Modal.showConfirm({
      content: '确认删除该文件吗？'
    }).then(() => {
      deleteAction(path, master)
    })
  }
  // toggle main file
  const toggleMain = (path, checked) => {
    // set main file
    if (checked) {
      // default expand
      onExpand(true, path)
    }

    setMain(path, checked)
  }

  const onExpand = (expanded, path) => {
    let keys = state.expandedKeys
    if (expanded) {
      keys = [...new Set([...state.expandedKeys, path])]
    } else {
      keys = state.expandedKeys.filter(item => item !== path)
    }

    state.updateExpandedKeys(keys)
  }
  const { expandedKeys } = state
  const { isMasterSlave } = model
  const width = 750

  return (
    <Wrapper>
      <Table
        props={{
          data: state.finalFiles,
          rowKey: 'path',
          ...(state.finalFiles.length <= 10
            ? { autoHeight: true }
            : { height: 330 }),
          width,
          expandedRowKeys: expandedKeys
        }}
        columns={
          [
            {
              header: '',
              props: {
                width: width * 0.05
              },
              cell: {
                props: { dataKey: 'path' },
                render: ({ rowData }) => {
                  const expanded = state.expandedKeys.includes(rowData.path)
                  return (
                    <>
                      {rowData.isMain ? (
                        <Icon
                          type='more'
                          onClick={() => onExpand(!expanded, rowData.path)}
                        />
                      ) : null}
                    </>
                  )
                }
              }
            },
            {
              header: '名称',
              props: {
                width: width * 0.35,
                resizable: true
              },
              cell: {
                props: {
                  dataKey: 'name'
                },
                render: ({ rowData, dataKey }) => {
                  const isSlave = rowData.master
                  return (
                    <span
                      className={`name ${isSlave ? 'slave' : ''}`}
                      title={rowData.path}>
                      {rowData.isDir ? (
                        <Icon type='folder' />
                      ) : (
                        <Icon type='file' />
                      )}
                      {rowData[dataKey]}
                    </span>
                  )
                }
              }
            },
            {
              header: '来源',
              props: {
                width: width * 0.1,
                resizable: true
              },
              cell: {
                props: {
                  dataKey: 'name'
                },
                render: ({ rowData }) => {
                  const from = rowData.from
                  return <span>{from === 'local' ? '本地' : '服务器'}</span>
                }
              }
            },
            {
              header: '大小',
              props: {
                width: width * 0.15,
                resizable: true
              },
              cell: {
                props: { dataKey: 'size' },
                render: ({ rowData, dataKey }) => {
                  if (rowData.isDir) {
                    return '--'
                  } else {
                    return formatFileSize(rowData[dataKey])
                  }
                }
              }
            },
            {
              header: '状态',
              props: {
                width: width * 0.15,
                resizable: true
              },
              cell: {
                props: {
                  dataKey: 'status'
                },
                render: ({ rowData, dataKey }) => (
                  <div>{statusMap[rowData[dataKey]]}</div>
                )
              }
            },
            ...(isMasterSlave
              ? [
                  {
                    header: '主文件',
                    props: {
                      width: width * 0.1
                    },
                    cell: {
                      props: {
                        dataKey: 'path'
                      },
                      render: ({ rowData }) => {
                        return (
                          <>
                            {rowData.canSetMain && (
                              <Switch
                                size='small'
                                onClick={(checked, e) => e.stopPropagation()}
                                checked={rowData.isMain}
                                onChange={checked =>
                                  toggleMain(rowData.path, checked)
                                }
                              />
                            )}
                          </>
                        )
                      }
                    }
                  }
                ]
              : []),
            {
              header: '操作',
              props: {
                width: width * 0.1
              },
              cell: {
                props: {
                  dataKey: 'path'
                },
                render: ({ rowData, dataKey }) => {
                  const uploadDisabled = !rowData.isMain
                  const from = rowData.from

                  return (
                    <div className='operate' onClick={e => e.stopPropagation()}>
                      {uploadDisabled ? (
                        <Icon className='disabled' type='upload' />
                      ) : from === 'local' ? (
                        <LocalUploader
                          upload={uploadLocalFile}
                          beforeUpload={beforeUploadLocalFile}
                          data={{
                            master: rowData.path
                          }}>
                          <Tooltip placement='top' title='上传从文件'>
                            <Icon type='upload' />
                          </Tooltip>
                        </LocalUploader>
                      ) : (
                        <ServerUploader
                          onUpload={files => {
                            uploadServerFile(files, rowData.path)
                          }}>
                          {upload => (
                            <Tooltip placement='top' title='上传从文件'>
                              <Icon type='upload' onClick={upload} />
                            </Tooltip>
                          )}
                        </ServerUploader>
                      )}
                      <Icon
                        type='delete'
                        onClick={() =>
                          deleteFile(rowData[dataKey], rowData.master)
                        }
                      />
                    </div>
                  )
                }
              }
            }
          ] as any
        }
      />
    </Wrapper>
  )
})

export default FileTree
