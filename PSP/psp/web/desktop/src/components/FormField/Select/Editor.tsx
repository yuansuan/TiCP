import { Input, Select, Switch, Radio, message } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'

import Field from '@/domain/Applications/App/Field'
import { Modal } from '@/components'
import BaseEditor from '../BaseEditor'
import Options from '../Options'
import OptionsScript from '../OptionsScript'
import { FormItem, Label } from '../style'

const scriptPathReg = /^(\/[\w-.]+)+$/
interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class SelectEditor extends React.Component<{ viewModel: any }> {
  public render() {
    const { viewModel } = this.props

    return (
      <>
        <FormItem>
          <Label>选项数据来源：</Label>
          <Radio.Group
            onChange={e => {
              // 切换警告
              Modal.showConfirm({
                title: '确认',
                content: '切换选项数据来源，可能会清除掉之前的选项和预设'
              }).then(() => {
                viewModel.optionsFrom = e.target.value
                // 重置选项 和 optionsScript
                viewModel.options = []
                viewModel.defaultValue = ''
                // viewModel.optionsScript = ''
              })
            }}
            value={viewModel.optionsFrom}>
            <Radio value={'custom'}>自定义</Radio>
            <Radio value={'script'}>脚本</Radio>
          </Radio.Group>
        </FormItem>
        {viewModel.optionsFrom === 'custom' ? (
          <Options viewModel={viewModel} />
        ) : (
          <OptionsScript viewModel={viewModel} />
        )}
        <FormItem>
          <Label>预设：</Label>
          <Select
            value={viewModel.defaultValue}
            onChange={(value: string) => (viewModel.defaultValue = value)}>
            <Select.Option key={-1} value={''}>
              &lt;空选项&gt;
            </Select.Option>
            {viewModel.options.map((item, index) => (
              <Select.Option key={index} value={item}>
                {item}
              </Select.Option>
            ))}
          </Select>
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

export default (props: IProps) => {
  const onConfirm = viewModel => {
    if (
      viewModel.optionsScript !== '' &&
      !scriptPathReg.test(viewModel.optionsScript)
    ) {
      message.error('请输入正确格式的脚本路径')
      return
    }

    props.onConfirm(viewModel)
  }

  const onCancel = viewModel => {
    viewModel.reset()
    props.onCancel(viewModel)
  }

  return (
    <BaseEditor {...{ ...props, onConfirm, onCancel }}>
      {({ viewModel }) => <SelectEditor viewModel={viewModel} />}
    </BaseEditor>
  )
}
