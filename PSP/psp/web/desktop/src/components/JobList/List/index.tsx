import * as React from 'react'
import { inject, observer } from 'mobx-react'
import { Table } from '@/components'
import { EllipsisWrapper, JobDetail } from '@/components'
import { Wrapper } from './style'
import { getJobStatus } from '../common'
import { history, Http } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { observable } from 'mobx'
import { Modal, Icon } from '@/components'
import { JobTableType } from '@/domain/JobList'
import { currentUser } from '@/domain'
import { fromEvent } from 'rxjs'
import { debounceTime, startWith } from 'rxjs/operators'
import { downloadFile } from '@/utils/FileDownload'
interface ListProps {
  workspaceId?: number
  isSubJob?: boolean
  loading?: boolean
  isHistory?: boolean
  jobList: any[]
  isPopupJobDetail: boolean // 是否用弹窗显示 JobDetail
  columns: string[]
  hasCheckbox: boolean
  selectedRowKeys?: string[]
  reSubmit?: boolean
  updateSelectedRowKeys?: (key: string[]) => void
  updateOrder?: (orderBy: string, orderAsc: boolean) => void
  updateCurrentIndex?: (currentIndex: number) => void
  filterByType?: (type, value, currentPage: number) => void
}

const mapDataKeyToSortKey = {
  id: 'id',
  jobName: 'name',
  submitTime: 'create_time',
  startTime: 'start_time',
  endTime: 'end_time'
}

@inject((stores: any) => {
  const { selectedRowKeys, updateOrder, updateSelectedRowKeys, filterByType } =
    stores.store
  return {
    selectedRowKeys,
    updateOrder,
    updateSelectedRowKeys,
    filterByType
  }
})
@observer
export default class List extends React.Component<ListProps> {
  columnSorter = ({ sortType, sortKey }) => {
    if (sortType === '') {
      this.props.updateOrder('', undefined)
    } else {
      this.props.updateOrder(mapDataKeyToSortKey[sortKey], sortType === 'asc')
    }
  }

  // feature: "monitor"
  startVNCMonitor = async cwd => {
    return Http.get('/vnc/monitor', {
      params: {
        cwd,
        username: currentUser.name
      }
    })
  }

  getJobCWD = async rowData => {
    const { isHistory } = this.props
    // 获取当前 job 的 cwd (job_dir)
    const res = await Http.get('/job/detail', {
      params: {
        job_id: rowData.id,
        job_array_index: rowData.arrayIndex,
        job_table_type: isHistory
          ? JobTableType.FinishedTable
          : JobTableType.LiveTable
      }
    })

    const cwd = res.data.job.job_dir

    return cwd
  }

  downloadDCVFile = async rowData => {
    const cwd = await this.getJobCWD(rowData)
    const filePath = `${cwd}/.spooler_action/console.vnc`
    downloadFile([filePath], currentUser.id)
  }

  showJobDetails = rowData => {
    const {
      isSubJob = false,
      isPopupJobDetail,
      isHistory,
      reSubmit
    } = this.props

    if (isPopupJobDetail) {
      // open Modal
      Modal.show({
        title: '作业详情',
        footer: null,
        content: ({ onCancel, onOk }) => (
          <JobDetail
            jobId={rowData.id}
            reSubmitCallback={onOk}
            isHistory={isHistory}
            isSubJob={isSubJob}
            reSubmit={reSubmit}
            jobType={rowData.jobType}
          />
        ),
        width: 1200
      })
    } else {
      // workspaceId === -1 ? 'clusterjob' : 'job'
      history.push(
        `/${'job'}/${
          rowData.id
        }?isHistory=${isHistory}&isSubJob=${isSubJob}&hasReSubmit=${reSubmit}&jobType=${
          rowData.jobType
        }`
      )
    }
  }

