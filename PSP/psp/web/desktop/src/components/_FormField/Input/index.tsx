/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observer } from 'mobx-react'
import * as React from 'react'

import { ValidInput } from '@/components'
import Container from '../Container'
import { runInTyping } from '../utils'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class InputItem extends React.Component<IProps> {
  private validate = value => {
    return true
  }

  public render() {
    const { formModel, model } = this.props
    const { id, defaultValue, hidden } = model
    if (!formModel[id]) return null

    return (
      <Container {...this.props}>
        <ValidInput
          data-test={`FormField-Input-${id}`}
          disabled={hidden}
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

    runInTyping(formModel, () => (formModel[id].value = e.target.value))
  }
}
