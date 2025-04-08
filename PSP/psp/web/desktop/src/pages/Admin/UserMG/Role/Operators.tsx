import * as React from 'react'

import { Role } from '@/domain/UserMG'
import { Button, Modal } from '@/components'
import RoleEditor from './RoleEditor'

export default class Toolbar extends React.Component {
  public render() {
    return (
      <Button type='primary' icon='add' ghost onClick={this.addRole}>
        新建角色
      </Button>
    )
  }
  private addRole = () => {
    Modal.show({
      title: '新建角色',
      bodyStyle: { height: 710, background: '#F0F5FD', overflow: 'auto' },
      width: 1130,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <RoleEditor
          role={new Role()}
          isAdd={true}
          onCancel={onCancel}
          onOk={onOk}
        />
      )
    }).catch(() => {})
  }
}
