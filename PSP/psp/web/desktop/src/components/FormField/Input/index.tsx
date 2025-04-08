import { message } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'

import { ValidInput } from '@/components'
import Container from '../Container'
import Editor from './Editor'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class InputItem extends React.Component<IProps> {
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

  private validate = value => {
    const { id } = this.props.model

    // hardcode: add validation for JOB_NAME field
    if (
      id === 'JOB_NAME' &&
      value !== '' &&
      !/^[a-zA-Z0-9][\[\]a-zA-Z0-9_-]{0,127}$/.test(value)
    ) {
      message.error(
        '作业名称不支持中文符号，必须在128个字符以内，由字母或数字开头，包含字母，数字，下划线，中括号或中划线'
      )

      return false
    }

    return true
  }

  public render() {
    const { formModel, model } = this.props
    const { id, defaultValue, hidden } = model

    return (
      <Container {...this.props}>
       <ValidInput
          // disabled={hidden}
          defaultValue={defaultValue}
          value={formModel[id]?.value}
          onChange={this.onChange}
          validator={this.validate}
        />
      </Container>
    )
  }

  private onChange = e => {
    const { formModel, model } = this.props
    const { id } = model

    formModel[id].value = e.target.value
  }
}
