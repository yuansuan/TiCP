import { observer } from 'mobx-react'
import * as React from 'react'

import { TabSuite } from '../components'
import { GroupList, RoleList, UserList } from '@/domain/UserMG'
import Operators from './Operators'
import Panel from './Panel'
import { defaultListQuery } from '../utils'

interface IProps {}

@observer
export default class UserManagement extends React.Component<IProps> {
  public render() {
    return (
      <TabSuite
        operators={<Operators />}
        freshen={() =>
          Promise.all([RoleList.fetch(), GroupList.fetch(), UserList.fetch()])
        }
        defaultListQuery={defaultListQuery}
        Panel={Panel}
        store={GroupList}
        searchPlaceholder='按用户组名搜索'
      />
    )
  }
}
