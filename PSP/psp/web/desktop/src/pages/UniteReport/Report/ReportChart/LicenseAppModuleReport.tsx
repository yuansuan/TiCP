// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'
import LineChart from '../../../../components/Chart/LineChart'
import ChartAction from '../../../../components/Chart/ChartAction'
import { exportImage, exportExcel } from '../../util'
import { IProps, Loading, SwitchChartType } from './common'
import { IconControlSummary } from './commonShowIcon'
import report from '@/domain/UniteReport'
import moment from 'moment'
import { observable } from 'mobx'
import { message } from 'antd'

const title = '许可证软件模块使用情况'

@observer
export default class LicenseAppModuleReport extends React.Component<IProps> {
  ref = React.createRef<HTMLDivElement>()

  @observable
  type: 'area' | 'line' = 'line'

  @observable
  loading = true

  @observable
  chartData = []

  @observable isShowSummary = true

  fetchChartData = async () => {
    this.loading = true
    try {
      const data = await report.getReportByType(this.props.reportType, this.props.reportDates, this.props.licenseId,
        this.props.licenseType)
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

    const sheetNameMap = {
      license_app_module_ut_avg: {
        name: '许可证软件模块使用情况',
        unit: '',
      },
    }

    const sheets = []

    Object.keys(this.chartData).forEach(key => {
      const data = this.chartData[key];
      if(data.length===0) return message.warn('许可证软件模块为空！')
      const maxLengthData = data.reduce((max, current) => (current.d.length > max.d.length ? current : max), data[0]);
    
      const sheetData = [];
      let keys = data.map(d => d.n);
      maxLengthData.d.forEach((item, index) => {
        const d = [moment(item.t).format('MM-DD HH:mm:ss')];
        keys.forEach(key => {
          const item = data.find(tmp => tmp.n === key);
          d.push(item.d[index] ? item.d[index].v : ''); 
        });
        sheetData.push(d);
      });
    
      const unitStr = sheetNameMap[key].unit ? `(${sheetNameMap[key].unit})` : '';
      keys = keys.map(k => `${k} ${unitStr}`);
      sheetData.unshift(['时间', ...keys]);
      sheets.push({
        sheetName: sheetNameMap[key].name,
        data: sheetData,
      });
    });
    exportExcel({
      excelName: `${title}报表`,
      sheets,
    })
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
                  { icon: 'area-chart', value: 'area', iconTitle: '面积图' },
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
                title={'许可证软件模块使用情况'}
                data={this.chartData['license_app_module_ut_avg'] || []}
                unit={''}
                timeFormat={'MM-DD HH:mm:ss'}
                summaryPosition={[ '50%', '-8%' ]}
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
