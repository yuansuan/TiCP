// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import LineChart from '../../../../components/Chart/LineChart'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage } from '../../util'
import { IProps, Loading, SwitchChartType, download } from './common'
import { IconControlSummary } from './commonShowIcon'
import report from '@/domain/UniteReport'
import qs from 'qs'
import { observable } from 'mobx'
import { message } from 'antd'

const title = '内存平均利用率'

@observer
export default class MEMUsageReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: 'area' | 'line' = 'line'

  @observable
  loading = true

  @observable
  exporting = false

  @observable
  chartData = []

  @observable isShowSummary = true

  get reportType() {
    return this.props.reportType.toLowerCase()
  }

  fetchChartData = async () => {
    this.loading = true
    try {
      const data = await report.getReportByType(
        this.props.reportType,
        this.props.reportDates
      )
      this.chartData = data
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
      '/api/v1/report/export/resourceUtAvg?' +
      qs.stringify(
        {
          start: this.props.reportDates[0],
          end: this.props.reportDates[1],
          type: this.props.reportType
        },
        { arrayFormat: 'repeat' }
      )

    console.debug(url)
    download(url, '')

    setTimeout(() => {
      this.exporting = false
    }, 2000)
  }

  toggleShowSummary = () => {
    this.isShowSummary = this.isShowSummary ? false : true
    this.fetchChartData()
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
              <IconControlSummary
                isShowSummary={this.isShowSummary}
                onClick={this.toggleShowSummary}
              />
              <SwitchChartType
                value={this.type}
                onClick={type => {
                  this.type = type
                  this.fetchChartData()
                }}
                configs={[
                  { icon: 'line-chart', value: 'line', iconTitle: '折线图' },
                  { icon: 'area-chart', value: 'area', iconTitle: '面积图' }
                ]}
              />
            </>
          )}
        />
        <div ref={this.ref}>
          {this.loading ? (
            <Loading tip='报表生成中...' />
          ) : (
            <LineChart
              type={this.type}
              title={title}
              data={this.chartData[this.reportType] || []}
              unit={'%'}
              min={0}
              max={100}
              summaryPosition={['50%', '-8%']}
              /* 
               由于后端给出了空数组的值，需要明确开始时间和结束时间，
               不能以数组的第一项为开始，最后一项为结束，进行过滤。
               计算最大值和最小值，以及平均值 
               */
              start={this.props.reportDates[0]}
              end={this.props.reportDates[1]}
              summaryWidth={this.props.width + 'px'}
              showSummary={this.isShowSummary}
              showDates={this.props.reportDates}
            />
          )}
        </div>
      </div>
    )
  }
}
