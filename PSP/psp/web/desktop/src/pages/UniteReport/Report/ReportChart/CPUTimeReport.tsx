// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage, exportExcel } from '../../util'
import { roundNumber } from '@/utils/formatter'
import {
  IProps,
  Loading,
  SwitchChartType,
  CoreDataTable,
  download
} from './common'
import report from '@/domain/UniteReport'
import { observable, toJS } from 'mobx'
import { message } from 'antd'
import qs from 'qs'
import ColumnChart from '../../../../components/Chart/ColumnChart'
import PieChart from '../../../../components/Chart/PieChart'
import SortedChart from '../../../../components/Chart/SortedChart'
import { PieWrapper } from '../style'
import { PRECISION } from '@/domain/common'

const title = '核时使用情况'
type chartType = 'horizontal' | 'vertical'

@observer
export default class CPUTimeReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: chartType | 'pie' | 'top10' = 'horizontal'

  @observable
  loading = true

  @observable
  exporting = false

  @observable
  chartData = {
    cpu_time_by_app: { data: [], fields: [], originData: [], name: '' },
    cpu_time_by_user: { data: [], fields: [], originData: [], name: '' }
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

      ;['cpu_time_by_app', 'cpu_time_by_user'].forEach(k => {
        let { name, original_data } = data[k]

        this.chartData[k] = {
          data: this.getData(name, original_data),
          fields: this.getFields(original_data),
          originData: original_data,
          name
        }
      })
      console.log(toJS(this.chartData))
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
      '/api/v1/report/export/cpuTimeSum?' +
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
                    data={this.chartData.cpu_time_by_app.data}
                    fields={this.chartData.cpu_time_by_app.fields}
                    unit={'小时'}
                    title={`${title}(按软件)`}
                    tooltipFomatter={this.tooltipFomatter}
                    showDates={this.props.reportDates}
                  />
                  <ColumnChart
                    type={this.type as chartType}
                    data={this.chartData.cpu_time_by_user.data}
                    fields={this.chartData.cpu_time_by_user.fields}
                    unit={'小时'}
                    title={`${title}(按用户)`}
                    tooltipFomatter={this.tooltipFomatter}
                    showDates={this.props.reportDates}
                  />
                </>
              )}
              {this.type === 'pie' && (
                <>
                  <PieWrapper>
                    <div className='chart'>
                      <PieChart
                        data={this.chartData.cpu_time_by_app.originData}
                        unit={'小时'}
                        title={`${title}(按软件)`}
                        tooltipFomatter={this.tooltipFomatter}
                        showDates={this.props.reportDates}
                      />
                    </div>
                    <div className='chart'>
                      <PieChart
                        title={`${title}(按用户)`}
                        data={this.chartData.cpu_time_by_user.originData}
                        unit={'小时'}
                        tooltipFomatter={this.tooltipFomatter}
                        showDates={this.props.reportDates}
                      />
                    </div>
                  </PieWrapper>
                  <PieWrapper>
                    <div className='data'>
                      <CoreDataTable
                        data={this.chartData.cpu_time_by_app.originData}
                        type={'软件'}
                      />
                    </div>
                    <div className='data'>
                      <CoreDataTable
                        data={this.chartData.cpu_time_by_user.originData}
                        type={'用户'}
                      />
                    </div>
                  </PieWrapper>
                </>
              )}
              {this.type === 'top10' && (
                <>
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    data={this.chartData.cpu_time_by_app.originData}
                    unit={'小时'}
                    top={10}
                    title={`${title}(按软件Top10)`}
                    tooltipFomatter={this.tooltipFomatter}
                    showDates={this.props.reportDates}
                  />
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    title={`${title}(按用户Top10)`}
                    data={this.chartData.cpu_time_by_user.originData}
                    top={10}
                    unit={'小时'}
                    tooltipFomatter={this.tooltipFomatter}
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
