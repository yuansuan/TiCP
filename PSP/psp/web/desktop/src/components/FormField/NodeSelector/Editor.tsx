import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, Select } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import BaseEditor from '../BaseEditor'
import { Http } from '@/utils'

interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class NodeSelectorEditor extends React.Component<{ viewModel: any }> {
  async componentDidMount() {
    const { viewModel } = this.props
    // TODO 动态获取节点
    const res = await Http.get('/node/list')
    const { nodeInfos } = res?.data
    viewModel.options = nodeInfos?.map(n => n.node_name) || []
  }

  render() {
    const { viewModel } = this.props

    return (
      <>
        <FormItem>
          <Label>使用说明：</Label>
          <p style={{ width: 400, wordBreak: 'break-all' }}>
            节点选择器为业务表单项，针对用户作业提交，需要自行选择节点，并填写核数。
            使用该组件后，在作业提交时，产生环境变量:
            NODE_SELECTOR=NODE_NAME:NODE_CORES,NODE_NAME:NODE_CORES
            例如：NODE_SELECTOR=node1:12,node2:8
            运维人员可以在脚本中解析环境变量 NODE_SELECTOR
          </p>
        </FormItem>
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
             maxLength={64}
            value={viewModel.postText}
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
      </>
    )
  }
}

export default (props: IProps) => (
  <BaseEditor {...props}>
    {({ viewModel }) => <NodeSelectorEditor viewModel={viewModel} />}
  </BaseEditor>
)
