import { observer } from 'mobx-react'
import * as React from 'react'

import { TabSuite } from '@/pages/Admin/UserMG/components'
import { defaultListQuery } from '@/pages/Admin/UserMG/utils'
import { RoleList } from '@/domain/UserMG'
import Operators from './Operators'
import Panel from './Panel'

interface IProps {}

@observer
export default class RoleManagement extends React.Component<IProps> {
  public render() {
    return (
      <TabSuite
        operators={<Operators />}
        freshen={RoleList.fetch}
        defaultListQuery={defaultListQuery}
        Panel={Panel}
        store={RoleList}
        searchPlaceholder='按角色名搜索'
      />
    )
  }
}
