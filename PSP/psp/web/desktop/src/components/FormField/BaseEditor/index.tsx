/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Button } from '@/components'
import { Input } from 'antd'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { createViewModel, IViewModel } from 'mobx-utils'
import * as React from 'react'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import { Footer, FormItemWrapper } from './style'

interface IProps {
  formModel?: any
  model: Field
  children?: any
  appId?: string
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
export default class BaseEditor extends React.Component<IProps, any> {
  @observable public viewModel: Field & IViewModel<Field>

  constructor(props: IProps) {
    super(props)

    this.viewModel = createViewModel(props.model)
  }

  public render() {
    const { viewModel } = this
    const { children, formModel, appId } = this.props

    return (
      <>
        <FormItemWrapper>
          <FormItem>
            <Label>
              <span className='required'>*</span>字段：
            </Label>
            <Input
              autoFocus
              maxLength={64}
              onFocus={e => e.target.select()}
              value={viewModel.label}
              onChange={e => (viewModel.label = e.target.value)}
            />
          </FormItem>
          <FormItem>
            <Label>
              <span className='required'>*</span>ID：
            </Label>
            <Input
              maxLength={64}
              value={viewModel.id}
              onChange={e => (viewModel.id = e.target.value)}
            />
          </FormItem>
          {children && children({ viewModel, formModel, appId })}
          <Footer>
            <Button onClick={this.onCancel}>取消</Button>
            <Button type='primary' onClick={this.onConfirm}>
              确定
            </Button>
          </Footer>
        </FormItemWrapper>
      </>
    )
  }

  private onCancel = () => {
    const { onCancel } = this.props
    const { viewModel } = this

    if (onCancel) {
      onCancel(viewModel)
    } else {
      viewModel.reset()
    }
  }

  private onConfirm = () => {
    const { viewModel } = this
    const { onConfirm } = this.props

    onConfirm(viewModel)
  }
}
