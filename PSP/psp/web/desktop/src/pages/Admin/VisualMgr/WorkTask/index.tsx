import * as React from 'react'
import styled from 'styled-components'
import { computed } from 'mobx'
import { observer } from 'mobx-react'

import { workTaskList } from '@/domain/Visual'
import { Table, Modal } from '@/components'

const Wrapper = styled.div`
  padding: 20px;
  .table {
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

@observer
export default class WorkTask extends React.Component<any> {
  constructor(props) {
    super(props)
    workTaskList.fetch()
  }
  componentDidMount = () => {}
  @computed
  get columns() {
    return [
      {
        props: {
          resizable: true,
        },
        header: 'ID',
        dataKey: 'id',
      },
      {
        props: {
          resizable: true,
          width: 200,
        },
        header: '应用名',
        dataKey: 'template_name',
      },
      {
        props: {
          resizable: true,
          width: 100,
        },
        header: '创建人',
        dataKey: 'user_name',
      },
      {
        props: {
          resizable: true,
          width: 150,
        },
        header: '工作站',
        dataKey: 'workstation_name',
      },
      {
        props: {
          resizable: true,
        },
        header: '操作系统',
        dataKey: 'os',
      },
      {
        props: {
          resizable: true,
          width: 200,
        },
        header: '开始时间',
        cell: {
          props: {
            dataKey: 'start_time',
          },
          render: ({ rowData, dataKey }) => {
            const startTime = rowData[dataKey]
            return <div>{startTime.toDateString()}</div>
          },
        },
      },
      {
        props: {
          resizable: true,
        },
        header: '状态',
        cell: {
          props: {
            dataKey: 'statusName',
          },
          render: ({ rowData, dataKey }) => {
            return (
              <div>
                <StatusBall color={rowData['statusColor']} />
                <span>{rowData[dataKey]}</span>
              </div>
            )
          },
        },
      },
      {
        props: {
          flexGrow: 1,
        },
        header: '操作',
        cell: {
          render: ({ rowData }) => {
            return (
              <div>
                <a onClick={() => this.endTask(rowData)}>关闭</a>
              </div>
            )
          },
        },
      },
    ]
  }
  endTask = worktask => {
    Modal.showConfirm({
      title: '关闭云工作站',
      content: '确定关闭云工作站？',
    })
      .then(() => {
        workTaskList.remove(worktask)
      })
      .catch(() => {})
  }
  render() {
    return (
      <Wrapper>
        <Table
          props={{
            autoHeight: true,
            data: [...workTaskList],
            rowKey: 'id',
          }}
          columns={this.columns}
        />
      </Wrapper>
    )
  }
}
