/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { TableProps as RTableProps } from 'rsuite-table/lib/Table'
import { ColumnProps as RTableColumnProps } from 'rsuite-table/lib/Column'
import { CellProps } from 'rsuite-table/lib/Cell'

export type INode = {
  [key: string]: any
  [key: number]: any
  key: number | string
  name: number | string
}

export type YSColumnProps = Omit<RTableColumnProps, 'width'> & {
  width?: number | string
}

export type ColumnProps = {
  props?: YSColumnProps
  headerClassName?: string
  header: string | (() => React.ReactElement<any>) | React.ReactElement<any>
  // if header is not a string, use this prop to map dataKey
  name?: string
  cell?: {
    props?: Partial<CellProps>
    render?: (props?: any) => React.ReactElement<any> | string
  }
  dataKey?: string

  sorter?: (option?: { sortType: string; sortKey: string }) => void
  filter?: {
    items: Array<INode>
    selectedKeys?: string[]
    updateSelectedKeys?: (keys: string[]) => void
    onChange?: (
      selectedKeys: string[],
      info?: { node?: INode; checked: boolean }
    ) => void
    searchable?: boolean
  }
}

export type PluginProps = {
  forceUpdate: Function
  hooks: {
    init: any
    beforeRender: any
    useLayoutEffect: any
    useEffect: any
  }
}

export type PluginContext = TableProps

export type PluginType =
  | ((props: PluginProps, options: any) => void)
  | {
      apply: (props: PluginProps, options: any) => any
    }

export type TableProps = {
  // TableId is the key to storage custom settings, use it to activate persistence
  // caveat: be care of duplicated
  tableId?: string
  // the defaults settings
  defaultConfig?: Array<
    { key: string; active?: boolean } | string | [string, boolean]
  >
  props: Omit<Partial<RTableProps>, 'rowKey'> & {
    rowKey: string | number
  }
  columns: Array<ColumnProps>
  // custom plugins
  plugins?: Array<PluginType>

  rowSelection?: Partial<{
    props: YSColumnProps
    defaultSelectedKeys: (string | number)[]
    selectedKeys: (string | number)[]
    selectedRowKeys: (string | number)[]
    onSelect: (rowKey?: string | number, checked?: boolean) => void
    onSelectAll: (keys?: (string | number)[]) => void
    onSelectInvert: () => void
    onChange: (kyes: (string | number)[]) => void
  }>

  onRow?: (
    rowData: object,
    rowIndex: number
  ) => Partial<{
    onClick: (rowData: object, rowIndex: number) => void
    onContextMenu: (rowData: object, rowIndex: number) => void
    onMouseEnter: (rowData: object, rowIndex: number) => void
    onMouseLeave: (rowData: object, rowIndex: number) => void
    onDoubleClick: (rowData: object, rowIndex: number) => void
  }>
}
