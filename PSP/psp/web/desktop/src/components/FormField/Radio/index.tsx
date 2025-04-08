import { Radio } from 'antd'
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
export default class RadioItem extends React.Component<IProps> {
  public static Editor = Editor

  constructor(props) {
    super(props)

    const { formModel, model } = props
    if (formModel) {
      formModel[model.id] = {
        ...model,
        value: model.value || model.defaultValue,
        values: model.values.length > 0 ? model.values : model.defaultValues
      }
    }
  }

  public render() {
    const { model } = this.props
    const { defaultValue, options } = model

    return (
      <Container {...this.props}>
        <Radio.Group defaultValue={defaultValue} onChange={this.onChange}>
          {options.map((option, index) => (
            <Radio key={index} value={option}>
              {option}
            </Radio>
          ))}
        </Radio.Group>
      </Container>
    )
  }

  private onChange = (e) => {
    const { formModel, model } = this.props
    const { id, defaultValue } = model

    const value = e.target.value || defaultValue
    formModel[id].value = value
  }
}
