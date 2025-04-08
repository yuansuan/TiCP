/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Select } from 'antd'
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
export default class SelectItem extends React.Component<IProps> {
  public render() {
    const { model, formModel } = this.props
    const { id, defaultValue, options } = model
    if (!formModel[id]) return null
    return (
      <Container {...this.props}>
        <Select
          defaultValue={defaultValue}
          value={formModel[id]?.value}
          onChange={this.onChange}>
          {options.map((option, index) => (
            <Select.Option title={option} key={index} value={option}>
              {option}
            </Select.Option>
          ))}
        </Select>
      </Container>
    )
  }

  private onChange = value => {
    const { formModel, model } = this.props
    const { id, defaultValue } = model

    runInTyping(formModel, () => (formModel[id].value = value || defaultValue))
  }
}
