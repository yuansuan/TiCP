/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Checkbox, Button, Tooltip, message } from 'antd'
import { Observer, observer } from 'mobx-react'
import {
  CaretDownOutlined,
  CaretRightOutlined,
  FolderOpenOutlined,
  FolderOutlined,
  FileOutlined,
} from '@ant-design/icons'
import { FileListStyle } from './style'
import { TreeTable } from '@/components'
import { useStore, useModel, Context } from './store'
import { SpinMask } from '@ys/components/dist/Mask/SpinMask'
import { Props } from './store'
import { currentFileList } from './Files'
import LocalUploader from '../LocalUploader'
import ServerUploader from '../ServerUploader'
import { EditableText, Icon } from '@/components'
import { Validator } from '@/utils'
import EditAction from './EditAction'
import { Modal } from '@/components'

const canExpandKey = 'isDir'

const EDITABLE_SIZE = 3 * 1024 * 1024

const isEditable = file => {
  if (file.size > EDITABLE_SIZE) {
    return {
      editable: false,
      message: '文件大小超过 3M',
    }
  }

  if (!file?.is_text) {
    return {
      editable: false,
      message: '非文本文件',
    }
  }

  return {
    editable: true,
    message: '编辑',
  }
}

const List = observer(() => {

  const isDisableDelBtn = (status) => (status !== 'done' && status !== 'aborted' && status !== 'paused' && status !== 'error')

  const store = useStore()

  const renderSubFileUploaderBtn = (record) => (
    record.isDir && 
    <>
      <LocalUploader
        upload={(params, isDir) => store.upload(params, isDir, true, record.path)}
        beforeUpload={params => store.beforeUpload(params, true)}
        data={{
          master: record.path,
        }}>
        <Tooltip placement='top' title='本地文件'>
          <Button icon='upload' type="link" style={{ marginRight: '5px' }}>
          </Button>
        </Tooltip>
      </LocalUploader>
      <ServerUploader
        onUpload={files => {
          store.upload(files, false, false, record.path)
        }}>
        {upload => (
          <Tooltip placement='top' title='远程文件'>
            <Button icon='upload' type="link" onClick={upload}>
            </Button>
          </Tooltip>
        )}
      </ServerUploader>
    </>
  ) 

  const beforeRename = value => {

    if (!value) {
      message.error(`重命名失败：文件名不能为空`)
      return false
    }

    const { error } = Validator.filename(value)
    if (error) {
      message.error(`重命名失败：${error.message}`)
      return false
    }

    return true
  }

  const renderViewFileAction = (record) => {
    return record.isDir ? 
    (<span style={{paddingLeft: 6}}>{record.name}</span>)
     : 
    (<EditAction
      title={'查看'}
      node={record}
      readOnly={true}>
      <span style={{paddingLeft: 6}}>{record.name}</span>
    </EditAction>)
  }

  const createEditableText = (record, onClick?: any) => (
    <EditableText
      Text={() => renderViewFileAction(record)}
      EditIcon={
        <Tooltip title='重命名'>
          <Icon type='edit-filled' />
        </Tooltip>
      }
      onClick={onClick}
      defaultValue={record.name}
      defaultShowEdit={record.isRenaming}
      beforeConfirm={beforeRename}
      onConfirm={name => store.rename(name, record)}
    />
  )

  const columns = [
    {
      title: '文件名',
      key: 'name',
      width: 320,
      headerRender: text => {
        return <div style={{ marginLeft: 22 }}>{text}</div>
      },
      render: (text, record, index, isExpand) => {
        return (
          <>
            {record.isDir ? (isExpand ? (
              <FolderOpenOutlined />
            ) : record[canExpandKey] ? (
              <FolderOutlined />
            ) : null) : <FileOutlined style={{paddingLeft: 6}}/>}
            <Observer>
              {
                () => createEditableText(record)
              }
            </Observer>
          </>
        )
      },
    },
    // {
    //   title: '来源',
    //   key: 'displayFrom',
    //   width: 90,
    // },
    {
      title: '文件大小',
      key: 'displaySize',
      width: 90,
    },
    {
      title: '状态',
      key: 'displayStatus',
      width: 85,
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
                      isMain: e.target.checked,
                    })
                    store.setMain(record.path, e.target.checked)
                  }}>
                  {!record.isDir ? '设置主文件' : '不可设置'}
                </Checkbox>
              </div>
            )}
          </Observer>
        )
      },
    },
    {
      title: '操作',
      key: 'options',
      width: 160,
      render: (text, record) =>
        (
          <>
           {renderSubFileUploaderBtn(record)}
           {record.isDir &&
              <Tooltip placement='top' title='创建子文件夹'>
                <Button
                  style={{padding: 0}}
                  type="link"
                  icon={'folder-add'}
                  disabled={isDisableDelBtn(record.status)}
                  onClick={() => {
                    store.createFolder(record)
                  }}>
                </Button>
              </Tooltip>
            }
            {!record.isDir && 
              <EditAction
                title={'编辑'}
                node={record}
                readOnly={false}>
                  <Tooltip 
                    placement='top' 
                    title={isEditable(record).message}>
                    <Button 
                      style={{padding: 0}}
                      type="link"
                      icon={'edit'}
                      disabled={!isEditable(record).editable}>
                    </Button>
                  </Tooltip>
              </EditAction>
            }
           <Tooltip placement='top' title='删除'>
              <Button
                style={{padding: 0}}
                type="link"
                icon={'delete'}
                disabled={isDisableDelBtn(record.status)}
                onClick={() => {
                  Modal.showConfirm({
                    title: '删除操作确认框',
                    content: `确认删除${record.isDir ? '文件夹' : '文件'}${record.name}? (请谨慎操作，当前操作无法恢复)`,
                  }).then(() => {
                    store.deleteNode(record.path, record.status === 'done')
                  })
                }}>
              </Button>
            </Tooltip>
          </>
        )
    },
  ]

  return (
    <FileListStyle className='uploadedFileTree'>
      <TreeTable
        dataSource={currentFileList.files}
        columns={columns}
        expandedKeys={store.expandedKeys}
        onExpand={store.onExpand}
        rowKey='path'
        canExpandKey={canExpandKey}
        expandIcon={isExpand =>
          isExpand ? <CaretDownOutlined /> : <CaretRightOutlined />
        }
        draggable={false}
        onScroll={store.onScroll}
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
