import * as React from 'react'
import { observer, disposeOnUnmount } from 'mobx-react'
import { Table, Modal } from '@/components'
import { Pagination, Input, DatePicker, Select, Popover, Button, message } from 'antd'
import { allApproveList, unapproveList, approvedList, statusColorMap } from '@/domain/Approve'
import { observable, reaction } from 'mobx'
import { ListWrapper } from './style'
import { untilDestroyed } from '@/utils/operators'
import { fromEvent } from 'rxjs'
import { debounceTime, startWith } from 'rxjs/operators'
import moment from 'moment'
import { ApproveTimeline } from '@/components/ApproveTimeline'
import { Form } from './Form'
import { currentUser } from '@/domain'

const RangePicker = DatePicker.RangePicker
const Option = Select.Option

const approveList = {
  'all': allApproveList, 
  'unapproved': unapproveList, 
  'approved': approvedList
}

@observer
export default class AllApproveList extends React.Component<any> {
  @observable allApproveList = approveList[this.props.type]
  @observable dates = []

  get options() {
    return {
      type: this.allApproveList.type,
      application_name: this.allApproveList.application_name, // 申请人
      // submitter_name: approveList.submitter_name, // 申请时间
      result: this.allApproveList.result,
      start_time: this.allApproveList.start_time,
      end_time: this.allApproveList.end_time,
      page_size: this.allApproveList.page_size,
      page_index: this.allApproveList.page_index,
    }
  }

  get columns(): any[] {
    return [
      {
        props: {
          width: 120,
          resizable: true,
        },
        header: '申请人',
        cell: {
          props: {
            dataKey: 'application_name',
          },
          render: ({ rowData }) => {
            return rowData.application_name || '--'
          },
        },
      },
      {
        props: {
          width: 220,
          resizable: true,
        },
        header: '申请时间',
        cell: {
          props: {
            dataKey: 'create_time_str',
          },
          render: ({ rowData }) => {
            return rowData.create_time_str
          },
        },
      },
      {
        props: {
          width: 160,
          resizable: true,
        },
        header: '申请操作类型',
        cell: {
          props: {
            dataKey: 'opt_type_str',
          },
          render: ({ rowData }) => {
            return rowData.opt_type_str
          },
        },
      },
      {
        props: {
          width: 800,
          resizable: true,
        },
        header: '申请操作内容',
        cell: {
          props: {
            dataKey: 'content',
          },
          render: ({ rowData }) => {
            return <p title={rowData.content}>{rowData.content}</p>
          },
        },
      },
      {
        props: {
          width: 120,
          fixed: 'right',
          resizable: true,
        },
        header: '审批人',
        cell: {
          props: {
            dataKey: 'approve_user_name',
          },
          render: ({ rowData }) => {
            return rowData.approve_user_name || '--'
          },
        },
      },
      this.props.type !== 'unapproved' ? {
        props: {
          width: 200,
          fixed: 'right',
          resizable: true,
        },
        header: '审批时间',
        cell: {
          props: {
            dataKey: 'approve_time_str',
          },
          render: ({ rowData }) => {
            return rowData.approve_time_str || '--'
          },
        },
      } : null,
      {
        props: {
          width: 120,
          fixed: 'right',
          resizable: true,
        },
        header: '审批结果',
        cell: {
          props: {
            dataKey: 'result_str',
          },
          render: ({ rowData }) => {
            const content = <ApproveTimeline approve={rowData} />

            return (
              <Popover
                placement='left'
                title={'申请信息'}
                content={content}
                trigger='click'>
                <Button type='link' style={{color: statusColorMap[rowData.status]}}>{rowData.result_str}</Button>
              </Popover>
            )
          },
        },
      },
      this.props.type !== 'approved' ? {
        props: {
          fixed: 'right',
          width: 200,
        },
        header: '操作',
        cell: {
          props: {
            dataKey: '',
          },
          render: ({ rowData }) => {
            return (
              <div style={{ paddingLeft: 10 }}>
                {
                  this.props.type === 'all' && <Button type='link' disabled={currentUser.name !== rowData.application_name || rowData.status !== 1} onClick={() => this.cancel(rowData)}>
                    撤销
                  </Button>
                }
                {
                  this.props.type === 'unapproved' && <Button type='link' onClick={() => this.accept(rowData)}>
                    同意
                  </Button>
                }
                {
                  this.props.type === 'unapproved' && <Button type='link' onClick={() => this.reject(rowData)}>
                    拒绝
                  </Button>
                }
              </div>
            )
          },
        },
      } : null,
    ].filter(Boolean)
  }

  cancel = rowData => {
    Modal.showConfirm({
      title: '撤销申请',
      content: (
        <>
          <p>{`确认撤销申请吗？`}</p>
          <p>{`申请内容: ${rowData.content}`}</p>
        </>
      ),
    }).then(() => {
      rowData.cancel('').then(res => {
        if (res.success) {
          message.success('撤销申请成功')
          this.refresh()
        } else {
          message.error(res.message || '撤销申请失败')
        }
      })
    })
  }

