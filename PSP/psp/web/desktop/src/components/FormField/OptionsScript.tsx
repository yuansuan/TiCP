import * as React from 'react'
import { action } from 'mobx'
import { observer } from 'mobx-react'
import { IViewModel } from 'mobx-utils'
import { Input, Button, message } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from './style'
import { Http } from '@/utils'

interface IProps {
  viewModel: Field & IViewModel<Field>
}

@observer
export default class OptionsScriptUI extends React.Component<IProps> {
  @action
  getOption = async () => {
    const { viewModel } = this.props
    // 获取选项
    // 路径格式校验，不为空校验
    if (viewModel.optionsScript) {
      if (!/^(\/[\w-.]+)+$/.test(viewModel.optionsScript)) {
        message.error('请输入正确格式的脚本路径')
        return
      }
      const res = await Http.post('/application/options', {
        script: viewModel.optionsScript
      })
      viewModel.options = [...res.data]
    }
  }

  @action
  onChange = e => {
    const { viewModel } = this.props
    viewModel.optionsScript = e.target.value
  }

  render() {
    const { viewModel } = this.props

    return (
      <FormItem>
        <Label>脚本路径：</Label>
        <Input value={viewModel.optionsScript} onChange={this.onChange} />
        <Button onClick={this.getOption} type='link'>
          获取选项
        </Button>
      </FormItem>
    )
  }
}
