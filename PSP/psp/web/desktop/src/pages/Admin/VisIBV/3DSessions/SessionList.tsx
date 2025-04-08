/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Table, Button, Icon, Modal } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from '../store'
import { CopyOutlined } from '@ant-design/icons'
import { copyToClipboard } from '@/utils/Validator'
import { message } from 'antd'
import styled from 'styled-components'
import { Spin } from 'antd'
import { StatusWrapper } from './style'
import { SESSION_STATUS_MAP } from '@/domain/Vis'
import { CloseSessionForm } from './CloseSessionForm'

import moment from 'moment'
export const COLOR_MAP = {
  等待资源: {
    color: '#2A8FDF',
    borderColor: '#C7E3F9'
  },
  启动中: {
    color: '#2A8FDF',
    borderColor: '#C7E3F9'
  },
  已启动: {
    color: '#63B03D',
    borderColor: '#D7F9C7'
  },
  出错: {
    color: '#EF5350',
    borderColor: '#F9D9D9'
  },
  不可用: {
    color: '#C5C5C5',
    borderColor: '#E6E4E4'
  },
  删除中: {
    color: '#FF9100',
    borderColor: '#FDEFC7'
  },
  已删除: {
    color: '#C5C5C5',
    borderColor: '#E6E4E4'
  }
}

const StyledDiv = styled.div`
  > .anticon {
    display: none;
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 18px;
  }

  &:hover {
    > .anticon {
      display: inline-block;
      &:hover {
        color: ${props => props.theme.primaryColor};
      }
    }
  }
`

const tableConfig = {
  id: 'session_table',
  columns: [
    'out_app_id', // 会话ID
    'software.name', //镜像名称
    'project_name',
    'software.desc', // 镜像描述
    'software_platform_str', // 操作平台
    'hardware.name', // 实例名称
    'status', //会话状态（全部、等待资源、启动中、已启动、删除中、已删除）
    'user_name', // 创建者
    'create_time', // 创建时间
    'start_time', // 开始时间
    'end_time', // 结束时间
    'duration', //时长（小时）
    'opts' //操作（删除会话，新做）
  ]
}

// if (!hasPerm) {
//   tableConfig.columns = tableConfig.columns.slice(0, -1)
// }
interface IProps {
  height?: number
}

