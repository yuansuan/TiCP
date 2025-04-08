import React, { useState } from 'react'
import { observer } from 'mobx-react'

import { Wrapper } from './style'
import { List } from './List'
import { normalLogList, adminLogList, securityLogList } from '@/domain/AuditLog'
import { Tabs } from 'antd'
import { currentUser } from '@/domain'

export default observer(function AuditLogPage() {
  const logTabs = [
    {
      label: '普通用户日志',
      key: 'normal',
      model: normalLogList,
      has: currentUser.hasNormalLogPerm,
    },
    {
      label: '系统管理员日志',
      key: 'admin',
      model: adminLogList,
      has: currentUser.hasSysAdminLogPerm,
    },
    {
      label: '安全管理员日志',
      key: 'security',
      model: securityLogList,
      has: currentUser.hasSecurityAdminLogPerm,
    }
  ].filter(t => t.has)

  const [activeKey, setActiveKey] = useState(logTabs[0].key)

  return (
    <Wrapper>
      <div className='body'>
        <Tabs activeKey={activeKey} onChange={activeKey => setActiveKey(activeKey)} >
          {logTabs.map(tab => {
            return (
              <Tabs.TabPane tab={tab.label} key={tab.key}>
                { activeKey === tab.key && <List logList={tab.model} label={tab.label}/>}
              </Tabs.TabPane>
            )
          })}
        </Tabs>
      </div>
    </Wrapper>
  )
})
