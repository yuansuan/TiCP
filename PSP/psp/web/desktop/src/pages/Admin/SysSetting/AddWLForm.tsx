import * as React from 'react'
import { Descriptions, Input, message, Select } from 'antd'
import { Label } from '@/components'
import { Button } from '@/components'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import { Http } from '@/utils'
import { isValidIp } from '@/utils/Validator'

const Wrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    flex-direction: column;
  }

  .formItem {
    width: 300px;
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
  margin-left: 10px;
  color: #f5222d;
`

interface IProps {
  onCancel: () => void
  onOk: () => void
}

@observer
export default class AddWLForm extends React.Component<IProps> {
  @observable loading = false
  @observable ip = ''
  @observable name = ''
  @observable ipErrMessage = ''
  @observable ipValidateErr = false
  @observable selectedName = []

  private inputRef = null

  constructor(props) {
    super(props)
    this.inputRef = React.createRef()
  }

  async componentDidMount() {
    this.inputRef.current.focus()
    const res = await Http.get('/sysconfig/whitelist/canadduser')
    this.selectedName = res.data?.map(item => item.name)
  }

  validate = (ip: string, name: string) => {
    if (ip.length === 0 && name.length === 0) {
      this.ipErrMessage = ''
      this.ipValidateErr = true
      return false
    } else if (ip.length === 0 && name.length !== 0) {
      this.ipErrMessage = ''
      this.ipValidateErr = false
      return true
    } else {
      if (isValidIp(ip)) {
        this.ipErrMessage = ''
        this.ipValidateErr = false
        return true
      } else {
        this.ipErrMessage = '请填写正确的IP'
        this.ipValidateErr = true
        return false
      }
    }
  }

  submit = async () => {
    if (this.validate(this.ip, this.name)) {
      this.loading = true

      try {
        const res = await Http.post('/sysconfig/whitelist', {
          ip: this.ip,
          name: this.name,
        })
        if (res.success) {
          message.success('添加规则成功')
          this.props.onOk()
        } else {
          message.error('添加规则失败')
        }
      } finally {
        this.loading = false
      }
    }
  }

  render() {
    return (
      <Wrapper>
        <Descriptions title='' column={1} style={{ margin: '0 0 50px 0' }}>
          <Descriptions.Item label={<Label>IP地址</Label>}>
            <div className='item'>
              <div>
                <Input
                  ref={this.inputRef}
                  style={{
                    borderColor: this.ipValidateErr ? '#f5222d' : 'inherit',
                  }}
                  className='formItem'
                  placeholder='请输入IP'
                  value={this.ip}
                  onChange={e => {
                    this.ip = e.target.value
                  }}
                />
                <Tips>{this.ipErrMessage}</Tips>
              </div>
            </div>
          </Descriptions.Item>
          <Descriptions.Item label={<Label>用户名</Label>}>
            <div className='item'>
              <Select
                allowClear
                showSearch
                value={this.name}
                defaultValue={this.name}
                className='formItem'
                onChange={value => {
                  this.name = value
                }}>
                {this.selectedName.map(key => {
                  return (
                    <Select.Option key={key} value={key}>
                      {key}
                    </Select.Option>
                  )
                })}
              </Select>
            </div>
          </Descriptions.Item>
        </Descriptions>
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