export const SessionList = observer(function SessionList(props: IProps) {
  const store = useStore()
  const { dataSource } = useLocalStore(() => ({
    get dataSource() {
      return store.model.list
    }
  }))

  async function handleCloseTask(rowData) {
    Modal.show({
      title: '删除会话',
      width: 600,
      footer: null,
      bodyStyle: { padding: 0 },
      content: ({ onCancel, onOk }) => (
        <CloseSessionForm
          rowData={rowData}
          onCancel={onCancel}
          onOk={() => {
            store.fetchSessionList()
            onOk()
          }}
        />
      )
    })
  }

  async function copySessionPassword(rowData: any) {
    navigator.clipboard.writeText(rowData.machine_password).then(
      () => {
        message.success(`成功复制 ${rowData.machine_password} 到粘贴板`)
      },
      () => {
        message.warn(`复制密码 ${rowData.machine_password} 失败`)
      }
    )
  }
  let columns = [
    {
      props: {
        width: 150,
        fixed: 'left'
      },
      header: '会话编号',
      cell: {
        props: {
          dataKey: 'out_app_id'
        },
        render: ({ rowData, dataKey }) => (
          <StyledDiv>
            <div>{rowData[dataKey]}</div>
            <CopyOutlined
              rev={'none'}
              onClick={() => {
                copyToClipboard(rowData[dataKey])
                message.success(`${rowData[dataKey]} 已复制到剪贴板`)
              }}
            />
          </StyledDiv>
        )
      }
    },
    {
      props: {
        width: 160,
        resizable: true
      },
      header: '镜像名称',
      dataKey: 'software.name',
      cell: {
        render: ({ rowData, dataKey }) => <span>{rowData?.software?.name}</span>
      }
    },
    {
      props: {
        width: 160,
        resizable: true
      },
      header: '项目名称',
      dataKey: 'project_name',
      cell: {
        render: ({ rowData, dataKey }) => <span>{rowData?.project_name}</span>
      }
    },
    {
      props: {
        width: 160,
        resizable: true
      },
      header: '镜像描述',
      dataKey: 'software.desc',
      cell: {
        render: ({ rowData, dataKey }) => <span>{rowData?.software?.desc}</span>
      }
    },
    {
      props: {
        width: 120,
        resizable: true
      },
      header: '操作平台',
      cell: {
        props: {
          dataKey: 'software_platform_str'
        },
        render: ({ rowData, dataKey }) => <span>{rowData[dataKey]}</span>
      }
    },
    {
      props: {
        width: 160,
        resizable: true
      },
      header: '实例名称',
      dataKey: 'hardware.name',
      cell: {
        render: ({ rowData, dataKey }) => <span>{rowData?.hardware.name}</span>
      }
    },
    {
      props: {
        width: 120,
        resizable: true
      },
      header: '会话状态',
      cell: {
        props: {
          dataKey: 'status'
        },
        render: ({ rowData, dataKey }) => {
          const getIcon = text => {
            if (text === '已启动' || '启动中') {
              return <Icon type='running' />
            } else if (text === '删除中' || '等待资源') {
              return <Icon type='loading' />
            } else {
              return null
            }
          }
          const text = SESSION_STATUS_MAP[rowData.status]

          return (
            <StatusWrapper>
              <div
                className='icon'
                style={{
                  background: COLOR_MAP[text]?.color,
                  border: `2px solid ${COLOR_MAP[text]?.borderColor}`
                }}
              />
              <div className='text'>{text}</div>
              <div className='icon-right'>{getIcon(text)}</div>
            </StatusWrapper>
          )
        }
      }
    },
    {
      props: {
        width: 160,
        resizable: true
      },
      header: '创建者',
      dataKey: 'user_name'
    },

    {
      props: {
        width: 200
      },
      header: '创建时间',
      dataKey: 'create_time'
    },
    {
      props: {
        width: 200
      },
      header: '开始时间',
      dataKey: 'start_time'
    },
    {
      props: {
        width: 200
      },
      header: '结束时间',
      dataKey: 'end_time'
    },
    {
      props: {
        width: 120
      },
      header: '时长（小时）',
      dataKey: 'duration'
    },
    // {
    //   props: {
    //     width: 140
    //   },
    //   header: '是否删除',
    //   cell: {
    //     props: {
    //       dataKey: 'deleted'
    //     },
    //     render: ({ rowData, dataKey }) => (
    //       <span>{rowData[dataKey] ? '是' : '否'}</span>
    //     )
    //   }
    // },
    {
      header: '操作',
      props: {
        fixed: 'right',
        width: 150
      },
      cell: {
        props: {
          dataKey: 'opts'
        },
        render: ({ rowData, dataKey }) => {
          let closeButtonDisabled = !['启动中', '已启动'].includes(
            SESSION_STATUS_MAP[rowData.status]
          )
          return (
            <>
              <Button
                type='link'
                disabled={closeButtonDisabled}
                onClick={() => {
                  handleCloseTask(rowData)
                }}>
                删除会话
              </Button>
            </>
          )
        }
      }
    }
  ]

  // if (!hasPerm) {
  //   columns = columns.slice(0, -1)
  // }

  if (store.fetching)
    return (
      <Spin
        style={{
          height: '250px',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center'
        }}
      />
    )
  else
    return (
      <Table
        tableId={tableConfig.id}
        defaultConfig={tableConfig.columns}
        props={{
          data: dataSource || [],
          rowKey: 'id',
          height: props.height - 120 || 400,
          locale: {
            emptyMessage: '没有会话数据',
            loading: '数据加载中...'
          }
        }}
        columns={columns as any}
      />
    )
})
