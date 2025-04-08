import * as React from 'react'
import { Input, Switch } from 'antd'
import { observer } from 'mobx-react'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import BaseEditor from '../BaseEditor'

interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class CheckboxEditor extends React.Component<{ viewModel: any }> {
  render() {
    const { viewModel } = this.props

    return (
      <>
        <FormItem>
          <Label>预设：</Label>
          <Switch
            checked={viewModel.defaultValue === 'true'}
            onChange={checked =>
              (viewModel.defaultValue = checked ? 'true' : 'false')
            }
          />
        </FormItem>
        <FormItem>
          <Label>帮助说明：</Label>
          <Input.TextArea
            value={viewModel.help}
            maxLength={255}
            onChange={e => (viewModel.help = e.target.value)}
          />
        </FormItem>
      </>
    )
  }
}

export default (props: IProps) => (
  <BaseEditor {...props}>
    {({ viewModel }) => <CheckboxEditor viewModel={viewModel} />}
  </BaseEditor>
)
