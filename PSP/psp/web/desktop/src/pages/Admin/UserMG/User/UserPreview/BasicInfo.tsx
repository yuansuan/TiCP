import * as React from 'react'
import { BasicInfoWrapper } from './style'
import { Modal } from '@/components'
import { User } from '@/domain/UserMG'
import { currentUser } from '@/domain'
import { message } from 'antd'
import { action, observable } from 'mobx'
import { observer } from 'mobx-react'

interface IProps {
  user: User
}

const textFlow = textStr => {
  if (!textStr) return ''
  let res = textStr
  if (textStr.length > 120) {
    res = `${textStr.slice(0, 60)}\n${textStr.slice(60, 120)}\n${textStr.slice(
      120
    )}`
  } else if (textStr.length > 60) {
    res = `${textStr.slice(0, 60)}\n${textStr.slice(60)}`
  } else {
    res = textStr
  }
  return res
}

@observer
export default class BaseInfo extends React.Component<IProps> {
  
  private reGen = () => {
    if (!this.props.user.enable_openapi) {
      return
    }
    
    Modal.showConfirm({
      content:"确定重新生成openapi凭证吗?"
    }).then(() => {
      this.props.user.genCert().then(
        res => { 
          message.success('重新生成成功')
        }
      )
    })
  }
  
  render() {
    
    const { name, email, mobile, enable_openapi } = this.props.user
    // console.info('props.user值为:', JSON.stringify(this.props.user, null, 2));


    const titleEmail = textFlow(email)
    return (
      <div>
        <BasicInfoWrapper>
        <div>
          <label>用户名：</label>
          <span title={name}>{name}</span>
        </div>

        <div>
          <label>邮箱：</label>
          <span title={titleEmail}>{email || '--'}</span>
        </div>

        <div>
          <label>电话：</label>
          <span title={mobile}>{mobile || '--'}</span>
        </div>
        
        </BasicInfoWrapper>

        <BasicInfoWrapper>
        {currentUser?.isOpenapiSwitchEnable && enable_openapi && (
          <div>
            <label>openapi凭证：</label>
            <span title={this.props.user.openapi_certificate}>{this.props.user.openapi_certificate || '--'}</span>
            <span className='text cert' onClick={this.reGen}>
              重新生成
            </span>
          </div>
        )}
        </BasicInfoWrapper>

      </div>
    )
  }
}
