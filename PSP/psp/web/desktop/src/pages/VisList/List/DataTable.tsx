/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect } from 'react'
import { Modal, Table } from '@/components'
import { useDispatch } from 'react-redux'
import { observer, useLocalStore } from 'mobx-react-lite'
import store from '@/reducers'
import styled from 'styled-components'
import {
  ListSessionResponse,
  ListSessionRequest,
  SESSION_STATUS_MAP,
  SESSION_STATUS_BUTTON_LOADING,
  statusMapping
} from '@/domain/Vis'
import { useStore } from './store'
import { ListDataWrapper, StatusWrapper } from '../style'
import { ColumnProps } from '@/components/Table'
import { Button, Pagination, Space, message, Tooltip, Divider } from 'antd'
import { CopyOutlined } from '@ant-design/icons'
import { copyToClipboard } from '@/utils/Validator'
import { Icon } from '@/components'
import { useInterval } from '@/utils/hooks'
import { UpdateForm } from './UpdateForm'
import { currentUser } from '@/domain'

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
interface IProps {
  response: ListSessionResponse
  request: ListSessionRequest
  loading: boolean
  deleteSession?: (id: string) => Promise<any>
  closeSession?: (id: string) => Promise<any>
  restartSession?: (id: string) => Promise<any>
  powerOnSession?: (id: string) => Promise<any>
  updateSession?: (
    id: string,
    autoClose: boolean,
    time?: string
  ) => Promise<any>
  openSession?: (id: string, row: any) => Promise<string>
  onPagination: (index: number, size: number) => void
  openRemoteApp?: (id: string, row: any) => Promise<string>
  height?: number
}

const colorMap = {
  success: {
    color: '#63B03D',
    borderColor: '#D7F9C7'
  },
  warn: {
    color: '#FF9100',
    borderColor: '#FDEFC7'
  },
  error: {
    color: '#EF5350',
    borderColor: '#F9D9D9'
  },
  primary: {
    color: '#2A8FDF',
    borderColor: '#C7E3F9'
  },
  canceled: {
    color: '#C5C5C5',
    borderColor: '#E6E4E4'
  }
}

const tableConfig = {
  id: 'user_session_table',
  columns: [
    'session.out_app_id', // 会话ID
    'session.software.name', //镜像名称
    'session.project_name', // 项目名称
    'session.software.platform', // 操作平台
    'session.hardware.name', // 实例名称
    'session.hardware.instance_type', // 实例类型
    'session.status', //会话状态（全部、等待资源、启动中、已启动、关闭中、已关闭）
    'session.user_name',
    'session.create_time', // 创建时间
    'session.start_time', // 开始时间
    'session.end_time', // 结束时间
    'opts' //操作（关闭会话，新做）
  ]
}

