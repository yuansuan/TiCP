import * as React from 'react'
import { Spin, Table } from 'antd'
import {
  AreaChartOutlined,
  PieChartOutlined,
  BarChartOutlined,
  LineChartOutlined,
  BarsOutlined
} from '@ant-design/icons'
import { toJS } from 'mobx'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'

const IconMap = {
  'line-chart': LineChartOutlined,
  'area-chart': AreaChartOutlined,
  'bar-chart': BarChartOutlined,
  'pie-chart': PieChartOutlined,
  bars: BarsOutlined
}

const sortfn = (a, b) =>
  b['value'] - a['value'] === 0
    ? a['key'] >= b['key']
      ? 1
      : -1
    : b['value'] - a['value']

export function Loading({ tip }) {
  return (
    <div
      style={{
        width: '100%',
        height: 400,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
      }}>
      <Spin tip={tip} />
    </div>
  )
}
export interface IProps {
  stopLoading: () => void
  reportType: string
  licenseId?: string
  licenseType?: string
  reportDates: number[]
  width?: string
}

interface config {
  icon: string
  value: string
  iconTitle: string
  rotate?: number
}

interface ISwitchChartTypeProps {
  onClick: (value) => void
  value: string
  configs: config[]
}

export function SwitchChartType(props: ISwitchChartTypeProps) {
  const { onClick, value, configs } = props

  const getStyle = (value, type) => ({
    color: value === type ? 'blue' : '#aaa',
    margin: '5px',
    fontSize: 20
  })

  /**
   *  [{icon: 'line-chart', value: 'line'},{icon: 'area-chart', value: 'area'}]
   *  [{icon: 'bar-chart', value: 'vertical'},{icon: 'bar-chart', value: 'horizontal', rotate:90}]
   *  [{icon: 'bar-chart', value: 'vertical'},{icon: 'bar-chart', value: 'horizontal', rotate:90},{icon: 'pie-chart', value: 'pie'}]
   */

  return (
    <div style={{ margin: '0 15px' }}>
      {configs.map(c => {
        const ICON = IconMap[c.icon]
        return (
          <ICON
            key={c.value}
            title={c.iconTitle}
            rotate={c.rotate || 0}
            style={getStyle(value, c.value)}
            onClick={() => onClick(c.value)}
          />
        )
      })}
    </div>
  )
}

export const JobDataTable = function ({ data, type }) {
  data = toJS(data)
  const total = data.reduce((s, v) => (s += v.value), 0)
  data.sort(sortfn)

  return (
    <Table
      footer={() => `作业总数: ${total || 0}`}
      pagination={false}
      size='small'
      columns={[
        {
          title: type,
          dataIndex: 'key',
          key: 'key'
        },
        {
          title: '作业数（个）',
          dataIndex: 'value',
          key: 'value'
        },
        {
          title: '作业数占比',
          dataIndex: '',
          key: '',
          render: (text, record) => {
            return `${roundNumber((record['value'] / total) * 100, PRECISION)}%`
          },
          sorter: (a, b) => a['value'] - b['value'],
          sortDirections: ['descend', 'ascend']
        }
      ]}
      dataSource={data}
    />
  )
}

export const CoreDataTable = function ({
  data,
  type,
  column_unit = '核时',
  unit = '小时'
}) {
  data = toJS(data)
  const total = data.reduce((s, v) => (s += v.value), 0)
  data.sort(sortfn)

  return (
    <Table
      footer={() =>
        `总${column_unit}: ${total ? roundNumber(total, PRECISION) : 0} ${unit}`
      }
      pagination={false}
      size='small'
      columns={[
        {
          title: type,
          dataIndex: 'key',
          key: 'key'
        },
        {
          title: `${column_unit}(${unit})`,
          dataIndex: 'value',
          key: 'value',
          render: value => {
            return `${roundNumber(value, PRECISION)}`
          }
        },
        {
          title: `${column_unit}占比`,
          dataIndex: '',
          key: '',
          render: (text, record) => {
            return total === 0
              ? '--'
              : `${roundNumber((record['value'] / total) * 100, PRECISION)}%`
          },
          sorter: (a, b) => a['value'] - b['value'],
          sortDirections: ['ascend', 'descend']
        }
      ]}
      dataSource={data}
    />
  )
}

export function download(href: string, name: string) {
  const a = document.createElement('a')
  a.style.display = 'none'
  document.body.appendChild(a)
  a.href = href
  a.download = name
  a.click()
  a.remove()
}
