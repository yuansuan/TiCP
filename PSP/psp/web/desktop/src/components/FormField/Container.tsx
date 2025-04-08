import { Icon } from '@/components'
import { Tooltip } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'

import { FormItem, InfoText, Label } from './style'

interface IProps {
  model
  showId?: boolean
}

@observer
export default class ItemContainer extends React.Component<IProps, any> {
  public render() {
    const { showId, model, children } = this.props

    return (
      <FormItem>
        <Label>
          <div className='info'>
            <div className='label'>
              <span className='text' title={model.label}>
                {model.required && <span className='required'>*</span>}
                {model.label}
              </span>
              {model.help && (
                <Tooltip placement='left' title={model.help}>
                  <Icon className='help' type='help-circle' />
                </Tooltip>
              )}
            </div>
            <div className='id' title={`(ID:${model.id})`}>
              {showId ? <span className='value'>(ID:{model.id})</span> : null}
            </div>
          </div>
        </Label>
        {children}
        {model.postText && <InfoText>（{model.postText}）</InfoText>}
      </FormItem>
    )
  }
}
