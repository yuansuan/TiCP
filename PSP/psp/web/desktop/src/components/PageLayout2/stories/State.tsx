/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import PageLayout from '..'

const { useStore } = PageLayout

function StatePreviwer() {
  const {
    menuExpanded: [menuExpanded],
  } = useStore()

  return <div>菜单状态：{menuExpanded ? '展开' : '收缩'}</div>
}

export const State = () => (
  <PageLayout>
    <StatePreviwer />
  </PageLayout>
)
