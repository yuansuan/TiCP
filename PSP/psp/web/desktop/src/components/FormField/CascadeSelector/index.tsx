import { Select } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'
import { observable } from 'mobx'

import Container from '../Container'
import Editor from './Editor'
import { Http } from '@/utils'
import styled from 'styled-components'

const Wrapper = styled.div``

interface IProps {
  model
  formModel: any
  showId?: boolean
  appId?: string
}

@observer
export default class CascadeSelectorItem extends React.Component<IProps> {
  public static Editor = Editor
  @observable nodeNums = {}

  constructor(props) {
    super(props)

    const { formModel, model } = props
    if (formModel) {
      formModel[model.id] = {
        ...model,
        value: model.value || model.defaultValue,
        values: model.values.length > 0 ? model.values : model.defaultValues
      }
    }
  }

  updateSelfOptions = async () => {
    const { model, formModel, appId } = this.props

    const jsonValue = JSON.parse(model.customJSONValueString || '{}')

    // 如果关联了上级选择器，获取上级选择器的值
    if (jsonValue?.preCascadeSelect?.id) {
      const preCascadeSelect = formModel[jsonValue?.preCascadeSelect?.id]
      // 过滤多加一个条件
      if (preCascadeSelect.value) {
        let resValue = await Http.get(
          `/app/schedulerResourceValue?app_id=${appId}&resource_type=${model.optionsFrom}&resource_sub_type=${preCascadeSelect.value}`
        )
        model.options = resValue.data.items
      } else {
        model.options = []
      }
    } else {
      let resValue = await Http.get(
        `/app/schedulerResourceValue?app_id=${appId}&resource_type=${model.optionsFrom}`
      )
      model.options = resValue.data.items
    }
  }

  updateNextOptions = async (value, isClearValue) => {
    const { model, formModel, appId } = this.props
    // 遍历所有 cascade_selector 组建，如果自己是别人的上级，则根据默认值更新对应的下级数据
    const cascadeSelectorKeys = Object.keys(formModel).filter(item => {
      if (formModel[item].type === 'cascade_selector') {
        const { preCascadeSelect } = JSON.parse(
          formModel?.[item]?.customJSONValueString || '{}'
        )
        if (preCascadeSelect?.id === model.id) {
          return true
        } else {
          return false
        }
      } else {
        return false
      }
    })

    await Promise.all(
      cascadeSelectorKeys.map(async item => {
        // formModel[item]
        let resValue = await Http.get(
          `/app/schedulerResourceValue?app_id=${appId}&resource_type=${formModel[item].optionsFrom}&resource_sub_type=${value}`
        )
        formModel[item].options = resValue.data.items
        if (isClearValue) {
          formModel[item].value = null
        }
      })
    )
  }

  async componentDidMount() {
    const { model } = this.props

    if (model.optionsFrom !== 'script' || model.optionsFrom !== 'custom') {
      // 说明走的自定义 API
      await this.updateSelfOptions()
      // 第一次刷新next options，不能清除默认值
      this.updateNextOptions(model.value, false)
    }
  }

  public render() {
    const { model, formModel } = this.props
    const { id, defaultValue, options } = model

    return (
      <Container {...this.props}>
        <Wrapper>
          <Select
            defaultValue={defaultValue}
            value={formModel[id]?.value}
            onDropdownVisibleChange={this.onDropdownVisibleChange}
            onChange={this.onChange}>
            {options.map((option, index) => (
              <Select.Option key={option.value} value={option.value}>
                {option.value + option.suffix}
              </Select.Option>
            ))}
          </Select>
        </Wrapper>
      </Container>
    )
  }

  private onDropdownVisibleChange = open => {
    if (open) {
      this.updateSelfOptions()
    }
  }

  private onChange = async value => {
    const { formModel, model } = this.props
    const { id } = model
    formModel[id].value = value

    // 更新关联的 select 选项
    await this.updateNextOptions(value, true)
  }
}
