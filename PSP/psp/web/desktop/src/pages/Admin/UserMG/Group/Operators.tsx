import * as React from 'react'

import { Group } from '@/domain/UserMG'
import { Button, Modal } from '@/components'
import GroupEditor from './GroupEditor'

export default class Toolbar extends React.Component {
  public render() {
    return (
      <Button type='primary' icon='add' ghost onClick={this.addGroup}>
        新建用户组
      </Button>
    )
  }

  private addGroup = () => {
    Modal.show({
      title: '新建用户组',
      bodyStyle: { height: 710, background: '#F0F5FD' },
      width: 1130,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <GroupEditor onCancel={onCancel} onOk={onOk} group={new Group()} />
      ),
    }).catch(() => {})
  }
}
