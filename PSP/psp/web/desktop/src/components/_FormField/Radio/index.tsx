/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Radio } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'

import Container from '../Container'
import { runInTyping } from '../utils'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class RadioItem extends React.Component<IProps> {
  public render() {
    const { model, formModel } = this.props
    const { defaultValue, options, id } = model
    if (!formModel[id]) return null
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

  private onChange = e => {
    const { formModel, model } = this.props
    const { id, defaultValue } = model

    const value = e.target.value || defaultValue

    runInTyping(formModel, () => (formModel[id].value = value))
  }
}
