/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Suite } from '@/components/FileMGT'
import { Page } from '@/components'
export default function FileMGT({ job }: any) {
  return (
    <Page header={null}>
      <div style={{ height: 'calc(100vh - 124px)' }}>
        <Suite job={job} />
      </div>
    </Page>
  )
}
