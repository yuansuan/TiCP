/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect, useCallback, useMemo } from 'react'
import { Observer } from 'mobx-react-lite'
import { IProps, IData } from './types'
import { Wrapper, StyledTable } from './style'
import { getAllParentKey } from './utils'
import { TableHeader, TableRow } from './components'

export const TreeTable = (props: IProps) => {
  const {
    columns,
    dataSource,
    rowKey = 'key',
    indentSize = 24,
    expandedKeys,
    onExpand = (keys: string[], data: IData) => {},
    defaultExpandAll = false,
    expandIcon = (isExpand: boolean) =>
      isExpand ? <div>-</div> : <div>+</div>,
    childrenField = 'children',
    onDragEnd = (key1: string, key2: string) => {},
    draggable = true
  } = props

  const [exdKeys, setExpandKeys] = useState<string[]>([])
  const finalExdKeys = useMemo(
    () => expandedKeys || exdKeys,
    [expandedKeys, exdKeys]
  )

  // 默认展开全部
  useEffect(() => {
    const pkeys = getAllParentKey(dataSource, rowKey)
    if (!expandedKeys && defaultExpandAll) {
      setExpandKeys(pkeys)
    }
  }, [dataSource])

  const toggleParentExpand = useCallback(
    (data: IData) => {
      const key = data[rowKey]
      let targetKeys = []
      if (finalExdKeys.includes(key)) {
        targetKeys = finalExdKeys.filter(i => i !== key)
      } else {
        targetKeys = [...finalExdKeys, key]
      }
      if (!expandedKeys) {
        // 非受控
        setExpandKeys(targetKeys)
      }
      onExpand(targetKeys, data)
    },
    [finalExdKeys, onExpand, expandedKeys, rowKey]
  )

  // TODO 支持draggable
  const renderRow = useCallback(
    (rowData: IData, rowIndex: number, level = 0) => {
      const rows = [
        <TableRow
          key={rowData[rowKey]}
          rowData={rowData}
          rowIndex={rowIndex}
          level={level}
          toggleParentExpand={toggleParentExpand}
          expandedKeys={finalExdKeys}
          indentSize={indentSize}
          expandIcon={expandIcon}
          rowKey={rowKey}
          childrenField={childrenField}
          columns={columns}
          // onDragEnd={onDragEnd}
          // draggable={draggable}
          // draggable={false}
        />
      ]
      if (rowData[childrenField] && finalExdKeys.includes(rowData[rowKey])) {
        rowData[childrenField].map((cld: IData, index: number) =>
          rows.push(...renderRow(cld, index, level + 1))
        )
      }
      return rows
    },
    [
      columns,
      rowKey,
      finalExdKeys,
      childrenField,
      toggleParentExpand,
      expandIcon,
      indentSize
      // onDragEnd,
      // draggable
    ]
  )

  return (
    <Wrapper>
      <Observer>
        {() => (
          // <DraggableWarpper>
          <StyledTable>
            <TableHeader columns={columns} />
            <tbody>
              {dataSource.map((rowData, rowIndex) =>
                renderRow(rowData, rowIndex)
              )}
            </tbody>
          </StyledTable>
          // </DraggableWarpper>
        )}
      </Observer>
    </Wrapper>
  )
}
