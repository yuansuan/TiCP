import React, { useState } from 'react'

import { Button } from 'antd'
import { observer } from 'mobx-react'
import organization from '@/domain/UserMG/UserOfOrgList'
import { message } from 'antd'
import { useStore } from './store'

type Props = {}

export const Toolbar = observer(function Toolbar({}: Props) {
  const store = useStore()
  const [fetch, loading] = store.getOrganization()
  const [syncLoading, setSyncLoading] = useState(false)

  async function synchronize() {
    try {
      setSyncLoading(true)
      const res = await organization.syncOrganizationStructure()
      if (res.success) {
        message.success('用户组织结构同步成功')
        fetch()
      } else {
        message.error('用户组织结构同步失败')
      }
    } finally {
      setSyncLoading(false)
    }
  }

  return (
    <Button type='primary' ghost onClick={synchronize} loading={syncLoading}>
      同步
    </Button>
  )
})
