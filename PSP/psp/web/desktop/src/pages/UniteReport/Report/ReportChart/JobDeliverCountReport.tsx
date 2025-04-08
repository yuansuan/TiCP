// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import LineChart from '../../../../components/Chart/LineChart'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage, exportExcel } from '../../util'
import { IProps, Loading, SwitchChartType, download } from './common'
import { IconControlSummary } from './commonShowIcon'
import report from '@/domain/UniteReport'
import moment from 'moment'
import qs from 'qs'
import { observable } from 'mobx'
import { message } from 'antd'

const title = '用户数与作业数情况'

@observer
export default class JobDeliverCountReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: 'area' | 'line' = 'line'

  @observable
  loading = true

  @observable
  exporting = false

  @observable
  chartData = {
    job_deliver_job_count: [],
    job_deliver_user_count: []
  }

  @observable isShowSummary = true

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
      '/api/v1/report/export/jobDeliverCount?' +
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
            <>
              <LineChart
                type={this.type}
                title={'提交作业的用户数(按天)'}
                data={this.chartData['job_deliver_user_count'] || []}
                timeFormat={'MM-DD'}
                unit={''}
                summaryPosition={['50%', '-8%']}
                summaryWidth={this.props.width + 'px'}
                showSummary={this.isShowSummary}
                showDates={this.props.reportDates}
              />
              <LineChart
                type={this.type}
                title={'投递作业数(按天)'}
                data={this.chartData['job_deliver_job_count'] || []}
                timeFormat={'MM-DD'}
                unit={''}
                summaryPosition={['50%', '-8%']}
                summaryWidth={this.props.width + 'px'}
                showSummary={this.isShowSummary}
                showDates={this.props.reportDates}
              />
            </>
          )}
        </div>
      </div>
    )
  }
}
