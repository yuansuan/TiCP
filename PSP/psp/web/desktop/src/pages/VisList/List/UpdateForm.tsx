import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { DatePicker, Radio } from 'antd'
import { UpdateWapper } from '../style'
import { Content, Sider, SiderTitle } from '../Create/styles'
import moment from 'moment'
import dayjs from 'dayjs'
import { Modal } from '@/components'

interface IProps {
  onCancel: () => void
  onOk: (autoClose: boolean, time?: string) => void | Promise<void>
  rowData: any
}

export const UpdateForm = observer(function UpdateForm({
  onCancel,
  onOk,
  rowData
}: IProps) {
  const [autoClose, setAutoClose] = useState(rowData.session?.is_auto_close)
  const [closeTime, setCloseTime] = useState(
    rowData.session?.auto_close_time?.seconds
  )
  const onChange = (date, dateString) => {
    setCloseTime(date?.unix())
  }

  const disabledDate = current => {
    return current && current < dayjs().startOf('day')
  }

  const range = (start, end) => {
    const result = []
    for (let i = start; i < end; i++) {
      result.push(i)
    }
    return result
  }
  const disabledDateTime = date => ({
    disabledHours: () => range(0, 24).splice(0, dayjs().hour() + 1)
  })

  const disabledBtn = () => {
    if (autoClose) {
      const nowTime = new Date().getTime() / 1000
      return (
        Number(closeTime) - nowTime < 3600 &&
        '所选关闭时间需大于当前时间一小时以上'
      )
    } else {
      return false
    }
  }

  return (
    <>
      <UpdateWapper>
        <Sider>
          <SiderTitle>关闭时间</SiderTitle>
        </Sider>
        <Content style={{ paddingLeft: '30px' }}>
          <Radio.Group
            onChange={e => setAutoClose(e.target.value)}
            value={autoClose}>
            <Radio value={false}>永不关闭</Radio>
            <Radio value={true}>
              <DatePicker
                inputReadOnly={true}
                disabled={!autoClose}
                value={
                  closeTime === 0
                    ? undefined
                    : closeTime && moment.unix(closeTime)
                }
                showTime={{ format: 'HH' }}
                format='YYYY-MM-DD HH:00'
                disabledDate={disabledDate}
                disabledTime={disabledDateTime}
                onChange={onChange}
                showNow={false}
                allowClear={false}
              />
            </Radio>
          </Radio.Group>
        </Content>
      </UpdateWapper>
      <Modal.Footer
        onCancel={onCancel}
        okButtonProps={{ disabled: disabledBtn() as any }}
        onOk={() => {
          onOk(autoClose, closeTime)
        }}
      />
    </>
  )
})
