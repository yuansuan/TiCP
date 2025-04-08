import * as React from 'react'
import { observer } from 'mobx-react'

import Editor from './Editor'
import Container from '../Container'

interface IProps {
  model
  showId?: boolean
}

@observer
export default class LabelItem extends React.Component<IProps> {
  static Editor = Editor

  render() {
    const { model, showId } = this.props

    return (
      <Container model={model} showId={showId}>
        <label>{model.label}</label>
      </Container>
    )
  }
}
