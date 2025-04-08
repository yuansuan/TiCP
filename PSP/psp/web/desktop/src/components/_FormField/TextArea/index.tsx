/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observer } from 'mobx-react'
import * as React from 'react'
import { Input, message } from 'antd'
import Container from '../Container'
import { runInTyping } from '../utils'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class TextAreaItem extends React.Component<IProps> {
  public render() {
    const { formModel, model } = this.props
    const { id, defaultValue, hidden } = model
    if (!formModel[id]) return null
    return (
      <Container {...this.props}>
        <Input.TextArea
          rows={10}
          maxLength={255}
          disabled={hidden}
          value={formModel[id].value || defaultValue}
          onChange={this.onChange}
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
