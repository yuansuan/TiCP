import * as React from 'react'
import { computed, observable, action } from 'mobx'
import { observer } from 'mobx-react'
import { Pagination, message } from 'antd'
import { Search } from '@/components'
import { Table, Button, Modal } from '@/components'

import { Wrapper, TopWrapper } from './style'
import { TABLE_CONF } from './const'
import { formatNodeAttrNumber } from '@/utils/formatter'
import { nodeManager, NodeActionLabel } from '@/domain/NodeMG'
import { PSPCanCloseStatus, PSPCanOpenStatus } from '../commonStatus'

const round = formatNodeAttrNumber

interface IProps {
  loading: boolean
  data: any[]
  onRowClick?: (rowData: any) => void
  context?: any
}

@observer
export default class NodeTable extends React.Component<IProps> {
  wrapperRef = null
  resizeObserver = null

  constructor(props) {
    super(props)
    this.wrapperRef = React.createRef()
  }

  @observable selectedRowKeys = []
  @observable sortKey = 'scheduler_status'
  @observable sortType = true
  @observable interval = null
  @observable height = 400
  @observable width = 800

  @action
  updateSelectedRowKeys = keys => (this.selectedRowKeys = keys)

  componentDidMount() {
    this.interval = setInterval(() => nodeManager.getNodeList(), 5000)

    this.resizeObserver = new ResizeObserver(entries => {
      for (let entry of entries) {
        this.height = entry.contentRect.height
        this.width = entry.contentRect.width
      }
    })

    this.resizeObserver.observe(this.wrapperRef.current)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      this.wrapperRef.current.style.paddingRight = 1 + 'px'
    }, 3000)
  }

  componentWillUnmount(): void {
    this.interval && clearInterval(this.interval)
    this.resizeObserver && this.resizeObserver.disconnect()
  }

  getDataProps = propName => {
    return [...new Set(this.props.data.map(item => item[propName]))].map(
      value => ({
        key: value,
        name: value
      })
    )
  }

  @computed
  get scheduler_states() {
    return this.getDataProps('scheduler_status')
  }

  @computed
  get node_status() {
    return this.getDataProps('node_status')
  }

  @computed
  get node_type() {
    return this.getDataProps('node_type')
  }

  @computed
  get resource_attr() {
    return this.getDataProps('resource_attr')
  }

  columnSorter = ({ sortType, sortKey }) => {
    if (sortType === '') {
      this.sortKey = 'node_name'
      this.sortType = true
    } else {
      this.sortKey = sortKey
      this.sortType = sortType === 'asc'
    }
  }

  @computed
  get columns(): any {
    const cols = [
      {
        props: {
          resizable: true,
          width: 200
        },
        // sorter: this.columnSorter,
        header: '机器名',
        cell: {
          props: {
            dataKey: 'node_name'
          },
          render: ({ rowData }) => {
            const { onRowClick } = this.props
            return (
              <a
                title={rowData.node_name}
                // style={{ color: '#1458E0' }}
                style={{ color: 'rgba(0,0,0,0.65)' }}
                onClick={() => {
                  onRowClick && onRowClick({ rowData })
                }}>
                {rowData.node_name}
              </a>
            )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '调度器状态',
        // filter: {
        //   items: this.scheduler_states,
        //   onChange: keys => {
        //     this.onDataPropChange('scheduler_status', keys)
        //   }
        // },
        dataKey: 'scheduler_status'
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '队列名称',
        dataKey: 'queue_name'
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        // filter: {
        //   items: this.node_type,
        //   onChange: keys => {
        //     this.onDataPropChange('node_type', keys)
        //   }
        // },
        header: '机器类型',
        // sorter: this.columnSorter,
        dataKey: 'node_type'
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '总核数',
        // sorter: this.columnSorter,
        dataKey: 'total_core_num'
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '使用核数',
        // sorter: this.columnSorter,
        dataKey: 'used_core_num'
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '空闲核数',
        // sorter: this.columnSorter,
        dataKey: 'free_core_num'
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '总内存',
        // sorter: this.columnSorter,
        dataKey: 'total_mem',
        cell: {
          props: {
            dataKey: 'total_mem'
          },
          render: ({ rowData }) => round(+rowData.total_mem)
        }
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '使用内存',
        // sorter: this.columnSorter,
        dataKey: 'used_mem',
        cell: {
          props: {
            dataKey: 'used_mem'
          },
          render: ({ rowData }) => round(+rowData.used_mem)
        }
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '空闲内存',
        // sorter: this.columnSorter,
        dataKey: 'free_mem',
        cell: {
          props: {
            dataKey: 'free_mem'
          },
          render: ({ rowData }) => round(+rowData.free_mem)
        }
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '可用内存',
        // sorter: this.columnSorter,
        dataKey: 'available_mem',
        cell: {
          props: {
            dataKey: 'available_mem'
          },
          render: ({ rowData }) => round(+rowData.available_mem)
        }
      }
    ]

    return cols
  }

  @observable dataProps = {
    scheduler_status: [],
    node_status: [],
    node_type: []
  }

  @computed
  get tableData() {
    // return this.bySort(this.filterDataProps(this.props.data))
    return this.filterDataProps(this.props.data)
  }

  private onSearch = (value: string) => {
    nodeManager.pageIndex = 1
    nodeManager.nodeName = value

    nodeManager.getNodeList()
  }

  private bySort = list => {
    return list.sort((a, z) => {
      // 升
      if (this.sortType) {
        return a[this.sortKey] > z[this.sortKey]
          ? 1
          : a[this.sortKey] === z[this.sortKey]
          ? 0
          : -1
        // 降
      } else {
        return z[this.sortKey] > a[this.sortKey]
          ? 1
          : z[this.sortKey] === a[this.sortKey]
          ? 0
          : -1
      }
    })
  }

  private filterDataProps(list) {
    return Object.keys(this.dataProps).reduce((list, prop) => {
      if (this.dataProps[prop].length === 0) {
        return list
      } else {
        return list.filter(node => this.dataProps[prop].includes(node[prop]))
      }
    }, list)
  }

  private onDataPropChange = (prop, states) => {
    nodeManager.pageIndex = 1
    this.dataProps[prop] = states
  }

  private onPageChange = current => {
    nodeManager.pageIndex = current
    nodeManager.getNodeList()
  }

  private onPageSizeChange = (current, size) => {
    nodeManager.pageSize = size
    nodeManager.pageIndex = current
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
    }
    this.updateSelectedRowKeys(keys)
  }

  private open = () => {
    Modal.show({
      title: '接受作业',
      content: '确定改变所选机器为接受作业状态吗？',
      onOk: async () => {
        await nodeManager.operate(this.selectedRowKeys, NodeActionLabel.open)
        this.updateSelectedRowKeys([])
        message.success('操作正在进行，请稍后等待刷新', 5)
      }
    })
  }

  private close = () => {
    Modal.show({
      title: '拒绝作业',
      content: '确定改变所选机器为拒绝作业状态吗？',
      onOk: async () => {
        await nodeManager.operate(this.selectedRowKeys, NodeActionLabel.close)
        this.updateSelectedRowKeys([])
        message.success('操作正在进行，请稍后等待刷新', 5)
      }
    })
  }

  get isEnableOpen() {
    return PSPCanOpenStatus.some(n => this.isEnableBySchedulerStatus(n, false))
  }

  get isEnableClose() {
    return PSPCanCloseStatus.some(n => this.isEnableBySchedulerStatus(n, false))
  }

  isEnableBySchedulerStatus = (status, isFazzy) => {
    // 当用户未选择任何 row
    if (this.selectedRowKeys.length === 0) return false

    let selectedRows = this.tableData.filter(n =>
      this.selectedRowKeys.includes(n.node_name)
    )

    return selectedRows.some(n =>
      isFazzy
        ? n.status.includes(status)
        : n.status === status
    )
  }

  render() {
    return (
      <Wrapper ref={this.wrapperRef}>
        <TopWrapper>
          <div>
            <Button
              style={{ marginRight: 20 }}
              ghost
              type='primary'
              disabled={!this.isEnableOpen}
              onClick={this.open}>
              接受作业
            </Button>
            <Button
              ghost
              type='primary'
              disabled={!this.isEnableClose}
              onClick={this.close}>
              拒绝作业
            </Button>
          </div>
          <Search
            onSearch={this.onSearch}
            debounceWait={300}
            placeholder={'请输入机器名称'}
          />
        </TopWrapper>
        <Table
          tableId={TABLE_CONF.id}
          defaultConfig={TABLE_CONF.columns}
          columns={this.columns}
          props={{
            data: this.tableData,
            height: this.height,
            loading: this.props.loading,
            shouldUpdateScroll: false,
            locale: {
              emptyMessage: '没有节点数据',
              loading: '数据加载中...'
            },
            rowKey: 'node_name'
          }}
          rowSelection={{
            selectedRowKeys: this.selectedRowKeys,
            onSelect: this.onSelect,
            onSelectAll: this.onSelectAll,
            onSelectInvert: this.onSelectInvert,
            props: {
              fixed: 'left'
            }
          }}
        />
        <div style={{ textAlign: 'center', padding: 10 }}>
          <Pagination
            showSizeChanger
            onChange={this.onPageChange}
            current={nodeManager.pageIndex}
            total={nodeManager.total}
            onShowSizeChange={this.onPageSizeChange}
          />
        </div>
      </Wrapper>
    )
  }
}
