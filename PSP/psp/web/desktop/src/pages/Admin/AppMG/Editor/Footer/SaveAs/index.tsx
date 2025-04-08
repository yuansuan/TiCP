import * as React from 'react'
import { Input, message } from 'antd'
import { observable, action } from 'mobx'
import { observer } from 'mobx-react'

import { Wrapper } from './style'
import { App } from '@/domain/Applications'
import { Button } from '@/components'

interface IProps {
  app: App
  defaultName?: string
  getScriptData: () => string
  onCancel: (next?: any) => void
  onOk: (next?: any) => void
}

@observer
export default class SaveAs extends React.Component<IProps> {
  @observable version = ''
  @action
  updateVersion = version => (this.version = version)
  @observable loading = false
  updateLoading = loading => (this.loading = loading)

  constructor(props) {
    super(props)

    if (props.defaultName) {
      this.updateVersion(props.defaultName)
    }
  }

  private onCancel = () => {
    const { onCancel } = this.props

    onCancel()
  }

  private onOk = async () => {
    const { onOk, app, getScriptData } = this.props

    if (!this.version) {
      message.error('版本不能为空')
      return
    }

    try {
      this.updateLoading(true)
      app.setScriptData(getScriptData())
      await app.saveAs(this.version).then(() => {
        message.success('另存为成功')
      })
      onOk(this.version)
    } finally {
      this.updateLoading(false)
    }
  }

  render() {
    const { loading } = this
    return (
      <Wrapper>
        <div className='body'>
          <div>
            <span>版本名称：</span>
            <Input
              autoFocus
              maxLength={64}
              onFocus={e => e.target.select()}
              value={this.version}
              onChange={e => this.updateVersion(e.target.value)}
            />
          </div>
        </div>
        <div className='footer'>
          <div className='footerMain'>
            <Button loading={loading} onClick={this.onCancel}>
              取消
            </Button>
            <Button loading={loading} type='primary' onClick={this.onOk}>
              确认
            </Button>
          </div>
        </div>
      </Wrapper>
    )
  }
}
