import React from 'react'
import styled from 'styled-components'
import { Button, Modal, Table } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from './store'
import { Http } from '@/utils'
import { Divider, message } from 'antd'
import HardwareEditor from './HardwareEditor'

const StyledLayout = styled.div`
  .name {
    cursor: pointer;

    &:hover {
      color: ${({ theme }) => theme.primaryHighlightColor};
    }
  }

  .rs-table-cell-content {
    text-overflow: unset;
  }
`
interface IProps {
  height: number
}

export const HardwareList = observer(function HardwareList(props: IProps) {
  const store = useStore()

  const state = useLocalStore(() => ({
    get dataSource() {
      return store.hardware.hardwareList.map(item => ({
        ...item,
        isPublish: item.enabled ? '已发布' : '未发布'
      }))
    }
  }))

  async function isPublish(rowData) {
    await Http.put('/vis/hardware/status', {
      id: rowData.id,
      enabled: !rowData['enabled']
    })
    store.refreshHardware()
    message.success(rowData['enabled'] ? '取消发布成功' : '发布成功')
  }

  function edit(rowData) {
    Modal.show({
      title: '编辑实例',
      width: 600,
      bodyStyle: { padding: 0, height: 640 },
      footer: null,
      content: ({ onCancel, onOk }) => (
        <HardwareEditor
          hardwareItem={rowData}
          onCancel={onCancel}
          onOk={() => {
            onOk()
            store.refreshHardware()
          }}
        />
      )
    })
  }
  const deleteHardware = ({ id }) => {
    Modal.confirm({
      title: '删除实例',
      content: '确认删除！',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        await Http.delete(`/vis/hardware?id=${id}`, {})
        store.refreshHardware()
      }
    })
  }
  return (
    <StyledLayout>
      <Table
        columns={[
          {
            props: {
              resizable: true,
              width: 220
            },
            header: '实例名称',
            dataKey: 'name'
          },
          {
            props: {
              resizable: true,
              width: 200
            },
            header: '实例描述',
            dataKey: 'desc'
          },
          {
            props: {
              resizable: true,
              width: 200
            },
            header: '实例类型',
            dataKey: 'instance_type'
          },
          {
            props: {
              resizable: true,
              width: 200
            },
            header: '实例机型系列',
            dataKey: 'instance_family'
          },

          {
            props: {
              resizable: true,
              width: 120
            },
            header: '实例最大带宽',
            dataKey: 'network'
          },
          {
            props: {
              resizable: true,
              width: 100
            },
            header: 'CPU核数',
            dataKey: 'cpu'
          },
          {
            props: {
              width: 100
            },
            header: 'GPU数量',
            dataKey: 'gpu'
          },
          {
            props: {
              resizable: true,
              width: 80
            },
            header: '内存',
            dataKey: 'mem'
          },
          // {
          //   props: {
          //     resizable: true,
          //     width: 100
          //   },
          //   header: '状态',
          //   cell: {
          //     props: {
          //       dataKey: 'isPublish'
          //     }
          //   }
          // },
          {
            props: {
              flexGrow: 1,
              fixed: 'right',
              minWidth: 200
            },
            header: '操作',
            cell: {
              props: {
                dataKey: 'id'
              },
              render: ({ rowData }) => (
                <div className='rowData'>
                  {/* <Button
                    type='link'
                    onClick={() => {
                      isPublish(rowData)
                    }}>
                    {rowData['enabled'] ? '取消发布' : '发布'}
                  </Button> */}
                  <Button type='link' onClick={() => edit(rowData)}>
                    编辑
                  </Button>
                  <Divider type='vertical' />
                  <Button
                    type='link'
                    style={{ padding: 0 }}
                    onClick={() => deleteHardware(rowData)}>
                    删除
                  </Button>
                </div>
              )
            }
          }
        ]}
        props={{
          height: props.height - 120 || 400,
          data: state.dataSource,
          rowKey: 'id',
          loading: store.loading,
          virtualized: true,
          isTree: true
        }}
      />
    </StyledLayout>
  )
})
