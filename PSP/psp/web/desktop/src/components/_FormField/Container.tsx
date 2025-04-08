/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Tooltip } from 'antd'
import { QuestionCircleFilled } from '@ant-design/icons'
import { observer } from 'mobx-react'
import React from 'react'
import { FormItem, Label } from './style'

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
                <Tooltip placement='top' title={model.help}>
                  <QuestionCircleFilled
                    style={{ margin: '0 2px', color: 'rgba(0,0,0,0.45)' }}
                  />
                </Tooltip>
              )}
              :
            </div>
            <div className='id' title={`(ID:${model.id})`}>
              {showId ? <span className='value'>(ID:{model.id})</span> : null}
            </div>
          </div>
        </Label>
        {children}
        {model.postText && <span className='post-text'>{model.postText}</span>}
      </FormItem>
    )
  }
}
