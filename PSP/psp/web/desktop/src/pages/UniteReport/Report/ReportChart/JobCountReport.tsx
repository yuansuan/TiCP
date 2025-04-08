// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage, exportExcel } from '../../util'
import {
  IProps,
  Loading,
  SwitchChartType,
  JobDataTable,
  download
} from './common'
import report from '@/domain/UniteReport'
import { observable } from 'mobx'
import { message } from 'antd'
import qs from 'qs'
import { PieWrapper } from '../style'
import ColumnChart from '../../../../components/Chart/ColumnChart'
import PieChart from '../../../../components/Chart/PieChart'
import SortedChart from '../../../../components/Chart/SortedChart'

const title = '作业投递数情况'
type chartType = 'horizontal' | 'vertical'

@observer
export default class JobCountReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: chartType | 'pie' | 'top10' = 'horizontal'

  @observable
  loading = true

  @observable
  exporting = false

  @observable
  chartData = {
    job_count_by_app: { data: [], fields: [], originData: [], name: null },
    job_count_by_user: { data: [], fields: [], originData: [], name: null }
  }

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

      ;['job_count_by_app', 'job_count_by_user'].forEach(k => {
        let { name, original_data } = data[k]

        this.chartData[k] = {
          data: this.getData(name, original_data),
          fields: this.getFields(original_data),
          originData: original_data,
          name
        }
      })
      // console.log(toJS(this.chartData))
    } catch (e) {
      message.error('获取报表数据失败')
      console.error(e)
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
      '/api/v1/report/export/jobCount?' +
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

  render() {
    return (
      <div>
        <ChartAction
          chartData={this.chartData}
          exportExcel={this.handleExportExcel}
          exportImage={this.handleExportImage}
          otherBtn={() => (
            <SwitchChartType
              value={this.type}
              onClick={type => (this.type = type)}
              configs={[
                {
                  icon: 'bar-chart',
                  value: 'horizontal',
                  rotate: 90,
                  iconTitle: '条形图'
                },
                { icon: 'pie-chart', value: 'pie', iconTitle: '饼图' }
              ]}
            />
          )}
        />
        <div ref={this.ref}>
          {this.loading ? (
            <Loading tip='报表生成中...' />
          ) : (
            <>
              {(this.type === 'vertical' || this.type === 'horizontal') && (
                <>
                  <ColumnChart
                    type={this.type as chartType}
                    data={this.chartData.job_count_by_app.data}
                    fields={this.chartData.job_count_by_app.fields}
                    unit={''}
                    title={`${title}(按软件)`}
                    showDates={this.props.reportDates}
                  />
                  <ColumnChart
                    type={this.type as chartType}
                    data={this.chartData.job_count_by_user.data}
                    fields={this.chartData.job_count_by_user.fields}
                    unit={''}
                    title={`${title}(按用户)`}
                    minTickInterval={1}
                    showDates={this.props.reportDates}
                  />
                </>
              )}
              {this.type === 'pie' && (
                <>
                  <PieWrapper>
                    <div className='chart'>
                      <PieChart
                        data={this.chartData.job_count_by_app.originData}
                        unit={''}
                        title={`${title}(按软件)`}
                        showDates={this.props.reportDates}
                      />
                    </div>
                    <div className='chart'>
                      <PieChart
                        title={`${title}(按用户)`}
                        data={this.chartData.job_count_by_user.originData}
                        unit={''}
                        showDates={this.props.reportDates}
                      />
                    </div>
                  </PieWrapper>
                  <PieWrapper>
                    <div className='data'>
                      <JobDataTable
                        data={this.chartData.job_count_by_app.originData}
                        type='软件'
                      />
                    </div>
                    <div className='data'>
                      <JobDataTable
                        data={this.chartData.job_count_by_user.originData}
                        type='用户'
                      />
                    </div>
                  </PieWrapper>
                </>
              )}
              {this.type === 'top10' && (
                <>
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    data={this.chartData.job_count_by_app.originData}
                    unit={''}
                    top={10}
                    title={`${title}(按软件Top10)`}
                    showDates={this.props.reportDates}
                  />
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    title={`${title}(按用户Top10)`}
                    data={this.chartData.job_count_by_user.originData}
                    unit={''}
                    top={10}
                    showDates={this.props.reportDates}
                  />
                </>
              )}
            </>
          )}
        </div>
      </div>
    )
  }
}
