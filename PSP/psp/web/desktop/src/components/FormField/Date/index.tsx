import { DatePicker } from 'antd'
import { observer } from 'mobx-react'
import moment from 'moment'
import * as React from 'react'

import Container from '../Container'
import Editor from './Editor'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

const dateFormat = 'YYYY-MM-DD HH:mm'

@observer
export default class DateItem extends React.Component<IProps> {
  public static Editor = Editor

  constructor(props) {
    super(props)

    const { formModel, model } = props
    if (formModel) {
      formModel[model.id] = {
        ...model,
        value: model.value || model.defaultValue,
        values: model.values.length > 0 ? model.values : model.defaultValues,
      }
    }
  }

  public render() {
    const { formModel, model } = this.props
    const { id, defaultValue } = model
    if (!formModel[id]) return null
    const { value } = formModel[id]

    return (
      <Container {...this.props}>
        <DatePicker
          showToday={false}
          style={{ width: '300px' }}
          defaultValue={defaultValue ? moment(defaultValue) : null}
          value={value ? moment(value) : null}
          format={dateFormat}
          showTime={{ format: 'HH:mm' }}
          onChange={date => {
            formModel[id].value = date && date.format(dateFormat)
          }}
        />
      </Container>
    )
  }
}
