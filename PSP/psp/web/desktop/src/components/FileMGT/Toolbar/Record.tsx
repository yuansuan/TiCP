import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { history } from '@/utils'

export const Record = observer(function Record() {
  const record = () => {
    history.push('/file/operation-record', history.location.state)
  }

  return <Button onClick={record}>文件操作记录</Button>
})
