// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { REPORT } from '../const'

interface IProps {
  stopLoading: () => void
  reportType: string
  reportDates: number[]
  licenseType: number[]
  licenseId: number[]
}

@observer
export default class ReportChart extends React.Component<IProps> {
  resizeObserver = null
  ref = null
  @observable width = 800

  constructor(props) {
    super(props)
    this.ref = React.createRef()
  }

  componentDidMount() {
    this.resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
        console.log(entry.contentRect.width)
        this.width = entry.contentRect.width
      }
    })
    this.resizeObserver.observe(this.ref.current)
  }

  componentWillUnmount() {
    this.resizeObserver.disconnect()
  }

  render() {
    const { reportType, stopLoading, reportDates,licenseType,licenseId } = this.props
    const Report = REPORT[reportType].ReportChart
    return (
      <div ref={this.ref}>
        <Report
          width={this.width}
          reportType={reportType}
          reportDates={reportDates}
          licenseId={licenseId}
          licenseType={licenseType}
          stopLoading={stopLoading}
        />
      </div>
    )
  }
}
