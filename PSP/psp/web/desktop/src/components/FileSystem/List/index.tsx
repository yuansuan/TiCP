import * as React from 'react'
import { computed, observable, action } from 'mobx'
import { observer } from 'mobx-react'
import { debounceTime } from 'rxjs/operators'
import { Subject, Observable } from 'rxjs'
import { message, Tooltip } from 'antd'
import moment from 'moment'

import { RootPoint, PathHistory, List as FileList } from '@/domain/FileSystem'
import { EditableText, Modal } from '@/components'
import { Icon } from '@/components'
import { Table } from '@/components'
import {
  createMobxStream,
  formatFileSize,
  Validator,
  eventEmitter
} from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { EditAction, DownloadAction } from '../Actions'
import { isTextFile } from '../Actions/Edit'
import { StyledPanel } from './style'

interface IProps {
  selectedKeys$: Subject<any>
  keyword$: Subject<any>
  resize$: Observable<any>
  fetching: boolean
  fileList: FileList
  currentPoint: RootPoint
  points: RootPoint[]
  history: PathHistory
  hasPerm: boolean
}

const defaultConfig = [
  'name',
  'type',
  'size',
  { key: 'modifiedTime', active: false },
  'operators'
]

const EDITABLE_SIZE = 3 * 1024 * 1024

const isPreviewable = file => isTextFile(file)
const isEditable = file => {
  if (file.size > EDITABLE_SIZE) {
    return {
      editable: false,
      message: '文件大小超过 3M'
    }
  }

  if (!isTextFile(file)) {
    return {
      editable: false,
      message: '非文本文件'
    }
  }

  return {
    editable: true
  }
}

@observer
export default class Panel extends React.Component<IProps> {
  @observable width: number = 0
  @observable height: number = 0
  @observable selectedRowKeys: string[] = []
  @observable sortType = ''
  @observable sortColumn = ''
  @observable filterTypes = []
  @observable keyword = ''
  @observable refreshing = false
  @action
  updateWidth = width => (this.width = width)
  @action
  updateHeight = height => (this.height = height)
  @action
  updateSelectedRowKeys = keys => (this.selectedRowKeys = keys)
  @action
  updateSortType = type => (this.sortType = type)
  @action
  updateSortColumn = column => (this.sortColumn = column)
  @action
  updateFilterTypes = types => (this.filterTypes = types)
  @action
  updateKeyword = keyword => (this.keyword = keyword)
  @action
  updateRefreshing = refreshing => (this.refreshing = refreshing)

  panelRef = null
  tableRef = null
  componentDidMount() {
    //fix bug: T21760
    eventEmitter.on('FILE_SYSTEM_FILE_DELETE_EVNET', () => {
      if (this.tableRef) {
        this.tableRef.scrollTop(0)
      }
    })
    const { resize$, keyword$ } = this.props

    resize$.pipe(untilDestroyed(this)).subscribe(() => {
      if (this.panelRef) {
        // hack: wait panel to render
        setTimeout(() => {
          this.updateWidth(this.panelRef.clientWidth)
          this.updateHeight(this.panelRef.clientHeight - 50)
        }, 0)
      }
    })

    // monitor the keyword change
    keyword$
      .pipe(untilDestroyed(this), debounceTime(300))
      .subscribe(this.updateKeyword)

    // export the selectedRowKeys$
    createMobxStream(() => this.selectedRowKeys)
      .pipe(untilDestroyed(this))
      .subscribe(keys => {
        this.props.selectedKeys$.next(keys)
      })

    // monitor the change of fileList
    const { fileList } = this.props
    createMobxStream(() => fileList.children.map(item => item.path))
      .pipe(untilDestroyed(this))
      .subscribe(paths => {
        // diff keys with selectedRowKeys
        const finalKeys = this.selectedRowKeys.filter(key =>
          paths.includes(key)
        )
        // update the selectedRowKeys
        this.updateSelectedRowKeys(finalKeys)
      })

    // monitor the change of fileList.path
    createMobxStream(() => fileList.parentPath)
      .pipe(untilDestroyed(this))
      .subscribe(() => {
        // reset filter types
        this.updateFilterTypes([])
      })

    createMobxStream(() => this.visibleTypes)
      .pipe(untilDestroyed(this))
      .subscribe(types => {
        // reset filter types
        this.updateFilterTypes(
          this.filterTypes.filter(type => types.includes(type))
        )
      })
  }
  componentWillUnmount() {
    eventEmitter.off('FILE_SYSTEM_FILE_DELETE_EVNET')
  }

