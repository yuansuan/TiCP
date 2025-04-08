import * as React from 'react'
import { computed, observable, action } from 'mobx'
import { observer } from 'mobx-react'
import { Pagination, message, DatePicker, Select, Input } from 'antd'
import { Table, Button, Modal } from '@/components'

import { ListWrapper, TopWrapper } from './style'
import { adminProjectMG, personalProjectMG, PROJECT_STATE_MAP, PROJECT_STATE_ENUM } from '@/domain/ProjectMG'
import { DatePicker_FORMAT, DatePicker_SHOWTIME_FORMAT, GeneralDatePickerRange } from '@/constant'
import { openUserSelector, onEdit, onAdd, onPreview, onChangeOwner } from './Form'
import { currentUser } from '@/domain'

const { RangePicker } = DatePicker
const { Option } = Select

interface IProps {
  onRowClick?: (rowData: any) => void
  context?: any
  isAdmin?: boolean
  isRefresh?: boolean
}

@observer
export default class List extends React.Component<IProps> {
  wrapperRef = null
  resizeObserver = null

  constructor(props) {
    super(props)
    this.wrapperRef = React.createRef()
  }

  @observable projectMG = this.props.isAdmin ? adminProjectMG : personalProjectMG
  @observable loading = true
  @observable selectedRowKeys = [] // TODO 支持多项，暂时保留改变量
  @observable height = 400
  @observable width = 800

  @action
  updateSelectedRowKeys = keys => (this.selectedRowKeys = keys)

  componentDidMount() {
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

    if (this.props.isRefresh) {
      this.projectMG.getList(this.props.isAdmin).finally(() => {
        this.loading = false
      })
    }
  }

  componentDidUpdate(prevProps: Readonly<IProps>, prevState: Readonly<{}>, snapshot?: any): void {
    if (prevProps.isRefresh !== this.props.isRefresh) {
      if (this.props.isRefresh) {
        this.loading = true
        this.projectMG.pageIndex = 1
        this.projectMG.getList(this.props.isAdmin).finally(() => {
          this.loading = false
        })
      }
    }
  }

  componentWillUnmount(): void {
    this.resizeObserver && this.resizeObserver.disconnect()
  }