export const DataTable = observer(
  ({
    request,
    response,
    loading,
    onPagination,
    closeSession,
    restartSession,
    powerOnSession,
    deleteSession,
    updateSession,
    openRemoteApp,
    height
  }: IProps) => {
    const { vis } = useStore()
    const dispatch = useDispatch()
    const state = useLocalStore(() => ({
      stopFetch: false,
      setStopFetch(bool) {
        this.stopFetch = bool
      },
      dataSource: new Array().fill({}) as any[],
      setDataSource(data) {
        this.dataSource = data
      },
      scrollPosition: { x: 0, y: 0 },
      setScrollPosition(pos) {
        this.scrollPosition = pos
      }
    }))

    const findVisExistLoading = response?.sessions.find(
      item => item.loading === true
    )

    useEffect(() => {
      if (loading || findVisExistLoading) {
        state.setStopFetch(true)
      } else {
        state.setStopFetch(false)
      }
    }, [loading, findVisExistLoading])

    const formatDataSource = response => {
      if (response?.length) {
        response.forEach(async item => {
          if (
            item.status === 'STARTED' ||
            item.status === 'POWERING ON' ||
            item.status === 'REBOOTING'
          ) {
            item.realStatus = 'READYING'
            if (item.realStatus === 'READIED') {
              state.setStopFetch(false)
              return
            }
            try {
              await vis
                .pollFetchRequest(item.session.id)
                .then(res => {
                  if (res.success == true) {
                    item.loading = !res.data.ready
                    if (res.data.ready) {
                      item['realStatus'] = 'READIED'
                      item.loading = false
                      item.realStatus = 'READIED'
                      state.setStopFetch(false)
                    }
                  } else {
                    state.setStopFetch(false)
                  }
                })
                .then(() => {
                  state.setDataSource(response.session)
                })
                .catch(() => {
                  state.setStopFetch(false)
                })
                .finally(() => {
                  if (findVisExistLoading || item.loading) {
                    state.setStopFetch(false)
                  }
                })
            } catch (err) {
              state.setStopFetch(false)
            }
          }
        })
      }
    }

    useInterval(
      () => {
        formatDataSource(response && response?.sessions)
      },
      state.stopFetch ? 5000 : null
    )

    const onOpenSession = async (sessionId: string, row: any) => {
      try {
        // const generaApps = await fetchSoftware()
        // const findCurrAPP = generaApps.find(app => app.id === sessionId)

        // const data = {
        //   type: findCurrAPP?.action,
        //   payload: 'togg'
        // }

        // store.dispatch(data)
        // 目前改为跳出系统打开
        const stream_url = window.atob(row.session.stream_url)
        if (stream_url) window.open(stream_url)
      } catch (err) {
        message.error('会话打开失败，请联系系统管理员')
      }
    }

    const onOpenRemoteApp = (id: string, row: any) => {
      openRemoteApp(id, row).then((url: string) => {
        if (url) {
          window.open(url)
        } else {
          message.error('远程应用不存在, 检查远程应用是否已经配置')
        }
      })
    }

    const onDeleteSession = (id: string) => {
      Modal.showConfirm({
        title: '删除会话',
        content: '是否确认删除该会话？',
        onOk() {
          return new Promise((resolve, reject) => {
            deleteSession(id)
              .then(
                res => {
                  message.success('删除会话成功')
                  resolve(res)
                },
                () => {
                  reject('删除失败')
                }
              )
              .catch(() => message.error('删除失败'))
          })
        }
      })
    }

    const onCloseSession = (id: string) => {
      Modal.showConfirm({
        title: '删除会话',
        content: '是否确认删除该会话？',
        onOk() {
          return new Promise((resolve, reject) => {
            closeSession(id)
              .then(
                res => {
                  message.success('删除会话成功，会话即将删除...')
                  dispatch({
                    type: 'desktop',
                    data: { sessionId: id },
                    payload: 'closeSession'
                  })
                  resolve(res)
                },
                () => {
                  reject('删除会话失败')
                }
              )
              .catch(() => message.error('删除会话失败'))
          })
        }
      })
    }

    const onRestartSession = (id: string) => {
      Modal.showConfirm({
        title: '重启会话',
        content: '是否确认重启该会话？',
        onOk() {
          return new Promise((resolve, reject) => {
            restartSession(id)
              .then(
                res => {
                  message.success('重启会话成功')
                  dispatch({
                    type: 'desktop',
                    data: { sessionId: id },
                    payload: 'restartSession'
                  })
                  resolve(res)
                },
                () => {
                  reject('重启会话失败')
                }
              )
              .catch(() => message.error('重启会话失败'))
          })
        }
      })
    }

    const onPowerOnSession = (id: string) => {
      Modal.showConfirm({
        title: '开启会话',
        content: '是否确认开启该会话？',
        onOk: async () => {
          return new Promise((resolve, reject) => {
            powerOnSession(id)
              .then(
                res => {
                  message.success('开启会话成功')
                  dispatch({
                    type: 'desktop',
                    data: { sessionId: id },
                    payload: 'restartSession'
                  })
                  resolve(res)
                },
                () => {
                  reject('开启会话失败')
                }
              )
              .catch(() => message.error('开启会话失败'))
          })
        }
      })
    }

    const onUpdateSession = rowData => {
      Modal.show({
        title: '修改关闭时间',
        content: ({ onCancel, onOk }) => {
          const OK = (autoClose, time) => {
            onOk()
            new Promise((resolve, reject) => {
              updateSession(rowData.session?.id, autoClose, time)
                .then(
                  res => {
                    message.success('修改会话关闭时间成功')
                    resolve(res)
                  },
                  () => {
                    reject('修改会话关闭时间失败')
                  }
                )
                .catch(() => message.error('修改会话关闭时间失败'))
            })
          }
          return <UpdateForm onCancel={onCancel} onOk={OK} rowData={rowData} />
        },
        footer: null
      })
    }

    const isSysAdmin = currentUser.hasSysMgrPerm ? true : false
    const isYourSession = rowData =>
      currentUser.name === rowData?.session?.user_name

    const columns: ColumnProps[] = [
      {
        header: '会话编号',
        props: {
          width: 150,
          fixed: 'left'
        },
        cell: {
          props: {
            dataKey: 'session.out_app_id'
          },
          render: ({ rowData }) => (
            <StyledDiv>
              <div>{rowData.session.out_app_id}</div>
              <CopyOutlined
                rev={'none'}
                onClick={() => {
                  copyToClipboard(rowData.session.out_app_id)
                  message.success(
                    `${rowData.session.out_app_id} 已复制到剪贴板`
                  )
                }}
              />
            </StyledDiv>
          )
        }
      },
      {
        header: '镜像名称',
        dataKey: 'session.software.name',
        props: {
          width: 200
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <Tooltip
                placement='topLeft'
                title={rowData.session.software.name}>
                {rowData.session.software.name}
              </Tooltip>
            )
          }
        }
      },
      {
        header: '项目名称',
        dataKey: 'session.project_name',
        props: {
          minWidth: 120,
          flexGrow: 1
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <Tooltip placement='topLeft' title={rowData.session.project_name}>
                {rowData.session.project_name}
              </Tooltip>
            )
          }
        }
      },
      {
        header: '操作平台',
        dataKey: 'session.software.platform',
        props: {
          width: 100,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <Tooltip
                placement='topLeft'
                title={rowData.session.software.platform}>
                {rowData.session.software.platform}
              </Tooltip>
            )
          }
        }
      },
      {
        header: '实例名称',
        dataKey: 'session.hardware.name',
        props: {
          minWidth: 200,
          align: 'center',
          flexGrow: 2
        },
        cell: {
          render: ({ rowData }) => {
            const name = rowData.session.hardware?.name
            return (
              <Tooltip placement='topLeft' title={name}>
                {name}
              </Tooltip>
            )
          }
        }
      },
      {
        header: '实例类型',
        dataKey: 'session.hardware.instance_type',
        props: {
          minWidth: 180,
          align: 'center',
          flexGrow: 2
        },
        cell: {
          render: ({ rowData }) => {
            const type = rowData.session.hardware?.instance_type
            return (
              <Tooltip placement='topLeft' title={type}>
                {type}
              </Tooltip>
            )
          }
        }
      },
      {
        header: '会话状态',
        dataKey: 'session.status',
        props: {
          width: 100,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
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
            const type = statusMapping[SESSION_STATUS_MAP[rowData.status]]

            return (
              <StatusWrapper>
                <div
                  className='icon'
                  style={{
                    background: colorMap[type]?.color,
                    border: `2px solid ${colorMap[type]?.borderColor}`
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
        header: '创建者',
        dataKey: 'session.user_name',
        props: {
          width: 120,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
            return `${rowData?.session?.user_name || '--'}`
          }
        }
      },
      {
        header: '创建时间',
        dataKey: 'session.create_time',
        props: {
          width: 200,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
            return `${rowData?.session?.create_time}` || '--'
          }
        }
      },
      {
        header: '开始时间',
        dataKey: 'session.start_time',
        props: {
          width: 200,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
            return `${rowData?.session?.start_time}` || '--'
          }
        }
      },
      {
        header: '结束时间',
        dataKey: 'session.end_time',
        props: {
          width: 200,
          resizable: true
        },
        cell: {
          render: ({ rowData }) => {
            return `${rowData?.session?.end_time}` || '--'
          }
        }
      },
      {
        header: '操作',
        dataKey: 'opts',
        props: {
          minWidth: 260,
          flexGrow: 2,
          fixed: 'right'
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <Space size={0}>
                {rowData.session.status === 'CLOSED' ? (
                  <>--</>
                ) : (
                  // <Button
                  //   type='link'
                  //   onClick={() => onDeleteSession(rowData.session.id)}>
                  //   删除
                  // </Button>
                  <>
                    <Tooltip
                      title={SESSION_STATUS_BUTTON_LOADING[rowData.realStatus]}
                      color={'#108ee9'}
                      key={'#108ee9'}>
                      <Button
                        type='link'
                        disabled={
                          isYourSession(rowData) || isSysAdmin
                            ? rowData.session.status !== 'STARTED' ||
                              (rowData.session.status === 'STARTED' &&
                                rowData.loading)
                            : true
                        }
                        loading={
                          rowData.session.status === 'STARTING' ||
                          (rowData.session.status === 'STARTED' &&
                            rowData.loading)
                        }
                        onClick={() =>
                          onOpenSession(rowData.session.id, rowData)
                        }>
                        打开
                      </Button>
                    </Tooltip>
                    <Divider type='vertical' style={{ margin: 2 }} />
                    {/*开机操作 */}
                    {['POWERING OFF', 'POWER OFF', 'POWERING ON'].includes(
                      rowData.session.status
                    ) ? (
                      <Button
                        type='link'
                        disabled={
                          rowData.session.status !== 'POWERING OFF' &&
                          rowData.session.status !== 'POWER OFF'
                        }
                        loading={
                          rowData.session.status === 'POWERING OFF' ||
                          rowData.session.status === 'POWERING ON'
                        }
                        onClick={() => onPowerOnSession(rowData.session.id)}>
                        {rowData.session.status === 'POWERING OFF'
                          ? '关机中'
                          : rowData.session.status === 'POWERING ON'
                          ? '开机中'
                          : '开机'}
                      </Button>
                    ) : (
                      <Button
                        type='link'
                        disabled={
                          rowData.session.status !== 'STARTED' &&
                          rowData.session.status !== 'UNAVAILABLE'
                        }
                        loading={rowData.session.status === 'REBOOTING'}
                        onClick={() => onRestartSession(rowData.session.id)}>
                        {rowData.session.status === 'REBOOTING'
                          ? '重启中'
                          : '重启'}
                      </Button>
                    )}

                    <Divider type='vertical' style={{ margin: 2 }} />
                    <Button
                      type='link'
                      danger
                      disabled={
                        isYourSession(rowData) || isSysAdmin
                          ? rowData.session.status === 'CLOSING' ||
                            rowData.session.status === 'CLOSED'
                          : true
                      }
                      onClick={() => onCloseSession(rowData.session.id)}>
                      删除
                    </Button>
                  </>
                )}
              </Space>
            )
          }
        }
      }
    ]

    return (
      <ListDataWrapper>
        <div className='main'>
          <Table
            columns={columns}
            tableId={tableConfig.id}
            defaultConfig={tableConfig.columns}
            props={{
              loading,
              rowKey: 'id',
              height: height,
              shouldUpdateScroll: false,
              data: response?.sessions || state.dataSource,
              locale: {
                emptyMessage: '还没有创建会话',
                loading: '数据加载中...'
              }
            }}></Table>
          <Pagination
            className='pagination'
            showSizeChanger
            pageSize={request.page_size || 10}
            onChange={onPagination}
            current={request.page_index}
            total={response?.total || 10}
          />
        </div>
      </ListDataWrapper>
    )
  }
)
