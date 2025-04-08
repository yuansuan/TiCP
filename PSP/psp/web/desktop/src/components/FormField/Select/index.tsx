import { Select } from 'antd'
import { observer } from 'mobx-react'
import { Modal } from '@/components'
import * as React from 'react'
import { CloseCircleOutlined } from '@ant-design/icons'
import Container from '../Container'
import Editor from './Editor'
import { Http, getSearchParamByKey } from '@/utils'
import { clusterCores } from '@/domain/ClusterCores'

interface IProps {
  model
  formModel: any
  showId?: boolean
}

@observer
export default class SelectItem extends React.Component<IProps> {
  public static Editor = Editor
  private isBindCloud = false

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

  async componentDidMount() {
    // const [_, res] = await Promise.all([
    //   clusterCores.getClusterCoreInfo(),
    //   Http.get('/bindcloud/status')
    // ])

    // this.isBindCloud = res.data.isBind

    const { model } = this.props
    if (model.optionsFrom === 'script') {
      const res = await Http.post('/application/options', {
        script: model.optionsScript
      })
      model.options = [...res.data]
    }
  }

  public render() {
    const { model, formModel } = this.props
    const { id, defaultValue, options } = model

    return (
      <Container {...this.props}>
        <Select
          defaultValue={defaultValue}
          value={formModel[id]?.value}
          allowClear={defaultValue === '' ? true : false}
          clearIcon={<CloseCircleOutlined title={'清除选项'} />}
          onChange={this.onChange}>
          {options.map((option, index) => (
            <Select.Option key={index} value={option}>
              {option}
            </Select.Option>
          ))}
        </Select>
      </Container>
    )
  }

  private validate = value => {
    // const { id } = this.props.model

    // hardcode: add validation for NUM_CPU field
    // if (id === 'NUM_CPU' && value !== '' && this.isBindCloud) {
    //   let num = Number(value)

    //   if (num > clusterCores.available_cores) {
    //     return false
    //   } else {
    //     return true
    //   }
    // }

    return true
  }

  private onChange = async value => {
    // 校验不通过
    if (!this.validate(value)) {
      try {
        await Modal.showConfirm({
          title: '确认',
          content:
            '当前所选核数大于本地集群剩余的可用核数, 会导致作业等待, 是否考虑进入【云端应用】进行作业提交, 点击【确认】按钮，进入【云端应用】页面，点击【取消】按钮，则继续当前作业操作 ？'
        })

        const appName = getSearchParamByKey(location.hash.split('?')[1], 'app')
        // go to cloud app to submit job
        // history.push(`/yscloudapps?app=${appName}&pushType=${HISTORY_PUSH_TYPE.NOT_NAV}`)
      } catch (e) {}
    }

    const { formModel, model } = this.props
    const { id, defaultValue } = model

    formModel[id].value = value || defaultValue
  }
}
