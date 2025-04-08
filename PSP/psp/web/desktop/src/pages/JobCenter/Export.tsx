/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { Button } from '@/components'
import { useStore } from './store'
import { jobCenterServer } from '@/server'

export const Export = observer(function Export() {
  const store = useStore()

  async function exportExcel() {
    const url = await jobCenterServer.export(store.query)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `jobs_${new Date().toISOString()}.xlsx`)
    document.body.appendChild(link)
    link.click()
  }

  return (
    <Button style={{ padding: 0 }} type='link' onClick={exportExcel}>
      导出
    </Button>
  )
})
