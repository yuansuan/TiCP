/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { useLocalStore, useObserver } from 'mobx-react-lite'
import { Dropdown, Input, Checkbox, Button, Tooltip, Switch, Modal } from 'antd'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { Icon } from '@/components'
import { FileActions } from './FileActions'
import { ToolbarStyle, FileSearchResultListStyle } from './style'
import { useStore as useJobStore } from '../../store'
import { buryPoint } from '@/utils'
import { WarningOutlined } from '@ant-design/icons'
import { Label } from '@/components'

const minimatch = require('minimatch')

interface IProps {
  fileTree: FileTree
}

export const Toolbar = ({ fileTree }: IProps) => {
  const jobStore = useJobStore()
  const store = useLocalStore(() => ({
    fileTree,
    val: undefined,
    visible: false,
    valMode: '关键字', // filename keyword or filepath glob
    selection: [],
    get selectedMainFile() {
      return store.fileTree.filterNodes(node => {
        return node.isFile && node.isMain === true
      })
    },
    get filterFiles() {
      if (!this.val) return []
      return store.fileTree.filterNodes(
        // 支持 glob way filter
        // TODO 过滤掉文件名中的空格
        node => {
          if (this.valMode === 'glob') {
            return node.isFile && minimatch(node.path, this.val)
          } else {
            return node.isFile && node.name.indexOf(this.val) > -1
          }
        }
      )
    },
    get indeterminate() {
      return (
        this.selection.length > 0 &&
        this.selection.length < this.filterFiles.length
      )
    },
    get checkAll() {
      return this.selection.length === this.filterFiles.length
    }
  }))

  React.useEffect(() => {
    store.fileTree = fileTree
  }, [fileTree])

  const onChange = selection => {
    store.selection = selection
  }

  const onCheckboxChange = e => {
    const { checked, value } = e.target
    buryPoint({
      category: '作业提交',
      action: '设置主文件'
    })
    fileTree.tapNodes(
      node => node.id === value,
      node => {
        node.isMain = checked
      }
    )
  }

  const checkedMainFiles = (checked, files) => {
    let fileIds = files.map(f => f.id)
    fileTree.tapNodes(
      node => fileIds.includes(node.id),
      node => {
        node.isMain = checked
      }
    )
  }

  const onCheckAll = e => {
    // clear main files
    checkedMainFiles(false, store.filterFiles)

    if (
      e.target.checked &&
      store.selection.length ===
        store.filterFiles.filter(file => !file.name.includes(' ')).length
    ) {
      store.selection = []
      return
    }

    if (e.target.checked) {
      store.selection = store.filterFiles
        .filter(file => !file.name.includes(' '))
        .map(file => file.id)
      buryPoint({
        category: '作业提交',
        action: '批量设置主文件'
      })
      checkedMainFiles(
        true,
        store.filterFiles.filter(file => !file.name.includes(' '))
      )
    } else {
      store.selection = []
    }
  }

  const changeSearchMode = () => {
    store.valMode = store.valMode === 'glob' ? '关键字' : 'glob'
    store.val = undefined
    store.visible = false
  }

  const syncSelectedMainFiles = () => {
    // 求交集 O(m+n)
    let mainFileIds = store.selectedMainFile.map(mf => mf.id)
    let fileIds = new Set(store.filterFiles.map(f => f.id) || [])
    let selectedIds = mainFileIds.filter(id => fileIds.has(id))
    store.selection = selectedIds || []
  }

  const setSelectToMainFile = () => {
    buryPoint({
      category: '作业提交',
      action: '批量设置主文件'
    })
    fileTree.tapNodes(
      node => store.selection.includes(node.id),
      node => {
        node.isMain = true
      }
    )
    store.visible = false
    store.selection = []
    store.val = undefined
  }

  const switchTempDirPathMode = value => {
    if (jobStore.fileTree.children.length === 0) {
      jobStore.setIsTempDirPath(!value)
      jobStore.setTempDirPath('')
      jobStore.resetFileTree()
      return
    }

    Modal.confirm({
      title: '作业计算模式切换',
      content: '当前上传的作业模型数据会全部丢失，确认是否切换？',
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        jobStore.setIsTempDirPath(!value)
        jobStore.setTempDirPath('')
        jobStore.resetFileTree()
      }
    })
  }

  const resetAllMainFile = () => {
    buryPoint({
      category: '作业提交',
      action: '重置主文件'
    })
    fileTree.tapNodes(
      () => true,
      node => {
        node.isMain = false
      }
    )
    store.selection = []
  }

  const FileSearchResultList = useObserver(() => (
    <FileSearchResultListStyle>
      {store.filterFiles.length ? (
        <>
          <li>
            <Checkbox
              onChange={onCheckAll}
              indeterminate={store.indeterminate}
              checked={store.checkAll}>
              全部
            </Checkbox>
          </li>
          <Checkbox.Group
            value={store.selection}
            onChange={onChange}
            className='dropdown-content'>
            {store.filterFiles.map(file => (
              <li key={file.id} title={file.name}>
                <Tooltip
                  title={file.name.includes(' ') ? '主文件名不能包含空格' : ''}>
                  <Checkbox
                    key={file.id}
                    onChange={onCheckboxChange}
                    disabled={file.name.includes(' ')}
                    value={file.id}>
                    {file.name}
                  </Checkbox>
                </Tooltip>
              </li>
            ))}
          </Checkbox.Group>
        </>
      ) : (
        <li>无匹配文件</li>
      )}
    </FileSearchResultListStyle>
  ))

  return useObserver(() => (
    <ToolbarStyle>
      <div
        style={{
          display: 'flex',
          alignItems: 'center'
        }}>
        <FileActions />
        <span>
          {!jobStore.isTempDirPath && jobStore.tempDirPath !== '' ? (
            <>
              <Tooltip
                title={
                  '指定工作目录后，对目录中文件的操作是不可逆的，请谨慎操作!'
                }>
                <WarningOutlined style={{ color: '#ec942c' }} rev={''} />
              </Tooltip>
              {
                <span style={{ marginLeft: 2 }}>
                  工作目录：{jobStore?.tempDirPath?.replace(/^\.\//, '')}
                </span>
              }
            </>
          ) : (
            ''
          )}
        </span>
      </div>
      <div className='right-search'>
        {/* <span className='label'>选择主文件：</span>
        <Dropdown
          visible={store.visible}
          onVisibleChange={e => {
            store.visible = e
            syncSelectedMainFiles()
          }}
          trigger={['click']}
          overlay={FileSearchResultList}
          placement='bottomLeft'>
          <Input
            addonBefore={<a onClick={changeSearchMode}>{store.valMode === 'glob' ? 'glob' : '关键字'}</a>}
            value={store.val}
            onChange={e => {
              let _tmp = e.target.value.trim()
              store.val = _tmp
              syncSelectedMainFiles()
              if (_tmp) {
                store.visible = true
              }
            }}
            placeholder={store.valMode === 'glob' ? 'glob模式' : '文件名中关键字'}
          />
        </Dropdown> */}

        <div>
          <span style={{ marginRight: 2 }}>原地计算</span>
          <Tooltip title='指定文件目录后, 计算结果将存在这个文件目录中, 不会额外创建文件目录!'>
            <Icon style={{ width: 15, height: 15 }} type={'help-circle'} />
          </Tooltip>
          <Switch
            checkedChildren='开启'
            unCheckedChildren='关闭'
            style={{ marginLeft: 5, marginRight: 20 }}
            checked={!jobStore.isTempDirPath}
            onChange={value => switchTempDirPathMode(value)}
          />
        </div>
        <Button onClick={resetAllMainFile}>重置主文件</Button>
      </div>
    </ToolbarStyle>
  ))
}
