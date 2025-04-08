/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { ReactElement } from 'react'

export type expandIcon = (isExpand: boolean) => ReactElement
export type onDragEndFn = (key1: string, key2: string) => void

export interface IProps {
  dataSource: IData[] // 表格数据
  columns: IColumn[] // 列数据
  rowKey?: string // 每一行数据的key
  childrenField?: string // children的字段, 默认为'children'
  expandIcon?: expandIcon // 展开/收起图标
  expandedKeys?: string[] // 受控展开的key, 值为rowKey字段的值数组
  defaultExpandAll?: boolean // 是否默认展开全部, 当传入expandedKeys时无效
  indentSize?: number // 每级的缩进
  onExpand?: (expanded: string[], record: IData) => any // row被展开的回调
  draggable?: boolean
  onDragEnd?: onDragEndFn // 拖拽完成后回调
}

export interface IColumn {
  title: string // 列标题
  key: string // 取data里的哪个属性
  width?: number // 列宽
  render?: (text?: any, record?: IData, index?: number) => ReactElement | string
}

export interface IData {
  children?: IData[]
  [key: string]: any
}
