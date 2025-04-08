import { Switch } from 'antd'
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
export default class CheckboxItem extends React.Component<IProps> {
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
    const { model, formModel } = this.props
    const { id } = model

    return (
      <Container {...this.props}>
        <Switch
          defaultChecked={model.defaultValue === 'true'}
          onChange={(checked) =>
            (formModel[id].value = checked ? 'true' : 'false')
          }
        />
      </Container>
    )
  }
}
