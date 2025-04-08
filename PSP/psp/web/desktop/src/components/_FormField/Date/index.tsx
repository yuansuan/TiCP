/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { DatePicker } from 'antd'
import { observer } from 'mobx-react'
import moment from 'moment'
import * as React from 'react'

import Container from '../Container'
import { runInTyping } from '../utils'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

const dateFormat = 'YYYY-MM-DD HH:mm'

@observer
export default class DateItem extends React.Component<IProps> {
  public render() {
    const { formModel, model } = this.props
    const { id, defaultValue } = model
    if (!formModel[id]) return null
    const { value } = formModel[id]
    return (
      <Container {...this.props}>
        <DatePicker
          style={{ width: '300px' }}
          defaultValue={defaultValue ? moment(defaultValue) : null}
          value={value ? moment(value) : null}
          format={dateFormat}
          showTime={{ format: 'HH:mm' }}
          onChange={date => {
            runInTyping(
              formModel,
              () => (formModel[id].value = date && date.format(dateFormat))
            )
          }}
        />
      </Container>
    )
  }
}
