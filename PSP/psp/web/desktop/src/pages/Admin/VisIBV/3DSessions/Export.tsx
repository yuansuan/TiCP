/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import qs from 'query-string'
import { Button } from '@/components'
import { useStore } from '../store'
import { currentUser } from '@/domain'

interface IProps {
  disabled?: boolean
}

export const Export = observer(function Export({disabled}: IProps) {
  const store = useStore()

  async function exportExcel() {
    const a = document.createElement('a')
    a.href = `/api/vis-ibv/session/excel?${qs.stringify(
      {
        // user_id: store?.user_id,
        // company_id: store?.company_id,
        status: store?.status,
        hardware_id: store?.hardware_id,
        software_id: store?.software_id,
        project_id: store?.project_id,
        is_admin: currentUser.hasSysMgrPerm,
      },
      {
        arrayFormat: 'bracket',
      }
    )}`
    a.style.display = 'none'
    document.body.appendChild(a)
    a.click()
  }

  return (
    <Button type={'link'} onClick={exportExcel} disabled={disabled}>
      导出
    </Button>
  )
})
