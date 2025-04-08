import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, Select } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import BaseEditor from '../BaseEditor'
import Options from '../Options'

interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class MultiSelectEditor extends React.Component<{ viewModel: any }> {
  render() {
    const { viewModel } = this.props

    return (
      <>
        <Options viewModel={viewModel} />
        <FormItem>
          <Label>预设：</Label>
          <Select
            mode='multiple'
            value={viewModel.defaultValues}
            onChange={(values: string[]) => (viewModel.defaultValues = values)}>
            {viewModel.options.map((item, index) => (
              <Select.Option key={index} value={item}>
                {item}
              </Select.Option>
            ))}
          </Select>
        </FormItem>
        <FormItem>
          <Label>右侧说明文字：</Label>
          <Input
            value={viewModel.postText}
            maxLength={64}
            onChange={e => (viewModel.postText = e.target.value)}
          />
        </FormItem>
        <FormItem>
          <Label>帮助说明：</Label>
          <Input.TextArea
           maxLength={255}
            value={viewModel.help}
            onChange={e => (viewModel.help = e.target.value)}
          />
        </FormItem>
        <FormItem>
          <Label>是否必填：</Label>
          <Switch
            checked={viewModel.required}
            onChange={checked => (viewModel.required = checked)}
          />
        </FormItem>
        <FormItem>
          <Label>是否隐藏：</Label>
          <Switch
            checked={viewModel.hidden}
            onChange={checked => (viewModel.hidden = checked)}
          />
        </FormItem>
      </>
    )
  }
}

export default (props: IProps) => (
  <BaseEditor {...props}>
    {({ viewModel }) => <MultiSelectEditor viewModel={viewModel} />}
  </BaseEditor>
)