  accept = rowData => {
    Modal.showConfirm({
      title: '同意申请',
      content: (
        <>
          <p>{`确认同意 ${rowData.application_name} 的申请吗？`}</p>
          <p>{`申请内容: ${rowData.content}`}</p>
        </>
      ),
    }).then(() => {
      rowData.accept('').then(res => {
        if (res.success) {
          message.success('同意申请成功')
        } else {
          message.error(res.message || '同意申请失败')
        }
        this.refresh()
      }).catch(() => {
        this.refresh()
      })
    })
  }

  reject = rowData => {
    Modal.show({
      title: '拒绝申请',
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = content => {
          rowData.reject(content).then(res => {
            if (res.success) {
              message.success('拒绝申请成功')
              this.refresh()
              onOk()
            } else {
              message.error(res.message || '拒绝申请失败')
            }
          })
        }
        return (
          <Form
            rowData={rowData}
            onOk={ok}
            onCancel={onCancel}
          />
        )
      },
      width: 600,
    })
  }


  get dataSource() {
    return this.allApproveList.list
  }

  @observable width
  @observable height

  tableContainerRef = null

  componentDidMount() {
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

  refresh = () => {
    this.allApproveList.fetch()
  }

  @disposeOnUnmount
  disposer = reaction(
    () => this.options,
    () => {
      this.refresh()
    }
  )

  render() {
    return (
      <ListWrapper>
        <div className='actions'>
          <div className='filterArea'>
            {this.props.type !== 'all' && <div className='item'>
              <span className='label'>申请人:</span>
              <Input
                style={{ width: 130 }}
                placeholder='申请人'
                value={this.allApproveList.application_name}
                onChange={e => {
                  this.allApproveList.updateApplicationName(e.target?.value?.trim())
                }}
              />
            </div>}
            <div className='item'>
              <span className='label'>申请操作类型:</span>
              <Select
                style={{ width: 160 }}
                showSearch
                optionFilterProp='children'
                value={this.allApproveList.type}
                onChange={value => {
                  this.allApproveList.updateOptType(value)
                }}>
                <Option value=''>全部</Option>
                {this.allApproveList.optTypes.map(t => (
                  <Option key={t.approve_type} value={t.approve_type}>
                    {t.name}
                  </Option>
                ))}
              </Select>
            </div>
            <div className='item'>
              <span className='label'>申请时间:</span>
              <RangePicker
                ranges={{
                  最近24小时: [moment().subtract(1, 'days'), moment()],
                  最近7天: [moment().subtract(7, 'days'), moment()],
                  最近30天: [moment().subtract(30, 'days'), moment()],
                }}
                // @ts-ignore
                value={this.dates}
                showTime={{ format: 'HH:mm' }}
                format='YYYY-MM-DD HH:mm'
                placeholder={['申请时间范围（开始）', '申请时间范围（结束）']}
                onChange={value => {
                   // for clear
                   if (!value) {
                    this.dates = null
                    this.allApproveList.updateTime(
                      null
                    )
                    return
                  }

                  if (value?.length === 0 || value?.length === 2) {
                    this.dates = value
                    if (value.length === 0) {
                      this.allApproveList.updateTime(
                        null
                      )
                    }
                    this.allApproveList.updateTime(
                      this.dates.map(m => m.unix())
                    )    
                  } 
                }}
                allowClear={true}
              />
            </div>
            {/* <div className='item'>
              <span className='label'>审批人:</span>
              <Input
                style={{ width: 130 }}
                placeholder='审批人'
                value={this.allApproveList.approve_user_name}
                onChange={e => {
                  this.allApproveList.uodateApproveUserName(e.target.value)
                }}
              />
            </div> */}
            {this.props.type !== 'unapproved' && <div className='item'>
              <span className='label'>审批结果:</span>
              <Select
                style={{ width: 160 }}
                showSearch
                optionFilterProp='children'
                value={this.allApproveList.result}
                onChange={value => {
                  this.allApproveList.updateResult(value)
                }}>
                <Option value=''>全部</Option>
                {Object.entries(this.allApproveList.resultTypes[this.allApproveList.listType])
                  .map(([key, value])=> (
                  <Option key={key} value={+key}>
                    {value}
                  </Option>
                ))}
              </Select>
            </div>}
          </div>
          <div className='btnArea'>
            <div></div>
          </div>
        </div>
        <div className='list' ref={ref => (this.tableContainerRef = ref)}>
          <Table
            columns={this.columns}
            props={{
              data: this.dataSource,
              rowKey: 'id',
              height: this.height,
              locale: {
                emptyMessage: '没有数据',
                loading: '数据加载中...',
              },
            }}
          />
        </div>
        <div className='footer'>
          <Pagination
            showSizeChanger
            pageSize={this.allApproveList.page_size}
            current={this.allApproveList.page_index}
            total={this.allApproveList.page_ctx.total}
            onChange={this.allApproveList.updateCurrentIndex.bind(this.allApproveList)}
            onShowSizeChange={this.allApproveList.updatePageSize.bind(this.allApproveList)}
          />
        </div>
      </ListWrapper>
    )
  }
}
