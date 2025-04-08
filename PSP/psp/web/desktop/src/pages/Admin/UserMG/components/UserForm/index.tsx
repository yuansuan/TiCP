import * as React from 'react'
import { Descriptions, Input, Select, Tooltip, Radio } from 'antd'
import { Label } from '@/components'
import { currentUser } from '@/domain'
import { Icon } from '@/components'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { User } from '@/domain/UserMG'
import moment from 'moment'

const Wrapper = styled.div`
  padding: 10px;
  margin-left: 150px;
  .formItem {
    width: 200px;
  }
`

const Msg = styled.span`
  padding: 0 10px;
  color: red;
`

interface IProps {
  user: User
  messages: any
  onChange: (type: string, value: any, isLDAP: boolean) => void
  onBlur?: (type: string, value: any) => void
  isEdit?: boolean
}

@observer
export default class UserForm extends React.Component<IProps> {
  private inputRef = null

  constructor(props) {
    super(props)
    this.inputRef = React.createRef()
  }

  componentDidMount() {
    this.inputRef.current.focus()
  }

  disabledDate = current => {
    return current && current < moment(new Date()).add(-1, 'days')
  }

  render() {
    const { user, onChange, onBlur, messages, isEdit } = this.props
    return (
      <Wrapper>
        <Descriptions title='' column={1}>
          <Descriptions.Item label={<Label required>登录名称</Label>}>
            <Input
              disabled={isEdit}
              ref={this.inputRef}
              className='formItem'
              maxLength={64}
              placeholder='请输入登录名称'
              value={user.name}
              onChange={e => {
                onChange('name', e.target.value, false)
              }}
              onBlur={e => {
                const msg = onBlur('name', e.target.value)
                messages.name = msg
              }}
            />
            <Tooltip title='登录名称只能包含字母数字下划线, 最大长度不能超过64位'>
              <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
            </Tooltip>
            <Msg>{messages.name}</Msg>
          </Descriptions.Item>

          <Descriptions.Item label={<Label required>密码</Label>}>
            <Input.Password
              className='formItem'
              placeholder='请输入密码'
              value={user.password}
              onChange={e => {
                onChange('password', e.target.value, false)
              }}
              onBlur={e => {
                const msg = onBlur('password', e.target.value)
                messages.password = msg
              }}
            />

            <Msg>{messages.password}</Msg>
          </Descriptions.Item>
          {/* TODO 需求待确定, 临时注释 */}
          {/* <Descriptions.Item label={<Label>状态</Label>}>
            <Radio.Group
              value={user.enabled}
              onChange={e => {
                onChange('enabled', e.target.value, false)
              }}>
              <Radio value={true}>启用（默认）</Radio>
              <Radio value={false}>禁用</Radio>
            </Radio.Group>
          </Descriptions.Item> */}

          <Descriptions.Item label={<Label>电话</Label>}>
            <Input
              className='formItem'
              placeholder='请输入电话'
              value={user.mobile}
              onChange={e => {
                onChange('mobile', e.target.value, false)
              }}
              onBlur={e => {
                const msg = onBlur('mobile', e.target.value)
                messages.mobile = msg
              }}
            />
            <Tooltip title='电话: 目前仅支持手机号码。'>
              <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
            </Tooltip>
            <Msg>{messages.mobile}</Msg>
          </Descriptions.Item>
          <Descriptions.Item label={<Label>邮箱</Label>}>
            <Input
              className='formItem'
              placeholder='请输入邮件'
              value={user.email}
              onChange={e => {
                onChange('email', e.target.value, false)
              }}
              onBlur={e => {
                const msg = onBlur('email', e.target.value)
                messages.email = msg
              }}
            />
            <Msg>{messages.email}</Msg>
          </Descriptions.Item>

          {
            currentUser?.isOpenapiSwitchEnable && (
              <Descriptions.Item label={<Label>允许调用openapi</Label>}>
                <Radio.Group
                  value={user.enable_openapi}
                  onChange={e => {
                    onChange('enable_openapi', e.target.value, false)
                  }}>
                  <Radio value={false}>禁用（默认）</Radio>
                  <Radio value={true}>启用</Radio>
                </Radio.Group>
              </Descriptions.Item>
            )
          }


        </Descriptions>
      </Wrapper>
    )
  }
}
