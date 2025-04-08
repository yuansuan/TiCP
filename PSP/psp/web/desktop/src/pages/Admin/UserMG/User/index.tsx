import * as React from 'react'
import { observer } from 'mobx-react-lite'
import { TabSuite } from '@/pages/Admin/UserMG/components'
import { defaultListQuery } from '@/pages/Admin/UserMG/utils'
import { RoleList, UserList } from '@/domain/UserMG'
import { currentUser } from '@/domain'
import Operators from './Operators'
import Panel from './Panel'

const UserManagement = observer(() => {
  return (
    <TabSuite
      freshen={() => Promise.all([RoleList.fetch(), UserList.fetch()])}
      defaultListQuery={defaultListQuery}
      operators={currentUser?.isLdapEnabled ? <Operators /> : null}
      Panel={Panel}
      store={UserList}
      searchPlaceholder='按用户名搜索'
    />
  )
})

export default UserManagement
