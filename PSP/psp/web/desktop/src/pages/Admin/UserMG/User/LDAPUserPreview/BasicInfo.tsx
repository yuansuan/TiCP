import * as React from 'react'
import { Descriptions } from 'antd'
import { User } from '@/domain/UserMG'
import { formatDateFromMilliSecWithTimeZone } from '@/utils/formatter'
import { StatsBall } from '@/components'
import { BaseInfoWraper } from './style'

interface IProps {
  user: User
}

export default class BaseInfo extends React.Component<IProps> {
  emptyItem = () => <Descriptions.Item label=''>{}</Descriptions.Item>

  render() {
    const { name, email, mobile, created_at, openapi_certificate, enable_openapi } = this.props.user

    return (
      <BaseInfoWraper>
        <Descriptions title='用户信息'>
          <Descriptions.Item label='登录名称'>{name}</Descriptions.Item>
          <Descriptions.Item label='电话'>{mobile}</Descriptions.Item>
          <Descriptions.Item label='邮件'>{email}</Descriptions.Item>
          {/* {this.emptyItem()} */}
          {enable_openapi && <Descriptions.Item label='openapi凭证'>{openapi_certificate}</Descriptions.Item>}
          <Descriptions.Item label='创建时间' span={3}>
            {formatDateFromMilliSecWithTimeZone(created_at)}
          </Descriptions.Item>
        </Descriptions>
      </BaseInfoWraper>
    )
  }
}
