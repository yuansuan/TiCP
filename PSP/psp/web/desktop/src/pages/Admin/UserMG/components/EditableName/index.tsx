import * as React from 'react'

import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { Input } from 'antd'

import { Icon } from '@/components'
import { EditableNameWrapper } from './style'

interface IProps {
  onClick?: any
  defaultValue?: string
  onConfirm?: (value: string) => void
  onCancel?: (value: string) => void
  className?: string
}

@observer
export default class EditableName extends React.Component<IProps> {
  @observable editing = false
  @action
  updateEditing = flag => (this.editing = flag)

  inputRef = null

  private edit = () => {
    this.updateEditing(true)
  }

  private onInput = e => {
    if (e.keyCode === 13) {
      this.onConfirm(e)
    } else if (e.keyCode === 27) {
      this.onCancel(e)
    }
  }

  private onConfirm = e => {
    const { value } = e.target
    const { onConfirm } = this.props

    onConfirm && onConfirm(value)

    this.updateEditing(false)
  }

  private onCancel = e => {
    const { value } = e.target
    const { onCancel } = this.props

    onCancel && onCancel(value)

    this.updateEditing(false)
  }

  render() {
    const { defaultValue = '', onClick, className } = this.props
    const { editing } = this

    return (
      <EditableNameWrapper className={className}>
        {editing ? (
          <Input
            ref={ref => (this.inputRef = ref)}
            autoFocus
            maxLength={64}
            defaultValue={defaultValue}
            onKeyDown={this.onInput}
            onBlur={this.onConfirm}
            onFocus={e => e.target.select()}
            size='small'
          />
        ) : (
          <div className='editor'>
            <div
              className={`style.value ${onClick ? 'isLink' : ''}`}
              onClick={onClick}>
              {defaultValue}
            </div>
            <div className='edit' onClick={this.edit}>
              <Icon type='edit' />
            </div>
          </div>
        )}
      </EditableNameWrapper>
    )
  }
}
