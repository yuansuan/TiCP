import * as React from 'react'
import styled from 'styled-components'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import { workStationList } from '@/domain/Visual'
import { Table } from '@/components'
import { Tag } from 'antd'

const Wrapper = styled.div`
  padding: 20px;
  .table {
  }
`

@observer
export default class WorkStation extends React.Component<any> {
  componentDidMount = () => {
    workStationList.fetch()
  }
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
        },
        header: '工作站',
        dataKey: 'name',
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
        },
        header: '总数',
        dataKey: 'up_limit',
      },
      {
        props: {
          resizable: true,
        },
        header: '空闲数',
        dataKey: 'free',
      },
      {
        props: {
          flexGrow: 1,
          // width: 300,
        },
        header: '安装软件',
        cell: {
          props: {
            dataKey: 'software_list',
          },
          render: ({ rowData, dataKey }) => {
            const softwares = rowData[dataKey].map(s => {
              return <Tag key={s.name}>{s.name}</Tag>
            })
            return <div>{softwares}</div>
          },
        },
      },
    ]
  }
  render() {
    return (
      <Wrapper>
        <Table
          props={{
            autoHeight: true,
            data: [...workStationList],
            rowKey: 'id',
          }}
          columns={this.columns}
        />
      </Wrapper>
    )
  }
}