  @computed
  get columns(): any {
    const cols = [
      {
        props: {
          // resizable: true,
          minWidth: 180,
          flexGrow: 1
        },
        header: '项目名称',
        cell: {
          props: {
            dataKey: 'project_name'
          },
          render: ({ rowData }) => {
            const { onRowClick } = this.props
            return (
              <a
                title={rowData.project_name}
                // style={{ color: '#1458E0' }}
                style={{ color: 'rgba(0,0,0,0.65)' }}
                onClick={() => {
                  onRowClick && onRowClick({ rowData })
                }}>
                {rowData.project_name}
              </a>
            )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 400
        },
        header: '项目周期(开始时间 - 结束时间)',
        cell: {
          props: {
            dataKey: 'start_time',
          },
          render: ({ rowData }) => {
            return (
              <span>
                {rowData['start_time']} - {rowData['end_time']}
              </span>
            )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 160
        },
        header: '管理员',
        dataKey: 'project_owner_name'
      },
      {
        props: {
          resizable: true,
          width: 120
        },
        header: '项目状态',
        cell: {
          props: {
            dataKey: 'state',
          },
          render: ({ rowData }) => {
            return (
              <span>
                {PROJECT_STATE_MAP[rowData.state] || '--'}
              </span>
            )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '创建时间',
        dataKey: 'create_time'
      },
      {
        props: {
          resizable: true,
          width: 320
        },
        header: '操作',
        cell: {
          props: {
            dataKey: 'opt'
          },
          render: ({ rowData }) => <>
            {
              (currentUser.hasProjectMgrPerm || currentUser.hasSysMgrPerm) ? (<>
                <Button type="link" onClick={() => this.viewProject(rowData)}>查看</Button>
                <Button type="link" disabled={this.isEditBtnDisabled(rowData)} onClick={() => this.editProject(rowData)}>编辑</Button>
                <Button type="link" onClick={() => this.memberMgr(rowData)}>成员管理</Button>
                <Button type="link" disabled={this.isDeleteBtnDisabled(rowData)} onClick={() => this.delete(rowData)}>删除</Button>
                <Button type="link" disabled={this.isTerminationBtnDisabled(rowData)} onClick={() => this.termination(rowData)}>终止</Button>
              </>) : (
                <Button type="link" onClick={() => this.viewProject(rowData)}>查看</Button>
              )
            } 
            {
              (currentUser.hasSysMgrPerm && this.props.isAdmin ) &&  <Button type="link" onClick={() => this.changeOwner(rowData)}>管理员转移</Button>
            }
          </>
        }
      }
    ]

    return cols
  }

  @computed
  get tableData() {
    return this.projectMG.list
  }

  private isAddBtnDisabled = () => {
    return currentUser.hasProjectMgrPerm || currentUser.hasSysMgrPerm ? false : '没有权限，只有项目管理员或拥有系统管理权限才可以操作' 
  }

  private isBtnDiabled = (rowDate) => {
    // 是该项目 owner 或 有系统管理权限
    if (currentUser.hasSysMgrPerm) {
      return false
    } else if (rowDate.isOwner) {
      return false
    } else {
      return '没有权限，只有项目所有者或拥有系统管理权限，才可以操作'
    }
  }

  private isEditBtnDisabled = (rowDate) => {
    const message = this.isBtnDiabled(rowDate)
    if (message) return message
    return rowDate.state !== PROJECT_STATE_ENUM.Init && rowDate.state !== PROJECT_STATE_ENUM.Running ? '只有项目状态为 初始化 和 已运行 时，才能编辑' : false
  }

  private isTerminationBtnDisabled = (rowDate) => {
    const message = this.isBtnDiabled(rowDate)
    if (message) return message
    return rowDate.state === PROJECT_STATE_ENUM.Completed || rowDate.state === PROJECT_STATE_ENUM.Terminated ? '已结束或已终止状态的项目，不能被终止' : false
  }

  private isDeleteBtnDisabled = (rowDate) => {
    const message = this.isBtnDiabled(rowDate)
    if (message) return message
    return rowDate.state === PROJECT_STATE_ENUM.Running ? '已运行中的项目，不能被删除，可以先进行终止操作' : false
  }

  private onSearch = () => {
    this.projectMG.pageIndex = 1
    this.projectMG.getList(this.props.isAdmin)
  }


  private onPageChange = current => {
    this.projectMG.pageIndex = current
    this.projectMG.getList(this.props.isAdmin)
  }

  private onPageSizeChange = (current, size) => {
    this.projectMG.pageSize = size
    this.projectMG.pageIndex = current
  }

  // private onSelectAll = keys => {
  //   this.updateSelectedRowKeys(keys)
  // }

  // private onSelectInvert = () => {
  //   this.updateSelectedRowKeys([])
  // }

  // private onSelect = (rowKey, checked) => {
  //   let keys = this.selectedRowKeys

  //   if (checked) {
  //     keys = [...keys, rowKey]
  //   } else {
  //     const index = keys.findIndex(item => item === rowKey)
  //     keys.splice(index, 1)
  //   }
  //   this.updateSelectedRowKeys(keys)
  // }
  private viewProject = (rowData) => {
    onPreview(rowData)
  }

  private editProject = (rowData) => {
    onEdit(rowData, () => this.projectMG.getList(this.props.isAdmin))
  }

  private changeOwner = (rowData) => {
    onChangeOwner(rowData, () => this.projectMG.getList(this.props.isAdmin))
  }

  private memberMgr = (rowData) => {
    openUserSelector(rowData, async (users) => {
      await this.projectMG.changeMembers({
        project_id: rowData.id,
        user_ids: users.map(item => item.user_id)
      })
      message.success('操作项目成员成功')
      this.projectMG.getList(this.props.isAdmin)
    })
  }

  private termination = (rowData) => {
    Modal.show({
      title: '终止项目',
      content: `确定要终止项目 ${rowData.project_name} 吗？`,
      onOk: async () => {
        await this.projectMG.termination(rowData.id)
        this.updateSelectedRowKeys([])
        message.success(`终止项目 ${rowData.project_name} 成功`)
        this.projectMG.getList(this.props.isAdmin)
      }
    })
  }

  private addProject = () => {
    onAdd(() => this.projectMG.getList(this.props.isAdmin))
  }

  private delete = async (rowData) => {
    Modal.show({
      title: '删除项目',
      content: `确定要删除项目 ${rowData.project_name} 吗 ？${rowData.members.length > 1 ? '注意: 项目中还有其他成员。' : ''}`,
      onOk: async () => {
        await this.projectMG.delete(rowData.id)
        this.updateSelectedRowKeys([])
        message.success(`删除项目 ${rowData.project_name} 成功`)
        this.projectMG.getList(this.props.isAdmin)
      }
    })
  }

  render() {
    return (
      <ListWrapper ref={this.wrapperRef}>
        <TopWrapper>
          <div className='action'>
            <div className='filter'>
              <div className='item'>
                <span className='label'>项目名称: </span> 
                <Input value={this.projectMG.projectName} onChange={(e) => this.projectMG.projectName = e.target.value}/>
              </div>
              <div className='item'>
                <span className='label'>项目周期: </span>
                <RangePicker
                  value={[this.projectMG.startTime,this.projectMG.endTime]}
                  ranges={GeneralDatePickerRange}
                  showTime={{ format: DatePicker_SHOWTIME_FORMAT }}
                  format={DatePicker_FORMAT}
                  onChange={(dates) => {
                    this.projectMG.startTime = dates?.[0]
                    this.projectMG.endTime = dates?.[1]
                  }}
                  allowClear={true}
                />
              </div>
              <div className='item'>
                <span className="label">项目状态: </span>
                <Select
                  style={{width: 160}}
                  value={this.projectMG.state} 
                  onSelect={value => {
                    this.projectMG.state = value
                  }}
                  >
                  <Option key={'-1'} value="">全部</Option>
                  {
                    Object.keys(PROJECT_STATE_MAP).map(key => 
                      <Option 
                        key={key} 
                        value={key}> 
                        {PROJECT_STATE_MAP[key]} 
                      </Option>
                  )}
                </Select>
              </div>
            </div>
            <div>
              <Button className="btn" onClick={() => this.onSearch()}>查询</Button>
            </div>
          </div>
          <div className='action'>
            <Button disabled={this.isAddBtnDisabled()} className="btn" type="primary" onClick={() => this.addProject()}>新增项目</Button>         
          </div>
        </TopWrapper>
        <Table
          columns={this.columns}
          props={{
            data: this.tableData,
            height: this.height,
            loading: this.loading,
            shouldUpdateScroll: false,
            locale: {
              emptyMessage: '没有项目数据',
              loading: '数据加载中...'
            },
            rowKey: 'id'
          }}
          // rowSelection={{
          //   selectedRowKeys: this.selectedRowKeys,
          //   onSelect: this.onSelect,
          //   onSelectAll: this.onSelectAll,
          //   onSelectInvert: this.onSelectInvert,
          //   props: {
          //     fixed: 'left'
          //   }
          // }}
        />
        <div style={{ textAlign: 'center', padding: 10 }}>
          <Pagination
            showSizeChanger
            onChange={this.onPageChange}
            current={this.projectMG.pageIndex}
            total={this.projectMG.total}
            onShowSizeChange={this.onPageSizeChange}
          />
        </div>
      </ListWrapper>
    )
  }
}
