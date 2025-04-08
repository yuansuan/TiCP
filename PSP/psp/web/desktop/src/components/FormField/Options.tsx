import * as React from 'react'
import { action } from 'mobx'
import { IViewModel } from 'mobx-utils'
import { Icon } from '@/components'
import { Input } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label, Options } from './style'

interface IProps {
  viewModel: Field & IViewModel<Field>
}

export default class OptionsUI extends React.Component<IProps> {
  @action
  changeOption = (index, value) => {
    const { viewModel } = this.props
    const options = [...viewModel.options]
    options.splice(index, 1, value)

    viewModel.options = options
  }

  @action
  addOption = () => {
    const { viewModel } = this.props
    viewModel.options = [...viewModel.options, '']
  }

  @action
  removeOption = index => {
    const { viewModel } = this.props
    const options = [...viewModel.options]
    options.splice(index, 1)

    viewModel.options = options
  }

  render() {
    const { viewModel } = this.props

    return (
      <FormItem>
        <Label>选项：</Label>
        <Options>
          {viewModel.options.map((item, index) => (
            <div key={index} className='input-wrapper'>
              <Input
                autoFocus
                value={item}
                onChange={e => this.changeOption(index, e.target.value)}
              />
              <div className='right-option'>
                <Icon type='close' onClick={() => this.removeOption(index)} />
              </div>
            </div>
          ))}
          <div className='new' onClick={this.addOption}>
            新建选项
          </div>
        </Options>
      </FormItem>
    )
  }
}
