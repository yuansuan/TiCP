import React from 'react'
import { observer } from 'mobx-react'
import { Chart, Geom, Axis, Tooltip, Coord, Legend, Guide } from 'bizcharts'
import DataSet from '@antv/data-set'
import styled from 'styled-components'
import { Table } from '@/components'
import { computed } from 'mobx'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'
import { EllipsisWrapper } from '@/components'
import { legendFormatterByLength } from '@/utils/chartLegendFormatter'

const Wrapper = styled.div`
  .jobTitle {
    display: flex;
    align-items: center;

    .name {
      padding-left: 20px;
      font-size: 12px;
      color: #797979;
    }
  }

  .chartTable {
    display: flex;
    .chart {
      min-width: 240px;
      width: 40%;
    }
    .table {
      width: 60%;

      .rs-table-row {
        border-bottom: none;
      }
      .rs-table-row-header {
        background: #fff;
        .rs-table-cell-content {
          font-size: 14px;
          font-weight: 650;
          color: #333333;
          line-height: 34px;
        }
      }
      .rs-table-cell-content {
        background: rgba(255, 255, 255, 1);
        line-height: 30px;
      }
    }
  }
`

interface IProps {
  data: any
  total: number
}

@observer
export default class AppJobNumInfo extends React.Component<IProps> {
  @computed
  get tableData() {
    const { data } = this.props
    const top5Data = data.slice(0, 5)
    return top5Data
  }

  @computed
  get columns(): any {
    const cols = [
      {
        props: {
          width: 80,
          align: 'center',
        },
        header: 'TOP 5',
        dataKey: 'top5',
      },
      {
        props: {
          flexGrow: 2,
        },
        header: '应用名称',
        dataKey: 'key',
        cell: {
          render: ({ rowData }) => (
            <>
              <EllipsisWrapper>{rowData.key}</EllipsisWrapper>
            </>
          ),
        },
      },
      {
        props: {
          flexGrow: 1,
          align: 'center',
        },
        header: '作业数',
        dataKey: 'value',
      },
    ]
    return cols
  }

  render() {
    const { data } = this.props
    const { DataView } = DataSet
    const { Html } = Guide
    const dv = new DataView()
    dv.source(data).transform({
      type: 'percent',
      field: 'value',
      dimension: 'key',
      as: 'percent',
    })
    const cols = {
      percent: {
        formatter: val => {
          val = `${roundNumber(val * 100, PRECISION)}%`
          return val
        },
      },
    }

    return (
      <Wrapper>
        <div className='jobTitle'>
          <div className='title'>应用作业数</div>
          <span className='name'>过去24小时应用作业数(TOP5)</span>
        </div>
        <div className='chartTable'>
          <div className='chart'>
            <Chart
              height={200}
              width={250}
              padding={'auto'}
              data={dv}
              scale={cols}>
              <Coord type={'theta'} radius={0.8} innerRadius={0.77} />
              <Axis name='percent' />
              <Legend
                position='right'
                offsetY={-15}
                offsetX={-15}
                itemFormatter={val => legendFormatterByLength(val, 20)}
              />
              <Tooltip
                showTitle={false}
                itemTpl='<li><span style="background-color:{color};" class="g2-tooltip-marker"></span>{name}: {value}</li>'
              />
              <Guide>
                <Html
                  position={['50%', '50%']}
                  html={`<div style="color:#262626;font-size:18px;text-align: center;width: 10em;">应用总数<br><span style="color:#262626;font-size:18px">${
                    this.props.total || 0
                  }</span></div>`}
                  alignX='middle'
                  alignY='middle'
                />
              </Guide>
              <Geom
                type='intervalStack'
                position='percent'
                color='key'
                tooltip={[
                  'key*percent',
                  (key, percent) => {
                    percent = `${roundNumber(percent * 100, PRECISION)}%`
                    return {
                      name: key,
                      value: percent,
                    }
                  },
                ]}
                style={{
                  lineWidth: 1,
                  stroke: '#fff',
                }}></Geom>
            </Chart>
          </div>
          <div className='table'>
            <Table
              columns={this.columns}
              props={{
                data: this.tableData,
                headerHeight: 34,
                rowHeight: 30,
                rowKey: 'top5',
                width: 400,
                style: {
                  width: '100%',
                  height: '100%'
                },
                locale: {
                  emptyMessage: '没有数据',
                } as any,
              }}
            />
          </div>
        </div>
      </Wrapper>
    )
  }
}
