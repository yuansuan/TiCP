import { DatePicker } from 'antd'
import { observer } from 'mobx-react'
import moment from 'moment'
import * as React from 'react'

import Field from '@/domain/Applications/App/Field'
import BaseEditor from '../BaseEditor'
import { FormItem, Label } from '../style'

const dateFormat = 'YYYY-MM-DD HH:mm'

interface IProps {
  model: Field
  onCancel?: (viewModel?: any) => void
  onConfirm: (viewModel?: any) => void
}

@observer
class DateEditor extends React.Component<{ viewModel: any }> {
  public render() {
    const { viewModel } = this.props

    return (
      <>
        <FormItem>
          <Label>预设：</Label>
          <DatePicker
            showToday={false}
            style={{ width: '300px' }}
            {...(viewModel.defaultValue
              ? { defaultValue: moment(viewModel.defaultValue) }
              : {})}
            showTime={{ format: 'HH:mm' }}
            format={dateFormat}
            onChange={date =>
              (viewModel.defaultValue = date.format(dateFormat))
            }
          />
        </FormItem>
      </>
    )
  }
}

export default (props: IProps) => (
  <BaseEditor {...props}>
    {({ viewModel }) => <DateEditor viewModel={viewModel} />}
  </BaseEditor>
)
