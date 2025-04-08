import React from 'react'
import { Input, Radio } from 'antd'
import { BasicInfoWrapper, EditWrapper } from './style'
import { Validator } from '@/utils'
import { currentUser } from '@/domain'

enum MessageType {
  EMAIL = 'email',
  MOBILE = 'mobile',
}

export default function BasicInfo({
  name,
  email,
  mobile,
  enable_openapi,
  updateEmail,
  updateMobile,
  updateEnableOpenapi,
  updateError,
}) {
  const validateEmail = e => {
    const email = e.target.value
    const res = Validator.isValidEmail(email)
    let err = ''
    if (!email.length) {
      err = ''
    } else if (!res) {
      err = '邮箱格式错误'
    } else if (email.length > 64) {
      err = '邮箱的长度不能超过 64 个字符'
    }

    updateError(MessageType.EMAIL, err)
    updateEmail(email)
  }

  const validateMobile = e => {
    const mobile = e.target.value
    const res = Validator.isValidPhoneNumber(mobile)
    let err = ''
    if (!mobile.length) {
      err = ''
    } else if (!res) {
      err = '电话格式错误'
    }

    updateError(MessageType.MOBILE, err)

    updateMobile(mobile)
  }

  const infoMap = [
    {
      title: '邮箱',
      defaultValue: email,
      onChange: validateEmail,
    },
    {
      title: '电话',
      defaultValue: mobile,
      onChange: validateMobile,
    },
  ]
  console.info({enable_openapi})

  return (
    <BasicInfoWrapper>
      <label>
        <span>用户名：</span>
        <span title={name}>{name}</span>
      </label>

      {infoMap.map(i => (
        <EditWrapper key={i.title}>
          <span>{i.title}：</span>
          <label className='mgRight'>
            <Input defaultValue={i.defaultValue} onChange={i.onChange} />
          </label>
        </EditWrapper>
      ))}

      {
        currentUser?.isOpenapiSwitchEnable && (
          <>
            <label>允许调用openapi：</label>
            <Radio.Group
              value={enable_openapi}
              onChange={e => {
                updateEnableOpenapi(e.target.value)
              }}>
              <Radio value={false}>禁用</Radio>
              <Radio value={true}>启用</Radio>
            </Radio.Group>
          </>
        )
      }
    </BasicInfoWrapper>
  )
}
