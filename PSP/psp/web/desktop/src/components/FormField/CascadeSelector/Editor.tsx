import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, Select, message, Tooltip } from 'antd'

import Field from '@/domain/Applications/App/Field'
import { FormItem, Label } from '../style'
import BaseEditor from '../BaseEditor'
import { Http } from '@/utils'
import { Icon } from '@/components'
import { observable, toJS } from 'mobx'

interface IProps {
  formModel?: any
  model: Field
  appId?: string
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class CascadeSelectorEditor extends React.Component<{
  viewModel: any
  formModel: any
  appId?: string
}> {
  @observable
  resourceTypes = []

  @observable
  preCascadeSelectId = null

  oldID = this.props.viewModel.id

  // 更新预设值选项
  updateOptions = async (type, subType?) => {
    const { viewModel, appId } = this.props
    let url = `/app/schedulerResourceValue?app_id=${appId}&resource_type=${type}`
    if (subType) {
      url += `&resource_sub_type=${subType}`
    }
    let resValue = await Http.get(url)
    viewModel.options = resValue?.data?.items || []
  }

  // 动态获取资源类型选项
  updateResourceTypeOptions = async () => {
    const res = await Http.get('/app/schedulerResourceKey')
    this.resourceTypes = res.data.keys || []
  }

  async componentDidMount() {
    const { viewModel, formModel } = this.props

    await this.updateResourceTypeOptions()
    // 未设置 optionsFrom 的情况，默认获取资源类型的第一项
    if (
      !viewModel.optionsFrom ||
      viewModel.optionsFrom === 'script' ||
      viewModel.optionsFrom === 'custom'
    ) {
      viewModel.optionsFrom = this.resourceTypes[0]
    }

    // 获取上一级级联选择器 ID
    const customJSONValueString = JSON.parse(
      viewModel.customJSONValueString || '{}'
    )
    this.preCascadeSelectId =
      customJSONValueString?.preCascadeSelect?.id || null

    // 根据 optionsFrom 更新预设值选项
    await this.updateOptions(
      viewModel.optionsFrom,
      formModel?.[this.preCascadeSelectId]?.defaultValue
    )
  }

  async resourceTypeOnChange(value) {
    const { viewModel, formModel } = this.props
    viewModel.optionsFrom = value

    await this.updateOptions(
      value,
      formModel?.[this.preCascadeSelectId]?.defaultValue
    )
  }

  render() {
    const { viewModel, formModel } = this.props

    return (
      <>
        <FormItem>
          <Label>使用说明：</Label>
          <p style={{ width: 400, wordBreak: 'break-all' }}>
            关联选择器为业务表单项，选择器数据来源于定制的 API 数据，包括
            平台，资源，队列等数据
            可以指定上一级关联选择器，上一级关联选择器值的变化，会影响当前关联选择器选项的值。
          </p>
        </FormItem>
        <FormItem>
          <Label>选项数据来源：</Label>
          <Select
            value={viewModel.optionsFrom}
            onChange={(value: string) => this.resourceTypeOnChange(value)}>
            {this.resourceTypes.map((item, index) => (
              <Select.Option key={index} value={item}>
                {item}
              </Select.Option>
            ))}
          </Select>
        </FormItem>
        <FormItem>
          <Label>
            <div className='info'>
              <div className='label'>
                <span className='text'>预设：</span>
                <Tooltip
                  placement='left'
                  title={'由于数据是动态的，预设数据可能会无效'}>
                  <Icon className='help' type='help-circle' />
                </Tooltip>
              </div>
            </div>
          </Label>
          <Select
            value={viewModel.defaultValue}
            onChange={(value: string) => (viewModel.defaultValue = value)}>
            {viewModel.options.map((item, index) => (
              <Select.Option key={index} value={item.value}>
                {item.value + item.suffix}
              </Select.Option>
            ))}
          </Select>
        </FormItem>
        <FormItem>
          <Label>
            <div className='info'>
              <div className='label'>
                <span className='text'>上一级关联选择器：</span>
                <Tooltip
                  placement='left'
                  title={
                    '确保上一级关联选择器设置了预设值，这样当前的关联选择器预设才可能有候选项'
                  }>
                  <Icon className='help' type='help-circle' />
                </Tooltip>
              </div>
            </div>
          </Label>
          <Select
            allowClear
            value={this.preCascadeSelectId}
            onChange={async (value: string) => {
              if (!viewModel.id) {
                message.error('请填写ID，再进行选择')
                return
              }

              if (value) {
                const selectedJsonValue = JSON.parse(
                  formModel?.[value]?.customJSONValueString || {}
                ) as any

                // TODO 循环依赖检测 （while循环）
                if (selectedJsonValue?.preCascadeSelect?.id === viewModel.id) {
                  message.error(
                    `${selectedJsonValue?.preCascadeSelect?.id} 和 ${viewModel.id} 不能互为上一级关联选择器`
                  )
                  return
                }
              }

              this.preCascadeSelectId = value

              viewModel.customJSONValueString = JSON.stringify({
                preCascadeSelect: value ? toJS(formModel[value]) : null
              })

              // 更新当前选项，由于你选择了上级数据
              if (value && formModel[value].defaultValue) {
                await this.updateOptions(
                  viewModel.optionsFrom,
                  formModel[value].defaultValue
                )
                viewModel.defaultValue = null // 清理预设值
              } else {
                await this.updateOptions(viewModel.optionsFrom)
                viewModel.defaultValue = null // 清理预设值
              }
            }}>
            {Object.keys(formModel)
              .filter(item => {
                return (
                  item !== this.oldID &&
                  formModel?.[item]?.id !== undefined &&
                  item !== viewModel.id &&
                  formModel[item]?.type === 'cascade_selector'
                )
              })
              .map((item, index) => (
                <Select.Option key={index} value={item}>
                  {formModel[item]?.label}
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
    {({ viewModel, formModel, appId }) => (
      <CascadeSelectorEditor
        viewModel={viewModel}
        formModel={formModel}
        appId={appId}
      />
    )}
  </BaseEditor>
)
