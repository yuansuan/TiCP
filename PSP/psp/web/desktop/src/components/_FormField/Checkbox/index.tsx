/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Switch } from 'antd'
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
export default class CheckboxItem extends React.Component<IProps> {
  public render() {
    const { model, formModel } = this.props
    const { id } = model
    if (!formModel[id]) return null

    return (
      <Container {...this.props}>
        <Switch
          defaultChecked={model.defaultValue === 'true'}
          onChange={checked =>
            runInTyping(
              formModel,
              () => (formModel[id].value = checked ? 'true' : 'false')
            )
          }
        />
      </Container>
    )
  }
}