  @computed
  get visibleFiles() {
    const { fileList } = this.props
    const { keyword, filterTypes } = this

    return fileList.children.filter(item => {
      // match keyword
      return (
        (!keyword || item.name.includes(keyword)) &&
        // match filter type
        (filterTypes.length === 0 || filterTypes.includes(item.type))
      )
    })
  }

  @computed
  get visibleTypes() {
    const { fileList } = this.props
    return [...new Set(fileList.children.map(item => item.type))]
  }

  @computed
  get dataSource() {
    let data = [...this.visibleFiles]

    // sort
    if (this.sortColumn && this.sortType) {
      switch (this.sortColumn) {
        // sort by name
        case 'name': {
          data = data.sort((x, y) => {
            if (this.sortType === 'asc') {
              return x.name.localeCompare(y.name)
            } else {
              return y.name.localeCompare(x.name)
            }
          })
          break
        }
        // sort by size
        case 'size': {
          data = data.sort((x, y) => {
            if (this.sortType === 'asc') {
              return x.size - y.size
            } else {
              return y.size - x.size
            }
          })
          break
        }
        // sort by modifiedTime
        case 'modifiedTime': {
          data = data.sort((x, y) => {
            const xTime = new Date(x.modifiedTime).getTime()
            const yTime = new Date(y.modifiedTime).getTime()

            return this.sortType === 'asc' ? xTime - yTime : yTime - xTime
          })
          break
        }
      }
    }

    return data.map(item => ({
      id: item.id,
      isFile: item.isFile,
      isText: item.is_text,
      name: item.name || '--',
      size: Number.isNaN(parseFloat(item.size)) ? '--' : item.size,
      type: item.type || '--',
      path: item.path ? item.path : '--',
      modifiedTime: item.modifiedTime || '--',
      isSymLink: item.is_sym_link
    }))
  }

  @computed
  get statisticsInfo() {
    const fileNum = this.dataSource.filter(item => item.isFile).length
    return `当前目录包含 ${
      this.dataSource.length - fileNum
    } 个文件夹，${fileNum} 个文件`
  }

  private refresh = async () => {
    const {
      currentPoint,
      fileList: { parentPath }
    } = this.props
    try {
      this.updateRefreshing(true)
      await currentPoint.service.fetch(parentPath)
    } finally {
      this.updateRefreshing(false)
    }
  }

  private onSelectAll = keys => {
    this.updateSelectedRowKeys(keys)
  }

  private onSelectInvert = () => {
    this.updateSelectedRowKeys([])
  }

  private onSelect = (rowKey, checked) => {
    let keys = this.selectedRowKeys

    if (checked) {
      keys = [...keys, rowKey]
    } else {
      const index = keys.findIndex(item => item === rowKey)
      keys.splice(index, 1)
      // create new ref to activate observer
      keys = [...keys]
    }

    this.updateSelectedRowKeys(keys)
  }

  private beforeRename = value => {
    const { error } = Validator.filename(value)
    if (error) {
      message.error(`重命名失败：${error.message}`)
      return false
    }

    return true
  }

  private onRename = ({ name, path }) => {
    const { fileList, currentPoint: point } = this.props
    const targetNode = fileList.filterFirstNode(item => item.path === path)

    if (!name || name === targetNode.name) {
      return
    }

    if (targetNode) {
      point.service
        .rename({
          path: targetNode.path,
          newName: name
        })
        .then(() => {
          message.success('重命名成功')
        })
    }
  }

  private beforeDownload = file => {
    if (!this.props.hasPerm) {
      return Promise.reject()
    }

    let content = ''
    if (!isTextFile(file)) {
      content = '非文本文件无法预览，是否下载？'
    }

    if (content) {
      return Modal.showConfirm({
        content
      })
    } else {
      return true
    }
  }

