/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { Checkbox } from 'antd'
import { Modal, Icon, Button } from '@/components'
import styled from 'styled-components'

const StyledLayout = styled.div`
  > .list {
    margin: 8px 0;
    padding: 24px 16px;
    max-height: 324px;
    overflow: auto;
    background-color: ${({ theme }) => theme.backgroundColorBase};

    > .item {
      font-size: 14px;
      margin: 5px 0;
      display: flex;
      align-items: center;

      > .name {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        max-width: calc(100% - 30px);
      }

      .anticon {
        color: ${({ theme }) => theme.primaryColor};
        margin: 0 8px;
      }
    }
  }

  > .footer {
    padding-top: 12px;
  }
`

type Node = {
  uid?: string
  name: string
  isFile: boolean
}

type Props = {
  coverable?: boolean
  type?: 'conflict'
  actionName: string
  items: Array<Node>
}

type ModalActions = {
  onCancel: () => void
  onOk: (props: Node[]) => void
}

const FailureFileContent = function FailureFileContent({
  coverable = true,
  type = 'conflict',
  actionName,
  items,
  onCancel,
  onOk
}: Props & ModalActions) {
  const [checkedList, setCheckedList] = useState(
    items.map(item => ({
      checked: true,
      item
    }))
  )

  function onChange(index, checked, item) {
    const list = [...checkedList]
    list.splice(index, 1, { checked, item })
    setCheckedList(list)
  }

  function onSelectAll(checked) {
    const list = [...checkedList]
    list.forEach(item => (item.checked = checked))
    setCheckedList(list)
  }

  function onConfirm() {
    return checkedList.filter(item => item.checked).map(item => item.item)
  }

  return (
    <StyledLayout>
      <div className='tip'>
        {actionName}操作{!coverable && '失败'}：以下{items.length}
        个文件/文件夹已存在
      </div>
      <div className='list'>
        {coverable && (
          <Checkbox
            indeterminate={
              !checkedList.every(item => item.checked) &&
              checkedList.some(item => item.checked)
            }
            checked={checkedList.every(item => item.checked)}
            onChange={({ target: { checked } }) => onSelectAll(checked)}>
            {' '}
            全选
          </Checkbox>
        )}
        {items.map((item, index) => {
          return (
            <div className='item' key={index}>
              {coverable && (
                <Checkbox
                  checked={checkedList[index]?.checked}
                  onChange={({ target: { checked } }) =>
                    onChange(index, checked, item)
                  }
                />
              )}
              <Icon type={item.isFile ? 'file_table' : 'folder_close'} />
              <span className='name'>{item.name}</span>
            </div>
          )
        })}
      </div>
      <Modal.Footer
        className='footer'
        CancelButton={null}
        OkButton={
          (coverable && (
            <>
              <Button type='default' onClick={() => onOk([])}>
                取消
              </Button>
              <Button
                disabled={
                  !checkedList.some(item => item.checked) &&
                  '未选中需要覆盖的文件'
                }
                type='primary'
                onClick={() => onOk(onConfirm())}>
                覆盖
              </Button>
            </>
          )) || (
            <Button
              type='primary'
              onClick={() => {
                onOk(null)
              }}>
              确定
            </Button>
          )
        }
      />
    </StyledLayout>
  )
}

export async function showFailure(props: Props) {
  return await Modal.show({
    title: `文件${props.actionName}`,
    content: ({ onCancel, onOk: onContentOk }) => (
      <FailureFileContent {...props} onCancel={onCancel} onOk={onContentOk} />
    ),
    footer: null
  })
}
