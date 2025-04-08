// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage } from '../../util'
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

const title = '可视化使用时长'
type chartType = 'horizontal' | 'vertical'

@observer
export default class VisualUsageDurationReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: chartType | 'pie' | 'top10' = 'horizontal'

  @observable
  loading = true

  @observable
  chartData = {
    usage_duration_by_software: {
      data: [],
      fields: [],
      originData: [],
      name: ''
    },
    usage_duration_by_user: { data: [], fields: [], originData: [], name: '' }
  }

  @observable
  exporting = false

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

      ;['usage_duration_by_software', 'usage_duration_by_user'].forEach(k => {
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
      '/api/v1/vis/statistic/report/duration/export?' +
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
                    data={this.chartData.usage_duration_by_software.data}
                    fields={this.chartData.usage_duration_by_software.fields}
                    unit={'小时'}
                    title={`${title}(按软件)`}
                    tooltipFomatter={this.tooltipFomatter}
                    showDates={this.props.reportDates}
                  />
                  <ColumnChart
                    type={this.type as chartType}
                    data={this.chartData.usage_duration_by_user.data}
                    fields={this.chartData.usage_duration_by_user.fields}
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
                        data={
                          this.chartData.usage_duration_by_software.originData
                        }
                        unit={'小时'}
                        title={`${title}(按软件)`}
                        tooltipFomatter={this.tooltipFomatter}
                        showDates={this.props.reportDates}
                      />
                    </div>
                    <div className='chart'>
                      <PieChart
                        title={`${title}(按用户)`}
                        data={this.chartData.usage_duration_by_user.originData}
                        unit={'小时'}
                        tooltipFomatter={this.tooltipFomatter}
                        showDates={this.props.reportDates}
                      />
                    </div>
                  </PieWrapper>
                  <PieWrapper>
                    <div className='data'>
                      <CoreDataTable
                        data={
                          this.chartData.usage_duration_by_software.originData
                        }
                        type={'软件'}
                        column_unit={'时长'}
                        unit={'小时'}
                      />
                    </div>
                    <div className='data'>
                      <CoreDataTable
                        data={this.chartData.usage_duration_by_user.originData}
                        type={'用户'}
                        column_unit={'时长'}
                        unit={'小时'}
                      />
                    </div>
                  </PieWrapper>
                </>
              )}
              {this.type === 'top10' && (
                <>
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    data={this.chartData.usage_duration_by_software.originData}
                    unit={'小时'}
                    top={10}
                    title={`${title}(按软件Top10)`}
                    tooltipFomatter={this.tooltipFomatter}
                    showDates={this.props.reportDates}
                  />
                  <SortedChart
                    padding={['auto', '10%', 'auto', 'auto'] as any}
                    title={`${title}(按用户Top10)`}
                    data={this.chartData.usage_duration_by_user.originData}
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
