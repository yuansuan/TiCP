// Copyright (C) 2019 LambdaCal Inc.

import * as React from 'react'
import { observer } from 'mobx-react'

import { Select, Button, DatePicker } from 'antd'
import { ActionWrapper } from './style'
import {
  REPORT_TYPE_LABEL,
  REPORT_TYPE_OPTIONS,
  REPORT,
  CLOUD_REPORT_TYPE_OPTIONS,
  VISUAL_REPORT_TYPE_OPTIONS
} from '../const'
import { observable, action } from 'mobx'
import moment, { Moment } from 'moment'
import {
  DatePicker_FORMAT,
  DatePicker_SHOWTIME_FORMAT,
  DeployMode
} from '@/constant'
import { sysConfig } from '@/domain'

const { RangePicker } = DatePicker

const Option = Select.Option

interface IValues {
  reportType: string
  dates: number[]
}

interface IProps {
  loading: boolean
  showSelect: boolean
  onSubmit: (values: IValues) => void
}

type Duration = '24h' | '7d' | '30d' | '180d' | '360d' | 'custom'

const durations = [
  {
    name: '24h',
    label: '24小时',
    value: 1
  },
  { name: '7d', label: '7天', value: 7 },
  { name: '30d', label: '30天', value: 30 },
  { name: '180d', label: '180天', value: 180 },
  { name: '360d', label: '360天', value: 360 }
]

@observer
export default class Action extends React.Component<IProps> {
  @observable reportType = REPORT_TYPE_OPTIONS[0]
  @observable dates: Moment[] = [moment().subtract(7, 'days'), moment()] // 默认七天
  @observable disabledDates = false
  @observable duration: Duration = '7d'

  @action
  submit = () => {
    this.props.onSubmit({
      reportType: this.reportType,
      dates: this.dates.map(d => d.valueOf())
    })
  }

  onChange = (type, value) => {
    this[type] = value

    if (type === 'reportType' && REPORT[value] && REPORT[value].disableDates) {
      this.disabledDates = true
    } else {
      this.disabledDates = false
    }

    if (type !== 'dates') {
      this.submit()
    }
  }

  onDurationChange = (duration, num) => {
    this['dates'] = [moment().subtract(num, 'days'), moment()]
    this.duration = duration

    this.submit()
  }

  render() {
    return (
      <ActionWrapper>
        {this.props.showSelect && (
          <div className='item'>
            <label>报表类型: </label>
            <Select
              showSearch
              style={{ width: 200 }}
              placeholder='选择报表类型'
              onChange={value => {
                this.onChange('reportType', value)
              }}
              value={this.reportType}
              optionFilterProp='children'>
              {sysConfig?.globalConfig?.enable_visual
                ? REPORT_TYPE_OPTIONS.concat(VISUAL_REPORT_TYPE_OPTIONS).map(
                    option => (
                      <Option key={option} value={option}>
                        {REPORT_TYPE_LABEL[option]}
                      </Option>
                    )
                  )
                : REPORT_TYPE_OPTIONS.map(option => (
                    <Option key={option} value={option}>
                      {REPORT_TYPE_LABEL[option]}
                    </Option>
                  ))}
            </Select>
          </div>
        )}

        {!this.disabledDates && (
          <div className='item'>
            <label>时间: </label>
            <div className='time'>
              {durations.map(d => (
                <Button
                  key={d.name}
                  type='link'
                  onClick={() => {
                    this.onDurationChange(d.name, d.value)
                  }}
                  style={{
                    padding: '0 5px',
                    color:
                      d.name === this.duration && !this.disabledDates
                        ? 'blue'
                        : 'inherit'
                  }}>
                  {d.label}
                </Button>
              ))}
            </div>
            <RangePicker
              // ranges={{
              //   最近24小时: [moment().subtract(1, 'days'), moment()],
              //   最近7天: [moment().subtract(7, 'days'), moment()],
              //   最近30天: [moment().subtract(30, 'days'), moment()],
              //   最近180天: [moment().subtract(180, 'days'), moment()],
              // }}
              // @ts-ignore
              value={this.dates}
              showTime={{ format: DatePicker_SHOWTIME_FORMAT }}
              format={DatePicker_FORMAT}
              placeholder={['开始时间', '结束时间']}
              onChange={value => {
                this.onChange('dates', value)
              }}
              allowClear={false}
              style={{
                borderRadius: 5,
                border:
                  'custom' === this.duration && !this.disabledDates
                    ? '1px solid blue'
                    : 'inherit'
              }}
              onOk={() => {
                this.duration = 'custom'
                this.submit()
              }}
            />
          </div>
        )}
      </ActionWrapper>
    )
  }
}
