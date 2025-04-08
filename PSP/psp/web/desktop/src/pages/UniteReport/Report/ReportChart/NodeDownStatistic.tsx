// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import LineChart from '../../../../components/Chart/LineChart'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage } from '../../util'
import {
  CoreDataTable,
  IProps,
  Loading,
  SwitchChartType,
  download
} from './common'
import report from '@/domain/UniteReport'
import qs from 'qs'
import { observable } from 'mobx'
import { message } from 'antd'
import PieChart from '@/components/Chart/PieChart'
import { roundNumber } from '@/utils'
import { PRECISION } from '@/domain/common'
import { SinglePieWrapper } from '../style'

const title = '节点宕机统计'

@observer
export default class NodeDownStatistic extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: 'pie' | 'area' | 'line' = 'line'

  @observable
  loading = true

  @observable
  exporting = false

  @observable
  chartData = {
    node_down_number_rate: [],
    node_down_number: {
      data: [],
      fields: [],
      originData: [],
      name: ''
    }
  }

  @observable isShowSummary = true

  getFields(originData) {
    return originData.map(d => d.key)
  }

  getData(name, originData) {
    return originData.map(d => {
      return {
        name,
        [d.key]: d.value
      }
    })
  }

  fetchChartData = async () => {
    this.loading = true
    try {
      const data = await report.getReportByType(
        this.props.reportType,
        this.props.reportDates
      )
      this.chartData.node_down_number_rate[0] =
        data?.node_down_number_rate || []
      ;['node_down_number'].forEach(k => {
        let { name, original_data } = data[k]
        this.chartData[k] = {
          data: this.getData(name, original_data),
          fields: this.getFields(original_data),
          originData: original_data,
          name
        }
      })
    } catch (e) {
      message.error('获取报表数据失败')
    } finally {
      this.props.stopLoading()
      this.loading = false
    }
  }

  async componentDidMount() {
    this.fetchChartData()
  }

  async componentDidUpdate(prevProps) {
    if (
      prevProps.reportType !== this.props.reportType ||
      prevProps.reportDates[0] !== this.props.reportDates[0] ||
      prevProps.reportDates[1] !== this.props.reportDates[1]
    ) {
      this.fetchChartData()
    }
  }

  handleExportImage = () => {
    exportImage(this.ref.current, `${title}报表`)
  }

  handleExportExcel = () => {
    if (this.exporting) {
      return
    }

    this.exporting = true
    message.info('正在生成导出文件，请耐心等待')

    let url =
      '/api/v1/report/export/nodeDownStatistics?' +
      qs.stringify(
        { start: this.props.reportDates[0], end: this.props.reportDates[1] },
        { arrayFormat: 'repeat' }
      )

    console.debug(url)
    download(url, '')

    setTimeout(() => {
      this.exporting = false
    }, 2000)
  }

  tooltipFomatter = (key, value) => {
    console.log(key, value)
    return {
      name: key,
      value: roundNumber(value, PRECISION)
    }
  }

  render() {
    return (
      <div>
        <ChartAction
          chartData={this.chartData}
          exportExcel={this.handleExportExcel}
          exportImage={this.handleExportImage}
          otherBtn={() => (
            <>
              <SwitchChartType
                value={this.type}
                onClick={type => {
                  this.type = type
                  this.fetchChartData()
                }}
                configs={[
                  { icon: 'line-chart', value: 'line', iconTitle: '折线图' },
                  { icon: 'pie-chart', value: 'pie', iconTitle: '饼图' }
                ]}
              />
            </>
          )}
        />
        <div ref={this.ref}>
          {this.loading ? (
            <Loading tip='报表生成中...' />
          ) : (
            <>
              {this.type === 'line' && (
                <LineChart
                  type={this.type}
                  title={title}
                  data={this.chartData?.node_down_number_rate || []}
                  unit={'%'}
                  min={0}
                  max={120}
                  /* 
               由于后端给出了空数组的值，需要明确开始时间和结束时间，
               不能以数组的第一项为开始，最后一项为结束，进行过滤。
               计算最大值和最小值，以及平均值 
               */
                  start={this.props.reportDates[0]}
                  end={this.props.reportDates[1]}
                  summaryPosition={['50%', '-8%']}
                  summaryWidth={this.props.width + 'px'}
                  showSummary={this.isShowSummary}
                  showDates={this.props.reportDates}
                />
              )}
              {this.type === 'pie' && (
                <>
                  <SinglePieWrapper>
                    <div className='chart'>
                      <PieChart
                        data={this.chartData?.node_down_number?.originData}
                        unit={'次'}
                        title={`${title}`}
                        tooltipFomatter={this.tooltipFomatter}
                        showDates={this.props.reportDates}
                      />
                    </div>
                  </SinglePieWrapper>
                  <SinglePieWrapper>
                    <div className='data'>
                      <CoreDataTable
                        data={this.chartData?.node_down_number?.originData}
                        type={'节点名称'}
                        column_unit={'次数'}
                        unit={'次'}
                      />
                    </div>
                  </SinglePieWrapper>
                </>
              )}
            </>
          )}
        </div>
      </div>
    )
  }
}
