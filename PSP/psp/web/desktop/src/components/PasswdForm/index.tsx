import * as React from 'react'
import styled from 'styled-components'
import { Input, Descriptions, message } from 'antd'
import { Modal } from '@/components'
import { Label, PasswordStrengthChecker } from '@/components'
import { sysConfig, currentUser } from '@/domain'
import passwdChecker from '@/utils/passwdChecker'
import { observer } from 'mobx-react'
import { observable, computed } from 'mobx'

// 密码只能包含以下字符
const pwdReg = /[^A-Za-z0-9!#$%&*]/

const Msg = styled.span`
  padding: 0 0 0 10px;
  color: red;
`

const tipPasswdKey = {
  Number: '数字',
  Char_Upper: '大写字母',
  Char_Lower: '小写字母',
  Char_Special: '特殊字符'
}

@observer
class PasswdForm extends React.Component<any> {
  @observable oldPasswd = null
  @observable newPasswd = null
  @observable passwdAgain = null

  @observable oldPasswdMsg = null
  @observable newPasswdMsg = null
  @observable passwdAgainMsg = null

  get passwdChecker() {
    const { userConfig } = sysConfig
    passwdChecker.config({
      strengthCheck: true,
      maxLength: userConfig.password?.max_len,
      minLength: userConfig.password?.min_len
    })

    return passwdChecker
  }

  get tips() {
    const { userConfig } = sysConfig
    const { strength, min_len, max_len } = userConfig.password

    let strengthTips =
      strength.length === 0
        ? ''
        : `，至少包含${strength.map(s => tipPasswdKey[s]).join('、')}`

    return `密码至少${min_len}位，最长${max_len}位${strengthTips}`
  }

  onChange = (type, value) => {
    this[type] = value
  }

  validate = (type, value) => {
    const { userConfig } = sysConfig
    const { strength } = userConfig.password

    if (type === 'oldPasswd') {
      if (!value) {
        this.oldPasswdMsg = '原密码不能为空'
      } else {
        this.oldPasswdMsg = null
      }
    }
    if (type === 'newPasswd') {
      if (!value) {
        this.newPasswdMsg = '新密码不能为空'
      } else {
        if (pwdReg.test(value)) {
          this.newPasswdMsg = '密码只能包含数字，大小写字母和特殊字符 !#$%&*'
        } else {
          const res = this.passwdChecker.test(value)
          const requireErrors = Object.values(res.requiredTestErrors)
          const strengthErrorKeys = Object.keys(
            res.strengthCheckTestErrors
          ).filter(k => strength.includes(k))

          if (requireErrors.length !== 0) {
            this.newPasswdMsg = requireErrors[0]
          } else if (strengthErrorKeys.length !== 0) {
            this.newPasswdMsg =
              res.strengthCheckTestErrors[strengthErrorKeys[0]]
          } else {
            this.newPasswdMsg = null
          }
        }
      }
    }
    if (type === 'passwdAgain') {
      if (!value) {
        this.passwdAgainMsg = '确认密码不能为空'
      } else if (value !== this.newPasswd) {
        this.passwdAgainMsg = '确认密码与新密码不一致'
      } else {
        this.passwdAgainMsg = null
      }
    }
  }

  @computed
  get strength() {
    const res = this.newPasswd ? this.passwdChecker.test(this.newPasswd) : {}
    const { strengthCheckTestsPassed } = res

    let strength = ''
    // strengthCheckTestsPassed 如果小于等于 2 密码强度弱，大于 2 且小于等于 3 密码强度中， 大于3 密码强度强
    if (strengthCheckTestsPassed <= 2) {
      strength = 'weak'
    } else if (strengthCheckTestsPassed <= 3) {
      strength = 'medium'
    } else if (strengthCheckTestsPassed >= 4) {
      strength = 'strong'
    }

    return strength
  }

  validateAll = () => {
    this.validate('oldPasswd', this.oldPasswd)
    this.validate('newPasswd', this.newPasswd)
    this.validate('passwdAgain', this.passwdAgain)
  }

  ok = () => {
    this.validateAll()

    if (this.oldPasswdMsg || this.newPasswdMsg || this.passwdAgainMsg) {
      return
    }

    currentUser
      .updatePwd(this.props.username, this.oldPasswd, this.newPasswd)
      .then(res => {
        if (res.data.success) {
          message.success('密码修改成功')
          this.props.onOk()
        }
      })
  }

  render() {
    const { tips } = this.props
    return (
      <>
        <div style={{ textAlign: 'center', padding: '0 0 20px 0' }}>
          {tips && tips}
        </div>
        <Descriptions column={1} colon={false}>
          <Descriptions.Item label={<Label required>原密码:</Label>}>
            <Input
              value={this.oldPasswd}
              style={{ width: 300 }}
              type='password'
              autoFocus
              placeholder='请输入'
              onChange={e => this.onChange('oldPasswd', e.target.value)}
              onBlur={e => this.validate('oldPasswd', e.target.value)}
            />
            <Msg>{this.oldPasswdMsg && this.oldPasswdMsg}</Msg>
          </Descriptions.Item>
          <Descriptions.Item label={<Label required>新密码:</Label>}>
            <Input
              value={this.newPasswd}
              style={{ width: 300 }}
              placeholder='请输入'
              type='password'
              onChange={e => this.onChange('newPasswd', e.target.value)}
              onBlur={e => this.validate('newPasswd', e.target.value)}
            />
            <Msg>{this.newPasswdMsg && this.newPasswdMsg}</Msg>
          </Descriptions.Item>
          <Descriptions.Item label={<Label>{''}</Label>}>
            <PasswordStrengthChecker
              strength={this.strength}
              tips={this.tips}
              style={{ width: 300 }}
            />
          </Descriptions.Item>
          <Descriptions.Item label={<Label required>确认密码:</Label>}>
            <Input
              value={this.passwdAgain}
              style={{ width: 300 }}
              placeholder='请再次输入新密码'
              type='password'
              onChange={e => this.onChange('passwdAgain', e.target.value)}
              onBlur={e => this.validate('passwdAgain', e.target.value)}
            />
            <Msg>{this.passwdAgainMsg && this.passwdAgainMsg}</Msg>
          </Descriptions.Item>
        </Descriptions>
        <Modal.Footer onCancel={this.props.onCancel} onOk={this.ok} />
      </>
    )
  }
}

export default PasswdForm
