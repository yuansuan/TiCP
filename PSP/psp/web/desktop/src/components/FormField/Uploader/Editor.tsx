import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, Radio, Tooltip } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import BaseEditor from '../BaseEditor'
import { Icon } from '@/components'

interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class UploaderEditor extends React.Component<{ viewModel: any }> {
  render() {
    const { viewModel } = this.props
    const { isMasterSlave } = viewModel

    if (!viewModel.fileFromType) {
      viewModel.fileFromType = 'local'
    }

    return (
      <>
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
          <Label>上传类型：</Label>
          <Radio.Group
            options={[
              { label: '本地文件', value: 'local' },
              { label: '服务器文件', value: 'server' },
              { label: '本地文件和服务器文件', value: 'local_server' }
            ]}
            value={viewModel.fileFromType}
            onChange={e => (viewModel.fileFromType = e.target.value)}
          />
        </FormItem>
        {/** 默认只支持一种方式，带主文件方式*/}
        {/* <FormItem>
          <Label>是否支持主文件(与老的主从文件模式互斥): </Label>
          <Switch
            disabled={isMasterSlave}
            checked={viewModel.isSupportMaster}
            onChange={checked => (viewModel.isSupportMaster = checked)}
          />
        </FormItem>
        <FormItem>
          <Label>是否支持主从文件选择：</Label>
          <Switch
            disabled={isSupportMaster}
            checked={viewModel.isMasterSlave}
            onChange={checked => (viewModel.isMasterSlave = checked)}
          />
        </FormItem> */}
        {isMasterSlave ? (
          <FormItem>
            <Label>
              <span className='required'>*</span>从文件关键字：
            </Label>
            <Input
              style={{ width: '350px' }}
              maxLength={64}
              placeholder='多个关键字使用;分隔'
              value={viewModel.masterIncludeKeywords}
              onChange={e => (viewModel.masterIncludeKeywords = e.target.value)}
            />
          </FormItem>
        ) : null}
        {isMasterSlave ? (
          <FormItem>
            <Label>
              <span className='required'>*</span>从文件后缀名：
            </Label>
            <Input
              style={{ width: '350px' }}
              maxLength={64}
              placeholder='多个后缀名使用;分隔'
              value={viewModel.masterIncludeExtensions}
              onChange={e =>
                (viewModel.masterIncludeExtensions = e.target.value)
              }
            />
          </FormItem>
        ) : null}
        <FormItem>
          <Label>
            是否选择工作目录
            <Tooltip placement='left' title={'重提交情况下，该选项禁用'}>
              <Icon className='help' type='help-circle' />
            </Tooltip>
            ：
          </Label>
          <Switch
            checked={viewModel.isSupportWorkdir}
            onChange={checked => (viewModel.isSupportWorkdir = checked)}
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
    {({ viewModel }) => <UploaderEditor viewModel={viewModel} />}
  </BaseEditor>
)
