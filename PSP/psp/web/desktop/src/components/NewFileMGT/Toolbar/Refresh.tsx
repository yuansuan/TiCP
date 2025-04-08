/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { useStore } from '../store'

export const Refresh = observer(function Refresh() {
  const store = useStore()
  const [refresh] = store.useRefresh()

  return <Button onClick={refresh}>刷新</Button>
})
