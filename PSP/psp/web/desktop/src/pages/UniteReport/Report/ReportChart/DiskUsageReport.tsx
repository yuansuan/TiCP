// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage, exportExcel } from '../../util'
import { IProps, Loading, SwitchChartType } from './common'
import report from '@/domain/UniteReport'
import { observable } from 'mobx'
import ColumnChart from '../../../../components/Chart/ColumnChart'
import { message } from 'antd'

const title = '磁盘使用情况'

@observer
export default class DiskUsageReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: 'horizontal' | 'vertical' = 'vertical'

  @observable
  loading = true

  @observable
  chartData = {
    diskUsageByDir: { data: [], fields: [] },
    diskUsageByUser: { data: [], fields: [] },
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
    if (!this.chartData) return

    const map = {
      diskUsageByDir: {
        name: '磁盘使用情况(按目录)',
        header: ['', ...this.chartData.diskUsageByDir.fields],
      },
      diskUsageByUser: {
        name: '磁盘使用情况(按用户)',
        header: ['', ...this.chartData.diskUsageByUser.fields],
      },
    }

    const sheets = []
    Object.keys(map).forEach(key => {
      const sheet = this.chartData[key].data.map(item => {
        let rowData = [`${item.name} (GB)`]
        for (let i = 1; i < map[key].header.length; i++) {
          rowData.push(item[map[key].header[i]])
        }
        return rowData
      })
      sheet.unshift(map[key].header)
      sheets.push({
        sheetName: map[key].name,
        data: sheet,
      })
    })
    exportExcel({
      excelName: title,
      sheets,
    })
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
                { icon: 'bar-chart', value: 'vertical', iconTitle: '柱状图' },
                {
                  icon: 'bar-chart',
                  value: 'horizontal',
                  iconTitle: '条形图',
                  rotate: 90,
                },
              ]}
            />
          )}
        />
        <div ref={this.ref}>
          {this.loading ? (
            <Loading tip='报表生成中...' />
          ) : (
            <>
              <ColumnChart
                type={this.type}
                data={this.chartData.diskUsageByDir?.data}
                fields={this.chartData.diskUsageByDir?.fields}
                unit={'GB'}
                title={`${title}(按共享存储)`}
              />
              <ColumnChart
                type={this.type}
                data={this.chartData.diskUsageByUser?.data}
                fields={this.chartData.diskUsageByUser?.fields}
                unit={'GB'}
                title={`${title}(按用户)`}
              />
            </>
          )}
        </div>
      </div>
    )
  }
}
