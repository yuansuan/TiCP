import React from 'react'
import { Checkbox, Tooltip } from 'antd'
import { observer } from 'mobx-react'

import { Button, Modal, Icon } from '@/components'
import { StyledToolbar } from './style'
import { observable } from 'mobx'

interface IProps {
  isReversed: boolean
  find: () => void
  fetchContent: (initial?: boolean) => Promise<any>
  freshen: () => void
  updateIsReversed: (flag: boolean) => void
  readOnly?: boolean
}

@observer
export default class Toolbar extends React.Component<IProps> {
  @observable isAutoRefrsh = false
  
  private intervalId = null

  private onFind = () => {
    this.props.find()
  }

  private reverse = () => {
    const { updateIsReversed, isReversed } = this.props

    updateIsReversed(!isReversed)
  }

  private freshen = () => {
    if (this.props.readOnly) {
      this.props.freshen()
      return
    }
    Modal.showConfirm({
      content: '确认要刷新文件吗？（刷新文件会丢失已编辑的内容）',
    }).then(() => {
      this.props.freshen()
    })
  }

  private onAutoRefreshCheck = (checked) => {
    this.isAutoRefrsh = checked
    
    if (checked) {
      this.props.freshen()
      this.intervalId && clearInterval(this.intervalId)

      this.intervalId = setInterval(() => {
        this.props.freshen()
      }, 5 * 1000)
    } else {
      clearInterval(this.intervalId)
    }
  }

  componentWillUnmount(): void {
    clearInterval(this.intervalId)
  }

  render() {
    const { isReversed } = this.props

    return (
      <StyledToolbar>
        <div>
          <Button icon='refresh' disabled={this.isAutoRefrsh} onClick={this.freshen}>
            刷新
          </Button>

          <Checkbox
            className='reverse'
            checked={this.isAutoRefrsh}
            onChange={e => this.onAutoRefreshCheck(e.target.checked)}>
            自动刷新
            <Tooltip
              title={
                '每5秒自动刷新内容, 注意：编辑模式下，启用自动刷新会导致编辑内容丢失'
              }>
              <Icon type={'help-circle'} />
            </Tooltip>
          </Checkbox>

          <Checkbox
            className='reverse'
            checked={isReversed}
            onClick={this.reverse}>
            跳转到底部
            <Tooltip
              title={
                '提前获取该文件底部内容 (如果文件内容大于1M, 则提前获取文件尾部的1M内容)'
              }>
              <Icon type={'help-circle'} />
            </Tooltip>
          </Checkbox>
        </div>

        <div className='search'>
          <Button type='primary' icon='search' ghost onClick={this.onFind}>
            查询
          </Button>
        </div>
      </StyledToolbar>
    )
  }
}
