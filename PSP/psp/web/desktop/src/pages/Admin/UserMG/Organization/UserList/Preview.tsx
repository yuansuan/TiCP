import UserOfOrg from '@/domain/UserMG/UserOfOrg'
import { observer } from 'mobx-react'
import React from 'react'
import { PermPreview, RadiusItem, Section } from '../../components'
import BasicInfo from '../../User/UserPreview/BasicInfo'
import { UserEditorWrapper } from '../../User/UserPreview/style'

type Props = {
  user: UserOfOrg
}
export const UserPreview = observer(function UserPreview({ user }: Props) {
  return (
    <UserEditorWrapper>
      <BasicInfo user={user} />
      <Section title='用户角色：'>
        <RadiusItem itemList={user.roleNames || []} />
      </Section>
      <PermPreview perms={user.permList} />
    </UserEditorWrapper>
  )
})
