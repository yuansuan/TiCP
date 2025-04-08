import * as React from 'react'

import { Button, Modal } from '@/components'
import { User } from '@/domain/UserMG'
import currentUser from '@/domain/User'
import UserCreatorByFile from './UserCreatorByFile'
import UserCreater from './LDAPUserCreator'

export default class Toolbar extends React.Component<any> {
  public render() {
    return (
      <>
        <Button type='primary' icon='add' ghost onClick={this.addUser}>
          新建用户
        </Button>
      </>
    )
  }

  private addUser = () => {
    Modal.show({
      title: '新建用户',
      bodyStyle: { height: 710, background: '#F0F5FD' },
      width: 900,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <UserCreater onCancel={onCancel} onOk={onOk} />
      )
    }).catch(() => {})
  }
  private addUserByFile = () => {
    Modal.show({
      title: '批量新建用户',
      bodyStyle: { height: 710, background: '#F0F5FD' },
      width: 900,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <UserCreatorByFile onOk={onOk} user={new User()} />
      )
    }).catch(() => {})
  }
}