  get columns(): any[] {
    const columnMap = {
      id: {
        header: '作业ID',
        dataKey: 'id',
        props: {
          width: 120,
          fixed: 'left',
          resizable: true
        },
        cell: {
          render: ({ rowData }) => (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                marginRight: 10
              }}>
              <a
                style={{ color: '#1458E0' }}
                onClick={() => this.showJobDetails(rowData)}>
                <EllipsisWrapper>{rowData.id}</EllipsisWrapper>
              </a>
            </div>
          )
        },
        sorter: this.columnSorter
      },
      jobName: {
        header: '作业名称',
        dataKey: 'jobName',
        props: {
          width: 200,
          fixed: 'left',
          resizable: true
        },
        cell: {
          render: ({ rowData }) => (
            <>
              <EllipsisWrapper
                title={`${
                  rowData.destClusterName &&
                  (/^\d+$/g.test(rowData.id) ? '一键爆发 ' : '泛超算云 ')
                }${rowData.jobName}`}>
                {rowData.destClusterName && (
                  <Icon
                    style={{ color: '#1A6EBA', padding: '0 10px 0 0' }}
                    type='cloud1'
                  />
                )}
                {rowData.jobName}
              </EllipsisWrapper>
            </>
          )
        },
        sorter: this.columnSorter
      },
      jobStatus: {
        header: '状态',
        dataKey: 'jobStatus',
        props: {
          width: 160,
          fixed: 'left',
          resizable: false
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  height: 56
                }}>
                {getJobStatus(rowData)}
              </div>
            )
          }
        }
      },
      app: {
        header: '应用程序',
        dataKey: 'app',
        props: {
          width: 160,
          fixed: 'left',
          resizable: true
        }
      },
      submitTime: {
        header: '提交时间',
        dataKey: 'submitTime',
        props: {
          width: 200,
          resizable: true
        },
        sorter: this.columnSorter
      },
      startTime: {
        header: '开始时间',
        dataKey: 'startTime',
        props: {
          width: 200,
          resizable: true
        },
        sorter: this.columnSorter
      },
      endTime: {
        header: '结束时间',
        dataKey: 'endTime',
        props: {
          width: 200,
          resizable: true
        },
        sorter: this.columnSorter
      },
      userName: {
        header: '提交用户',
        dataKey: 'userName',
        props: {
          width: 160,
          resizable: true
        }
      },
      queue: {
        header: '队列',
        dataKey: 'queue',
        props: {
          width: 160,
          resizable: true
        }
      },
      exHostStr: {
        header: '执行机器',
        dataKey: 'exHostStr',
        props: {
          width: 160,
          resizable: true
        }
      },
      numExecProcs: {
        header: '使用核数',
        dataKey: 'numExecProcs',
        props: {
          width: 160,
          resizable: true
        }
      },
      cpuTime: {
        header: '核时',
        dataKey: 'cpuTime',
        props: {
          width: 200,
          resizable: true
        }
      }
    }

    return this.props.columns.map(key => columnMap[key])
  }

  @observable width
  @observable height

  tableContainerRef = null

  componentDidMount() {
    fromEvent(window, 'resize')
      .pipe(untilDestroyed(this), startWith(''), debounceTime(300))
      .subscribe(() => {
        if (this.tableContainerRef) {
          // hack: wait container to render
          setTimeout(() => {
            this.width = this.tableContainerRef.clientWidth
            this.height = this.tableContainerRef.clientHeight
          }, 0)
        }
      })
  }

  get dataSource() {
    const { jobList } = this.props
    // format data
    return jobList
  }

  render() {
    const { selectedRowKeys, loading, hasCheckbox } = this.props
    const rowSelection = hasCheckbox
      ? {
          selectedRowKeys,
          onSelect: this.onSelect,
          onSelectAll: this.onSelectAll,
          onSelectInvert: this.onSelectInvert,
          props: {
            fixed: 'left'
          }
        }
      : undefined
    return (
      <Wrapper ref={ref => (this.tableContainerRef = ref)}>
        <Table
          columns={this.columns}
          props={{
            data: this.dataSource,
            rowKey: 'id',
            height: this.height,
            headerHeight: 54,
            rowHeight: 58,
            loading,
            virtualized: true, //添加虚拟滚动
            locale: {
              emptyMessage: '没有作业数据',
              loading: '数据加载中...'
            }
          }}
          rowSelection={rowSelection as any}
        />
      </Wrapper>
    )
  }
  private onSelectAll = keys => {
    this.props.updateSelectedRowKeys(keys)
  }

  private onSelectInvert = () => {
    this.props.updateSelectedRowKeys([])
  }

  private onSelect = (rowKey, checked) => {
    const { selectedRowKeys } = this.props
    let curKeys = [...selectedRowKeys]
    if (checked) {
      curKeys.push(rowKey)
    } else {
      curKeys = curKeys.filter(item => item !== rowKey)
    }
    this.props.updateSelectedRowKeys([...curKeys])
  }
}
