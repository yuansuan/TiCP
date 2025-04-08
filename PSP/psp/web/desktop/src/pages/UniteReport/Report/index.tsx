// Copyright (C) 2019 LambdaCal Inc.
import * as React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'

import Action from './Action'
import ReportChart from './ReportCharts'
import { REPORT_TYPE_OPTIONS } from '../const'
import moment from 'moment'

@observer
export default class Report extends React.Component<any> {
  @observable dates = [moment().subtract(7, 'days'), moment()].map(m =>
    m.valueOf()
  )
  @observable reportType = REPORT_TYPE_OPTIONS[0]

  @observable loading = false

  onSubmit = values => {
    this.dates = values.dates
    this.reportType = values.reportType
  }

  render() {
    return (
      <div>
        <Action
          onSubmit={this.onSubmit}
          showSelect={this.props.showSelect === false ? false : true}
          loading={this.loading}
        />
        <ReportChart
          reportType={this.props.chartType || this.reportType}
          reportDates={this.dates}
          licenseId={this.props.licenseId}
          licenseType={this.props.licenseType}
          stopLoading={() => {
            this.loading = false
          }}
        />
      </div>
    )
  }
}
