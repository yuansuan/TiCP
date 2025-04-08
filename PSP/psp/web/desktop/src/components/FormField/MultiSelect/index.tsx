import { Select } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'

import Container from '../Container'
import Editor from './Editor'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class SelectItem extends React.Component<IProps> {
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
    const { model, formModel } = this.props
    const { id, defaultValues, options } = model

    return (
      <Container {...this.props}>
        <Select
          mode='multiple'
          defaultValue={defaultValues}
          value={formModel[id].values}
          onChange={this.onChange}>
          {options.map((option, index) => (
            <Select.Option key={index} value={option}>
              {option}
            </Select.Option>
          ))}
        </Select>
      </Container>
    )
  }

  private onChange = values => {
    const { formModel, model } = this.props
    const { id } = model

    formModel[id].values = values
  }
}
