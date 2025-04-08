/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'
import { computed, observable } from 'mobx'
import { observer, Observer } from 'mobx-react'
import { nodeList, VM_OS_TYPE } from '@/domain/Visual'
import { Table, Button, Modal } from '@/components'
import { message, Descriptions } from 'antd'
import CreateVMForm from './CreateVM'
import MachineNode from '@/domain/Visual/MachineNode'
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

const IPWrapper = styled.div`
  display: flex;
  flex-flow: wrap;
  margin-left: 16px;

  .ip {
    height: 24px;
    line-height: 24px;
    width: 108px;
    margin-top: 2px;
  }
`
@observer
export default class VMMgr extends React.Component<any> {
  @observable loading = false
  @observable expendedRowKeys = []

  fresh = async isFirst => {
    if (isFirst) {
      // 第一次进来 clear old data
      nodeList.list = []
      this.loading = true
    }
    try {
      await nodeList.fetch()

      // 第一次进来 clear old data
      if (isFirst) {
        this.expendedRowKeys = nodeList.list.map(item => item.id)
      }
    } finally {
      this.loading = false
    }
  }

  _timer = null
  componentDidMount = async () => {
    await this.fresh(true)
    this._timer = setInterval(async () => {
      await this.fresh(false)
    }, 8 * 1000)
  }
  componentWillUnmount = () => {
    clearInterval(this._timer)
  }

  previewVMInfo = rowData => {
    Modal.show({
      title: '预览虚拟机信息',
      content: () => {
        const { gpu_domain, gpu_bus, gpu_slot, gpu_function } = rowData
        const domain_name = `${gpu_domain}:${gpu_bus}:${gpu_slot}:${gpu_function}`
        return (
          <>
            <Descriptions title='基本信息' column={1}>
              <Descriptions.Item label='虚拟机名称'>
                {rowData['name']}
              </Descriptions.Item>
              <Descriptions.Item label='操作系统'>
                {VM_OS_TYPE[rowData['image_os_type']] || '--'}
              </Descriptions.Item>
              <Descriptions.Item label='GPU'>
                {rowData['gpu_name'] && rowData['gpu_domain']
                  ? `${rowData['gpu_name']} -- ${domain_name}`
                  : '--'}
              </Descriptions.Item>
              <Descriptions.Item label='CPU核数'>
                {rowData['allocate_cpu'] || '--'}
              </Descriptions.Item>
              <Descriptions.Item label='内存(MB)'>
                {rowData['allocate_mem'] / 1024 || '--'}
              </Descriptions.Item>
              {/* <Descriptions.Item label="网络类型">{rowData['network'] || '--'}</Descriptions.Item>
            {
              rowData['network'] === 'static' && (<>
                 <Descriptions.Item label="IP地址">{rowData['ip_address'] || '--'}</Descriptions.Item>
                 <Descriptions.Item label="子网掩码">{rowData['ip_mask'] || '--'}</Descriptions.Item>
              </>)
            } */}
            </Descriptions>
            <Descriptions title='其它信息' column={1}>
              <Descriptions.Item label='虚拟机路径'>
                {rowData['image_path'] || '--'}
              </Descriptions.Item>
              <Descriptions.Item label='虚拟机状态'>
                {rowData['status'] || '--'}
              </Descriptions.Item>
              <Descriptions.Item label='Agent状态'>
                {rowData['resource_status'] || '--'}
              </Descriptions.Item>
            </Descriptions>
          </>
        )
      }
    })
  }

  @computed
  get columns() {
    return [
      {
        props: {
          resizable: true,
          width: 300
        },
        header: '主机名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => {
            const ips = rowData.agent_ip.split(',')
            return () =>
              rowData.children ? (
                <h4>{rowData.name}</h4>
              ) : (
                <Button type='link' onClick={() => this.previewVMInfo(rowData)}>
                  {rowData[dataKey]}
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
        dataKey: 'image_path'
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: 'IP地址',
        cell: {
          props: {
            dataKey: 'vm_ip'
          },
          render: ({ rowData, dataKey }) => {
            return rowData.children ? null : (
              <Observer>
                {() => (
                  <div>
                    <span>{rowData.vm_ip}</span>
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
        header: '虚拟机状态',
        cell: {
          props: {
            dataKey: 'status'
          },
          render: ({ rowData, dataKey }) => {
            // const vm = vmList.get(rowData['id'], !!rowData.children)
            return rowData.children ? null : (
              <Observer>
                {() => (
                  <div>
                    <StatusBall color={rowData.statusColor} />
                    <span>{rowData.status}</span>
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
        header: 'Agent状态',
        cell: {
          props: {
            dataKey: 'resource_status'
          },
          render: ({ rowData, dataKey }) => {
            // const vm = vmList.get(rowData['id'], !!rowData.children)
            return rowData.children ? null : (
              <Observer>
                {() => (
                  <div>
                    <StatusBall color={rowData.resourceStatusColor} />
                    <span>{rowData.resource_status}</span>
                  </div>
                )}
              </Observer>
            )
          }
        }
      },
      {
        props: {
          flexGrow: 1
        },
        header: '',
        cell: {
          props: {},
          render: ({ rowData }) => {
            return !rowData.children ? (
              <Button
                type='link'
                disabled={
                  'creating' === rowData.status || rowData.running_task_num > 0
                }
                onClick={() => {
                  this.onDeleteVM(rowData.id)
                }}>
                删除
              </Button>
            ) : (
              <Button
                type='link'
                onClick={() => {
                  this.onCreateVM(rowData.id)
                }}>
                新建虚拟机
              </Button>
            )
          }
        }
      }
    ]
  }
  async onDeleteVM(id: string) {
    await Modal.showConfirm({
      content: '确定要删除该虚拟机？'
    })

    Http.delete(`/visual/vm/${id}`, {
      data: {
        remove_resource: true,
        remove_vm_machine: true,
        remove_from_pbs: true
      },
      baseURL: ''
    })
      .then(() => {
        message.info('删除成功！')
        nodeList.fetch()
      })
      .catch(() => {
        message.error('删除失败！')
      })
  }
  onCreateVM(id: string) {
    const node = nodeList.list.find((n: MachineNode) => {
      return n.id === id
    })
    Modal.show({
      title: '新建虚拟机',
      bodyStyle: { height: 600, background: '#F0F5FD', overflow: 'auto' },
      width: 630,
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = () => {
          nodeList.fetch()
          onOk()
        }

        return (
          <div>
            <CreateVMForm node={node} onOk={ok}></CreateVMForm>
          </div>
        )
      }
    }).catch(() => {})
  }
  render() {
    const d = nodeList.list
    return (
      <Wrapper>
        <Table
          props={{
            autoHeight: true,
            data: d,
            loading: this.loading,
            rowKey: 'id',
            isTree: true,
            expandedRowKeys: this.expendedRowKeys,
            onExpandChange: (expanded, rowData) => {
              if (expanded) {
                this.expendedRowKeys.push(rowData['id'])
              } else {
                let index = this.expendedRowKeys.indexOf(rowData['id'])

                if (index !== -1) {
                  this.expendedRowKeys.splice(index, 1)
                }
              }
            },
            defaultExpandAllRows: true
          }}
          columns={this.columns as any}
        />
      </Wrapper>
    )
  }
}