  render() {
    const { fetching, currentPoint: point } = this.props
    const { width, height, selectedRowKeys, refreshing } = this
    const columns: any[] = [
      {
        props: {
          width: width * 0.2,
          resizable: true
        },
        header: '名称',
        sorter: ({ sortType, sortKey }) => {
          this.updateSortType(sortType)
          this.updateSortColumn(sortKey)
        },
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey, rowKey }) => {
            const { history, hasPerm } = this.props
            const getText = value => (
              <>
                <span style={{ color: 'black', marginRight: 3 }}>
                  {rowData.isFile ? (
                    rowData.isSymLink ? (
                      <Icon type='link-file' />
                    ) : (
                      <Icon type='file' />
                    )
                  ) : rowData.isSymLink ? (
                    <Icon type='link-folder' />
                  ) : (
                    <Icon type='folder' />
                  )}
                </span>
                <span title={value}>{value.replace(/ /g, '\u00a0')}</span>
              </>
            )

            const createReadOnlyText = (onClick?: any) => (
              <div style={{ cursor: 'pointer' }} onClick={onClick}>
                {getText(rowData[dataKey])}
              </div>
            )

            const createEditableText = (onClick?: any) => (
              <EditableText
                Text={getText}
                EditIcon={
                  <Tooltip title='重命名'>
                    <Icon type='edit-filled' />
                  </Tooltip>
                }
                onClick={onClick}
                defaultValue={rowData[dataKey]}
                defaultShowEdit={false}
                beforeConfirm={this.beforeRename}
                onConfirm={name => this.onRename({ name, path: rowKey })}
              />
            )
            return (
              <div className='nameCell'>
                {isPreviewable(rowData) ? (
                  <EditAction
                    title={rowData[dataKey]}
                    path={rowData.path}
                    point={point}
                    readOnly={true}>
                    {hasPerm ? createEditableText() : createReadOnlyText()}
                  </EditAction>
                ) : rowData.isFile ? (
                  <DownloadAction
                    point={point}
                    targets={[rowData.path]}
                    beforeDownload={() => this.beforeDownload(rowData)}>
                    {hasPerm ? createEditableText() : createReadOnlyText()}
                  </DownloadAction>
                ) : (
                  <>
                    {hasPerm
                      ? createEditableText(() =>
                          history.push({
                            source: point,
                            path: rowData.path
                          })
                        )
                      : createReadOnlyText(() =>
                          history.push({
                            source: point,
                            path: rowData.path
                          })
                        )}
                  </>
                )}
              </div>
            )
          }
        }
      },
      {
        props: {
          width: width * 0.2,
          resizable: true
        },
        header: '大小',
        sorter: ({ sortType, sortKey }) => {
          this.updateSortType(sortType)
          this.updateSortColumn(sortKey)
        },
        cell: {
          props: {
            dataKey: 'size'
          },
          render: ({ rowData, dataKey }) => {
            if (rowData.isFile) {
              return formatFileSize(rowData[dataKey])
            } else {
              return '--'
            }
          }
        }
      },
      {
        props: {
          width: width * 0.2,
          resizable: true
        },
        header: '类型',
        filter: {
          selectedKeys: this.filterTypes,
          updateSelectedKeys: this.updateFilterTypes,
          items: this.visibleTypes.map(type => ({
            key: type,
            name: type || '--'
          }))
        },
        cell: {
          props: {
            dataKey: 'type'
          },
          render: ({ rowData, dataKey }) => (
            <div className='typeCell' title={rowData[dataKey]}>
              {rowData[dataKey]}
            </div>
          )
        }
      },
      {
        props: {
          width: width * 0.2,
          resizable: true
        },
        header: '修改时间',
        sorter: ({ sortType, sortKey }) => {
          this.updateSortType(sortType)
          this.updateSortColumn(sortKey)
        },
        cell: {
          props: {
            dataKey: 'modifiedTime'
          },
          render: ({ rowData, dataKey }) =>
            moment.unix(rowData[dataKey]).format('YYYY/MM/DD HH:mm:ss')
        }
      },
      {
        props: {
          width: 90
        },
        header: '操作',
        cell: {
          props: {
            dataKey: 'operators'
          },
          render: ({ rowData }) => {
            const { editable, message } = isEditable(rowData)
            return (
              <div className='operators'>
                {editable && this.props.hasPerm ? (
                  <EditAction
                    title={rowData.name}
                    path={rowData.path}
                    point={point}
                    readOnly={false}>
                    <span className='item'>编辑</span>
                  </EditAction>
                ) : message ? (
                  <Tooltip title={message}>
                    <span className='disabled'>编辑</span>
                  </Tooltip>
                ) : (
                  <span className='disabled'>编辑</span>
                )}
              </div>
            )
          }
        }
      }
    ]

    return (
      <StyledPanel ref={ref => (this.panelRef = ref)}>
        <div className='list'>
          <Table
            tableId='FileManagement'
            props={
              {
                ref: ref => (this.tableRef = ref),
                loading: fetching || refreshing,
                height,
                data: this.dataSource,
                rowKey: 'path',
                rowClassName: 'tableRow',
                virtualized: true //添加虚拟滚动
              } as any
            }
            defaultConfig={defaultConfig}
            columns={columns}
            rowSelection={{
              selectedRowKeys,
              onSelect: this.onSelect,
              onSelectAll: this.onSelectAll,
              onSelectInvert: this.onSelectInvert
            }}
          />
        </div>
        <div className='footer'>
          <div className='right'>
            <span className='info'>{this.statisticsInfo}</span>
            <div className='refresh' onClick={this.refresh}>
              <Icon type='refresh' />
            </div>
          </div>
        </div>
      </StyledPanel>
    )
  }
}
