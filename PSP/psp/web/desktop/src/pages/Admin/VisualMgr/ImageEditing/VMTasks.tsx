/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import * as React from 'react'
import styled from 'styled-components'
import { Table, Modal } from '@/components'
import VirtualMachine from '@/domain/Visual/VirtualMachine'
import { Http } from '@/utils'
import { observer } from 'mobx-react'
import { computed, observable } from 'mobx'
import { message } from 'antd'

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
const statusMap = new Map([
  [1, '排队'],
  [2, '提交'],
  [3, '运行'],
  [4, '失败'],
  [5, '关闭']
])
const colorMap = new Map([
  [1, '#F5A623'],
  [2, '#52C41A '],
  [3, '#4A90E2'],
  [4, '#D0021B '],
  [5, '#9B9B9B']
])
@observer
export class VMTasks extends React.Component<{ vm: VirtualMachine }> {
  @observable tasks = []

  @computed
  get columns() {
    return [
      {
        props: {
          width: 80
        },
        header: 'ID',
        dataKey: 'id'
      },
      {
        props: {
          width: 160
        },
        header: '应用',
        dataKey: 'template_name'
      },
      {
        props: {
          width: 100
        },
        header: '创建人',
        dataKey: 'user_name'
      },
      {
        props: {
          resizable: true
        },
        header: '状态',
        cell: {
          props: {
            dataKey: 'statusName'
          },
          render: ({ rowData, dataKey }) => {
            return (
              <div>
                <StatusBall color={rowData['statusColor']} />
                <span>{rowData[dataKey]}</span>
              </div>
            )
          }
        }
      },
      {
        props: {
          resizable: true,
          width: 200
        },
        header: '开始时间',
        cell: {
          props: {
            dataKey: 'start_time'
          },
          render: ({ rowData, dataKey }) => {
            const startTime = new Date(rowData[dataKey])
            return <div>{startTime.toDateString()}</div>
          }
        }
      },
      {
        props: {
          flexGrow: 1
        },
        header: '操作',
        cell: {
          render: ({ rowData }) => {
            return rowData.status != 3 ? null : (
              <div>
                <a onClick={() => this.endTask(rowData)}>关闭</a>
              </div>
            )
          }
        }
      }
    ]
  }
  endTask = async task => {
    await Modal.showConfirm({
      title: '关闭云工作站',
      content: '确定关闭云工作站？'
    })
    await Http.post(
      '/visual/worktask/stop',
      {
        user_id: task.user_id,
        work_task_id: task.id
      },
      { baseURL: '' }
    )
      .then(() => {
        message.info('操作成功！')
        this.fetchTasks()
      })
      .catch(() => {
        message.error('操作失败！')
      })
  }
  componentDidMount = async () => {
    this.fetchTasks()
  }
  fetchTasks = async () => {
    const { vm } = this.props
    const res = await Http.get(`/visual/vm/${vm.id}/tasks`, { baseURL: '' })
    let tasks = res.data.worktask_list
    const userIdentities = tasks.map(t => {
      return { id: t.user_id }
    })
    const res2 = await Http.post('/user/batch', { userIdentities })
    let userMap = {}
    res2.data.list.forEach(user => {
      userMap[user.id] = user
    })
    tasks.forEach(t => {
      t.user_name = userMap[t.user_id].name
      t.statusName = statusMap.get(t.status)
      t.statusColor = colorMap.get(t.status)
    })

    this.tasks = tasks
  }
  render() {
    return (
      <div>
        <Table
          props={{
            autoHeight: true,
            data: this.tasks,
            rowKey: 'id'
          }}
          columns={this.columns}
        />
      </div>
    )
  }
}
