/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Select, Tag } from 'antd'
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
  tagRender(props) {
    const { label, closable, onClose } = props

    return (
      <Tag
        title={label}
        closable={closable}
        onClose={onClose}
        style={{
          maxWidth: '100%',
          marginTop: 2,
          marginRight: 4,
          marginBottom: 2,
          padding: '0 4px 0 8px',
          display: 'inline-flex',
          alignItems: 'center',
        }}>
        <span
          style={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
            display: 'inline-block',
            maxWidth: 'calc(100% - 6px)',
          }}>
          {label}
        </span>
      </Tag>
    )
  }

  public render() {
    const { model, formModel } = this.props
    const { id, defaultValues, options } = model
    if (!formModel[id]) return null
    return (
      <Container {...this.props}>
        <Select
          mode='multiple'
          showArrow={true}
          tagRender={this.tagRender}
          defaultValue={defaultValues}
          value={formModel[id].values}
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

  private onChange = values => {
    const { formModel, model } = this.props
    const { id } = model

    runInTyping(formModel, () => (formModel[id].values = values))
  }
}
