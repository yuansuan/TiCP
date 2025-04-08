import * as React from 'react'
import { observer, disposeOnUnmount } from 'mobx-react'
import { Table, Button } from '@/components'
import { Pagination, Input, DatePicker, Select } from 'antd'
import { LogList } from '@/domain/AuditLog'
import { observable, reaction, computed } from 'mobx'
import { ListWrapper } from './style'
import { untilDestroyed } from '@/utils/operators'
import { fromEvent } from 'rxjs'
import { debounceTime, startWith } from 'rxjs/operators'
import { exportFile } from './exportFile'
import { LogTimeline } from './LogTimeline'
import moment from 'moment'
import { replaceCustomTags } from '@/utils/formatter'
import SysConfig from '@/domain/SysConfig'

const RangePicker = DatePicker.RangePicker
const Option = Select.Option

interface IProps {
  logList: LogList
}
@observer
export class List extends React.Component<IProps> {
  @observable loading = false
  @observable dates = []
  @observable mode = false // false: table, true: timeline

  get options() {
    const {
      size,
      index,
      user_name,
      ip_address,
      operate_type,
      start_time,
      end_time
    } = this.props.logList

    return {
      size,
      index,
      user_name,
      ip_address,
      operate_type,
      operate_time: {
        start_time,
        end_time
      }
    }
  }

  get columns(): any[] {
    return [
      {
        props: {
          width: 120,
          fixed: 'left',
          resizable: true
        },
        header: '用户名称',
        cell: {
          props: {
            dataKey: 'user_name'
          }
        }
      },
      {
        props: {
          width: 150,
          fixed: 'left',
          resizable: true
        },
        header: 'IP 地址',
        cell: {
          props: {
            dataKey: 'ip_address'
          }
        }
      },
      {
        props: {
          width: 250,
          fixed: 'left',
          resizable: true
        },
        header: '操作时间',
        cell: {
          props: {
            dataKey: 'operate_time_str'
          }
        }
      },
      {
        props: {
          width: 120,
          fixed: 'left',
          resizable: true
        },
        header: '操作类型',
        cell: {
          props: {
            dataKey: 'operate_type'
          },
          render: ({ rowData }) => {
            return rowData.operate_type || ''
          }
        }
      },
      {
        props: {
          width: 800,
          resizable: true
        },
        header: '操作内容',
        cell: {
          props: {
            dataKey: 'operate_content'
          },
          render: ({ rowData }) => {
            const content = replaceCustomTags(rowData.operate_content)
            return <p title={content}>{content}</p>
          }
        }
      }
    ]
  }

  get dataSource() {
    return [...this.props.logList]
  }

  @observable width
  @observable height = 700

  tableContainerRef = null

  @computed
  get iconMap() {
    const { logList } = this.props
    // TODO for time line
    return {}
  }

  async componentDidMount() {
    // await this.props.logList.fetchOptTypes()
    this.refresh()

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

    setTimeout(() => {
      this.tableContainerRef.style.paddingRight = 1 + 'px'
    }, 3000)
  }

  refresh = async () => {
    this.loading = true
    try {
      await this.props.logList.fetch()
    } finally {
      this.loading = false
    }
  }

  @disposeOnUnmount
  disposer = reaction(
    () => this.options,
    () => {
      this.refresh()
    }
  )

  exportLogFile = () => {
    const { logList } = this.props
    const { user_name, start_time, end_time, operate_type, ip_address } =
      logList

    exportFile(
      {
        user_name,
        operate_type: operate_type ? +operate_type : null,
        start_time,
        end_time,
        ip_address
      },
      '审计日志'
    )
  }

  render() {
    const { logList } = this.props

    return (
      <ListWrapper>
        <div className='actions'>
          <div className='item'>
            <span className='label'>用户名称:</span>
            <Input
              placeholder='用户名称'
              value={logList.user_name}
              onChange={e => {
                logList.updateUsername(e.target.value)
              }}
            />
          </div>
          <div className='item'>
            <span className='label'>IP 地址:</span>
            <Input
              placeholder='IP 地址'
              value={logList.ip_address}
              onChange={e => {
                logList.updateIPAddress(e.target.value)
              }}
            />
          </div>
          <div className='item'>
            <span className='label'>操作类型:</span>
            <Select
              showSearch
              optionFilterProp='children'
              style={{ width: 160 }}
              value={logList.operate_type}
              onChange={value => {
                logList.updateOptType(value)
              }}>
              <Option value=''>全部</Option>
              {Object.entries(logList.LogOptionsTypesMap)
                .filter(([key, _]) =>
                  SysConfig.enableThreeMemberMgr ? true : key != '10'
                )
                .map(([key, name]) => (
                  <Option value={key} key={name}>
                    {logList.LogOptionsTypesMap[key]}
                  </Option>
                ))}
            </Select>
          </div>
          <div className='item'>
            <span className='label'>操作时间:</span>
            <RangePicker
              ranges={{
                最近24小时: [moment().subtract(1, 'days'), moment()],
                最近7天: [moment().subtract(7, 'days'), moment()],
                最近30天: [moment().subtract(30, 'days'), moment()]
              }}
              // @ts-ignore
              value={this.dates}
              showTime={{ format: 'HH:mm' }}
              style={{ width: 450 }}
              format='YYYY-MM-DD HH:mm'
              placeholder={['操作时间范围（开始）', '操作时间范围（结束）']}
              onChange={value => {
                // for clear
                if (!value) {
                  this.dates = null
                  logList.updateOptTime(null)
                  return
                }

                if (value.length === 0 || value.length === 2) {
                  this.dates = value
                  if (value.length === 0) {
                    logList.updateOptTime(null)
                  }

                  logList.updateOptTime(
                    this.dates.map(m => m.format('YYYY-MM-DD HH:mm:ss'))
                  )
                }
              }}
              allowClear={true}
            />
          </div>
          <div className='item'>
            <Button
              type='link'
              onClick={this.exportLogFile}
              disabled={this.dataSource.length === 0}>
              导出日志
            </Button>
          </div>
          {/* <div className='item'>
            <Switch
              checkedChildren='时间线'
              unCheckedChildren='表格'
              checked={this.mode}
              onChange={val => (this.mode = val)}
            />
          </div> */}
        </div>
        <div className='list' ref={ref => (this.tableContainerRef = ref)}>
          {this.mode ? (
            <div
              className={'logTimeline'}
              style={{ height: this.height, overflow: 'auto' }}>
              <LogTimeline
                logs={this.dataSource}
                iconMap={this.iconMap as any}
              />
            </div>
          ) : (
            <Table
              columns={this.columns}
              props={{
                data: this.dataSource,
                rowKey: 'id',
                height: this.height,
                loading: this.loading,
                locale: {
                  emptyMessage: '没有数据',
                  loading: '数据加载中...'
                }
              }}
            />
          )}
        </div>
        <div className='footer'>
          <Pagination
            showSizeChanger
            pageSize={logList.size}
            current={logList.index}
            total={logList.totals}
            onChange={logList.updateIndex.bind(logList)}
            onShowSizeChange={logList.updateSize.bind(logList)}
          />
        </div>
      </ListWrapper>
    )
  }
}
