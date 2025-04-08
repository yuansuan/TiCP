/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { useStore } from '../store'

type Props = {
  isSyncToLocal: boolean,
  userName: string,
}
export const Refresh = observer(function Refresh({ isSyncToLocal,userName }: Props) {
  const store = useStore()
  const [refresh] = store.useRefresh(isSyncToLocal,userName)
  return <Button id="JobDetailListRefreshBtn" onClick={refresh}>刷新</Button>
})
