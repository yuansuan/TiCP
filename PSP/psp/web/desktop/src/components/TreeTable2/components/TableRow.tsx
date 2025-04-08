import React, { useCallback, useMemo, useRef, useEffect } from 'react'
import { useDrag, useDrop } from 'react-dnd'
import styled from 'styled-components'
import { expandIcon, IData, IColumn, onDragEndFn } from '../types'
import { ItemTypes } from '../utils'

const Tr = styled.tr<{ indentSize: number }>`
  .icon-expand {
    text-align: center;
    width: ${props => props.indentSize + 'px'};
    height: ${props => props.indentSize + 'px'};
    line-height: ${props => props.indentSize - 2 + 'px'};
    cursor: pointer;
    font-size: 16px;
    font-style: normal;
  }
`

interface IProps {
  indentSize: number
  rowKey: string
  expandIcon: expandIcon
  toggleParentExpand: (data: IData) => void
  childrenField: string
  expandedKeys: string[]
  rowData: IData
  columns: IColumn[]
  rowIndex: number
  level: number
  onDragEnd: onDragEndFn
  draggable: boolean
}

export function TableRow(props: IProps) {
  const {
    indentSize,
    rowKey,
    expandIcon,
    toggleParentExpand,
    childrenField,
    expandedKeys,
    rowData,
    columns,
    rowIndex,
    level,
    onDragEnd,
    draggable,
  } = props

  const rowKeyData = useMemo(() => rowData[rowKey], [rowData, rowKey])
  const trRef = useRef<HTMLTableRowElement>(null)
  const [{ isDragging }, drag] = useDrag({
    item: { [rowKey]: rowKeyData, type: ItemTypes.TableRow, data: rowData },
    collect: monitor => ({
      isDragging: !!monitor.isDragging(),
    }),
    end(item, monitor) {
      if (monitor.didDrop()) {
        onDragEnd((item as any)[rowKey], monitor.getDropResult()[rowKey])
      }
    },
  })
  const [{ active }, drop] = useDrop({
    accept: ItemTypes.TableRow,
    canDrop: (item, monitor) => {
      if (!rowData[childrenField]) {
        return false
      }
      // 判断是否是移到了子集
      const isParent = (source: IData, target: IData): boolean => {
        if (!source[childrenField]) return false
        const index = source[childrenField].findIndex(
          (c: any) => c[rowKey] === target[rowKey]
        )
        if (index > -1) return true
        return isParent(source[childrenField], target)
      }

      return !isParent((item as any).data, rowData)
    },
    collect: monitor => ({
      active: monitor.isOver({ shallow: true }) && monitor.canDrop(),
    }),
    drop() {
      return { [rowKey]: rowKeyData, data: rowData }
    },
  })

  useEffect(() => {
    if (draggable && trRef.current) {
      drag(trRef.current)
      drop(trRef.current)
    }
  }, [trRef, drag, drop, draggable])

  const renderCell = useCallback(
    (column, rowData, rowIndex, colIndex, level) => {
      let data = rowData[column.key]
      if (column.render) {
        data = column.render(data, rowData, rowIndex)
      }
      const isExpand = expandedKeys.includes(rowData[rowKey])
      return (
        <td key={column.key}>
          <div
            style={{
              paddingLeft: colIndex === 0 ? `${level * indentSize}px` : 0,
              display: 'flex',
              alignItems: 'center',
            }}>
            {colIndex === 0 &&
              rowData[childrenField] &&
              React.cloneElement(expandIcon(isExpand), {
                onClick: () => toggleParentExpand(rowData),
                className: `icon-expand`,
              })}
            {data}
          </div>
        </td>
      )
    },
    [
      indentSize,
      rowKey,
      expandIcon,
      toggleParentExpand,
      childrenField,
      expandedKeys,
    ]
  )

  return (
    <Tr
      key={rowKeyData}
      ref={trRef}
      indentSize={indentSize}
      style={{
        opacity: isDragging ? 0.5 : 1,
        cursor: isDragging ? 'grabbing' : 'default',
        background: active ? '#efefef' : '',
      }}>
      {columns.map((c, colIndex) =>
        renderCell(c, rowData, rowIndex, colIndex, level)
      )}
    </Tr>
  )
}
