/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { useLocalStore, useObserver } from 'mobx-react-lite'
import { Dropdown, Input, Checkbox, Button, Tooltip } from 'antd'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { FileActions } from './FileActions'
import { ToolbarStyle, FileSearchResultListStyle } from './style'
import { buryPoint } from '@/utils'

const minimatch = require('minimatch')

interface IProps {
  fileTree: FileTree
}

export const Toolbar = ({ fileTree }: IProps) => {
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
      <FileActions />
      <div className='right-search'>
        <span className='label'>选择主文件aa：</span>
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
            addonBefore={
              <a onClick={changeSearchMode}>
                {store.valMode === 'glob' ? 'glob' : '关键字'}
              </a>
            }
            value={store.val}
            onChange={e => {
              let _tmp = e.target.value.trim()
              store.val = _tmp
              syncSelectedMainFiles()
              if (_tmp) {
                store.visible = true
              }
            }}
            placeholder={
              store.valMode === 'glob' ? 'glob模式' : '文件名中关键字'
            }
          />
        </Dropdown>
        <Button onClick={resetAllMainFile}>重置主文件</Button>
      </div>
    </ToolbarStyle>
  ))
}
