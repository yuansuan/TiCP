import { Select } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'
import { observable } from 'mobx'

import Container from '../Container'
import Editor from './Editor'
import { Http } from '@/utils'
import { InputNumber } from 'antd'
import styled from 'styled-components'

const Wrapper = styled.div`
  
`


const InputNumberWrapper = styled.div`
  margin: 10px 5px;
  display: flex;
  align-items: center;
  justify-content: flex-start;

  .label {
    width: 80px;
  }

  .ant-input-number-sm input {
    height: 22px;
    line-height: 22px;
  }
  .unit {
    padding: 5px;
  }
`

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class NodeSelectorItem extends React.Component<IProps> {
  public static Editor = Editor
  @observable nodeNums = {}

  constructor(props) {
    super(props)

    const { formModel, model } = props
    if (formModel) {
      formModel[model.id] = {
        ...model,
        value: model.value || model.defaultValue,
        values: model.values.length > 0 ? model.values : model.defaultValues,
      }
    }
  }

  async componentDidMount() {
    const { model } = this.props
    this.nodeNums = JSON.parse(model.customJSONValueString || '{}')
    // TODO 动态获取节点
    const res = await Http.get('/node/list')
    const { nodeInfos } = res?.data
    model.options = nodeInfos?.map(n => n.node_name) || []
  }

  public render() {
    const { model, formModel } = this.props
    const { id, defaultValues, options, customJSONValueString } = model
    const jsonValue = JSON.parse(customJSONValueString || '{}')

    return (
      <Container {...this.props}>
        <Wrapper>
          <Select
            mode='multiple'
            defaultValue={defaultValues}
            value={formModel[id].values}
            placeholder={'请选择节点'}
            onChange={this.onChange}>
            {options.map((option, index) => (
              <Select.Option key={index} value={option}>
                {option}
              </Select.Option>
            ))}
          </Select>
          {
            formModel[id]?.values?.map((option, index) => (
              <InputNumberWrapper key={index} >
                <div className="label">{option}</div>
                <InputNumber 
                  size="small"
                  style={{width: 180}}
                  min={0}
                  precision={0}
                  value={this.nodeNums[option] || 0} 
                  defaultValue={jsonValue[option] || 0} 
                  onChange={value => this.onNodeNumChange(option, value)}
                  />
                  <div className="unit">核</div>
              </InputNumberWrapper>  
            ))
          }
        </Wrapper>
      </Container>
    )
  }

  private onNodeNumChange = (option, value) => {
    this.nodeNums[option] = value
    const { formModel, model } = this.props
    const { id } = model

    formModel[id].customJSONValueString = JSON.stringify(this.nodeNums)
  }

  private onChange = values => {
    const { formModel, model } = this.props
    const { id } = model

    formModel[id].values = values

    Object.keys(this.nodeNums).forEach(key => { if (!values.includes(key)) this.nodeNums[key] = 0 })

    formModel[id].customJSONValueString = JSON.stringify(this.nodeNums)
  }
}
