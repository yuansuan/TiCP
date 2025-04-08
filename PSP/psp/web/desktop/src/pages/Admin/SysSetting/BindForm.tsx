import * as React from 'react'
import { Descriptions, Input, message, Spin } from 'antd'
import { Label } from '@/components'
import { Button } from '@/components'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import { Http } from '@/utils'

const Wrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    flex-direction: column;
  }

  .formItem {
    width: 360px;
  }

  .ant-descriptions-item {
    display: flex;
  }

  .ant-descriptions-item-label {
    padding-top: 5px;
  }

  .footer {
    position: absolute;
    display: flex;
    bottom: 0px;
    right: 0;
    width: 100%;
    line-height: 64px;
    height: 64px;
    background: white;

    .footerMain {
      margin-left: auto;
      margin-right: 8px;

      button {
        margin: 0 8px;
      }
    }
  }
`

const Tips = styled.span`
  font-family: PingFangSC-Regular;
  font-size: 12px;
  color: #999999;
  line-height: 22px;
`

interface IProps {
  onCancel: () => void
  onOk: () => void
}

@observer
export default class BindForm extends React.Component<IProps> {
  @observable loading = false
  @observable err = false
  @observable token = ''

  private inputRef = null

  constructor(props) {
    super(props)
    this.inputRef = React.createRef()
  }

  componentDidMount() {
    this.inputRef.current.focus()
  }

  validate = (token: string) => {
    if (token.length !== 0) {
      this.err = false
      return true
    } else {
      this.err = true
      return false
    }
  }

  submit = async () => {
    if ([this.validate(this.token)].every(r => r)) {
      this.loading = true
      let res = null
      try {
        res = await Http.post(
          '/bindcloud/bind/',
          {
            token: this.token.trim(),
          },
          { timeout: 0 }
        )
        if (res.data.success) {
          message.info('绑定泛超算云成功, 开始同步泛超算云应用，请稍候')
          this.props.onOk()
          await Http.put('/bindcloud/sync-apps')
          message.success('泛超算云应用同步成功')
        } else {
          message.error('绑定泛超算云失败')
        }
      } finally {
        this.loading = false
      }
    }
  }

  render() {
    return (
      <Wrapper>
        <Spin spinning={this.loading} tip={'API Token 绑定中, 请稍后...'}>
          <Descriptions title='' column={1} style={{ margin: '0 0 50px 0' }}>
            <Descriptions.Item label={<Label required>API Token</Label>}>
              <div className='item'>
                <Input
                  ref={this.inputRef}
                  style={{
                    borderColor: this.err ? '#f5222d' : 'inherit',
                  }}
                  className='formItem'
                  placeholder=''
                  value={this.token}
                  onBlur={e => {
                    this.validate(e.target.value.trim())
                  }}
                  onChange={e => {
                    this.token = e.target.value
                  }}
                />
                <Tips>可以向远算工作人员申请 API Token</Tips>
              </div>
            </Descriptions.Item>
          </Descriptions>
        </Spin>
        <div className='footer'>
          <div className='footerMain'>
            <Button
              type='primary'
              disabled={this.loading}
              onClick={this.submit}>
              确认
            </Button>
            <Button onClick={this.props.onCancel}>取消</Button>
          </div>
        </div>
      </Wrapper>
    )
  }
}
