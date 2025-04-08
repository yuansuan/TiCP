/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'
import { computed } from 'mobx'
import { observer, Observer } from 'mobx-react'
import { vmList } from '@/domain/Visual'
import { Table, Button, Modal } from '@/components'
import VirtualMachine from '@/domain/Visual/VirtualMachine'
import { message } from 'antd'
import { VMTasks } from './VMTasks'
import { VMDetail } from './VMDetail'
import { Http } from '@/utils'

const Wrapper = styled.div`
  padding: 20px;

  .rs-table-cell-expand-wrapper {
    float: left;
  }
`
const StatusBall = styled.div`
  display: inline-block;
  &:before {
    display: inline-block;
    content: '';
    height: 6px;
    width: 6px;
    border-radius: 50%;
    background: ${props => props.color};
    margin-right: 10px;
    margin-bottom: 2px;
  }
`
const PrepareWrapper = styled.div`
  padding-left: 16px;
`

const IPWrapper = styled.div`
  display: inline-block;
  margin-left: 16px;
  .ip {
    height: 26px;
    line-height: 26px;
  }
`
@observer
export default class VMImages extends React.Component<any> {
  _timer = null
  componentDidMount = () => {
    vmList.fetch()
    this._timer = setInterval(() => {
      vmList.fetch()
    }, 8 * 1000)
  }
  componentWillUnmount = () => {
    clearInterval(this._timer)
  }
  editImage = async (vm: VirtualMachine) => {
    vm.edit().catch(() => {
      message.error('编辑出错!')
    })
  }
  openEditPage = work_task => {
    const { link, user_id, id } = work_task
  }
  closeEditing = async (vm: VirtualMachine) => {
    await Modal.showConfirm({
      content: '确定要关闭镜像编辑任务？'
    })
    vm.close()
      .then(() => {
        message.info('关闭成功！')
      })
      .catch(() => {
        message.error('关闭失败!')
      })
  }
  offlineVM = (vm: VirtualMachine) => {
    vm.offline()
      .then(() => {
        message.info('操作成功!')
      })
      .catch(() => {
        message.error('操作失败!')
      })
  }
  onlineVM = (vm: VirtualMachine) => {
    vm.online()
      .then(() => {
        message.info('操作成功!')
      })
      .catch(() => {
        message.error('操作失败!')
      })
  }
  showTasks = (vm: VirtualMachine) => {
    Modal.show({
      title: '当前运行的任务',
      bodyStyle: { height: 310, background: '#F0F5FD', overflow: 'auto' },
      width: 830,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <div>
          <VMTasks vm={vm}></VMTasks>
        </div>
      )
    }).catch(() => {})
  }
  showDetail = (vm: VirtualMachine) => {
    Modal.show({
      title: '虚拟机详情',
      bodyStyle: { height: 310, background: '#F0F5FD', overflow: 'auto' },
      width: 630,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <div>
          <VMDetail vm={vm}></VMDetail>
        </div>
      )
    }).catch(() => {})
  }

  editResourceNum = async (num, vm) => {
    if (num === vm.resource_number) return

    await Http.post(
      `/visual/vm/updateResourceNumber/${vm.id}/${num}`,
      {},
      {
        baseURL: ''
      }
    )

    await vmList.fetch()
    message.success('操作成功！')
  }
  @computed
  get columns() {
    return [
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '主机名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => {
            const vm = vmList.get(rowData['id'], !!rowData.children)
            const hostName = rowData.name
            return () =>
              rowData.children ? (
                <div>{hostName}</div>
              ) : (
                <Button
                  type='link'
                  onClick={() => {
                    this.showDetail(vm)
                  }}>
                  {vm.name}
                </Button>
              )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '镜像路径',
        dataKey: 'path'
      },
      {
        props: {
          resizable: true
        },
        header: '状态',
        cell: {
          props: {
            dataKey: 'status'
          },
          render: ({ rowData, dataKey }) => {
            const vm = vmList.get(rowData['id'], !!rowData.children)
            return vm.children ? null : (
              <Observer>
                {() => (
                  <div>
                    <StatusBall color={vm.statusColor} />
                    <span>{vm.status}</span>
                  </div>
                )}
              </Observer>
            )
          }
        }
      },

      {
        props: {
          resizable: true
        },
        header: '当前任务数',
        cell: {
          props: {
            dataKey: 'running_task_num'
          },
          render: ({ rowData }) => {
            const vm = vmList.get(rowData['id'], !!rowData.children)
            return rowData.children ? null : vm.running_task_num &&
              !vm.editing ? (
              <Button
                type='link'
                onClick={() => {
                  this.showTasks(vm)
                }}>
                {vm.running_task_num}
              </Button>
            ) : (
              <span style={{ paddingLeft: 16 }}>{vm.running_task_num}</span>
            )
          }
        }
      },
      {
        props: {
          resizable: true
        },
        header: '总任务数',
        cell: {
          props: {
            dataKey: 'resource_number'
          }
          // render: ({ rowData, dataKey }) => {
          //   const vm = vmList.get(rowData['id'], !!rowData.children)
          //   return vm.children ? null : (
          //     <Observer>
          //       {() => (
          //         <EditableCell
          //           value={vm.resource_number}
          //           showEdit={vm.os_name === 'linux'}
          //           onChange={value => this.editResourceNum(value, vm)}
          //         />
          //       )}
          //     </Observer>
          //   )
          // },
        }
      },
      {
        props: {
          resizable: true
        },
        header: 'CPU核数',
        dataKey: 'allocate_cpu'
      },
      {
        props: {
          resizable: true
        },
        header: '内存',
        cell: {
          props: {
            dataKey: 'allocate_mem_giga'
          },
          render: ({ rowData }) => {
            const vm = vmList.get(rowData['id'], !!rowData.children)
            return rowData.children ? null : <span>{vm.allocate_mem_giga}</span>
          }
        }
      },
      {
        props: {
          flexGrow: 1
        },
        header: '操作',
        cell: {
          props: {},
          render: ({ rowData }) => {
            const vm = vmList.get(rowData['id'], !!rowData.children)
            return vm.children ? null : (
              <Observer>
                {() => (
                  <div>
                    {!vm.editing &&
                      vm.running_task_num == 0 &&
                      vm.status === 'online' && (
                        <div>
                          <Button
                            type='link'
                            onClick={() => this.editImage(vm)}>
                            编辑
                          </Button>
                          <Button
                            type='link'
                            onClick={() => this.offlineVM(vm)}>
                            下线
                          </Button>
                        </div>
                      )}
                    {vm.status === 'offline' && (
                      <Button type='link' onClick={() => this.onlineVM(vm)}>
                        上线
                      </Button>
                    )}
                    {vm.editing && vm.work_task && (
                      <div>
                        {vm.work_task.status === 3 && (
                          <Button
                            type='link'
                            onClick={() => this.openEditPage(vm.work_task)}>
                            打开
                          </Button>
                        )}
                        <Button
                          type='link'
                          onClick={() => this.closeEditing(vm)}>
                          关闭
                        </Button>
                      </div>
                    )}
                    {(vm.status === 'changing' ||
                      (vm.editing &&
                        vm.work_task &&
                        vm.work_task.status < 3)) && (
                      <PrepareWrapper>准备中...</PrepareWrapper>
                    )}
                  </div>
                )}
              </Observer>
            )
          }
        }
      }
    ]
  }
  render() {
    const d = [...vmList]
    return !d.length ? null : (
      <Wrapper>
        <Table
          props={{
            autoHeight: true,
            data: d,
            rowKey: 'id',
            isTree: true,
            defaultExpandAllRows: true
          }}
          columns={this.columns}
        />
      </Wrapper>
    )
  }
}
