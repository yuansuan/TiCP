/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observer, useLocalStore } from 'mobx-react-lite'
import React, { useEffect, useRef } from 'react'
import { message, Switch } from 'antd'
import { jobServer } from '@/server'

type Props = {
  id: string
  name: string
  paused: boolean
}

export const Switcher = observer(function Switcher({
  id,
  name,
  paused,
}: Props) {
  const lockHandler = useRef(undefined)
  const state = useLocalStore(() => ({
    locked: true,
    setLocked(bool) {
      this.locked = bool
    },
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    localChecked: !paused,
    setLocalChecked(bool) {
      if (this.locked) {
        this.localChecked = bool
      }
    },
  }))

  useEffect(() => {
    state.setLocalChecked(!paused)
  }, [paused])

  useEffect(() => {
    return () => clearTimeout(lockHandler.current)
  }, [])

  return (
    <Switch
      checkedChildren='开启'
      unCheckedChildren='暂停'
      loading={state.loading}
      checked={state.localChecked}
      onChange={async checked => {
        state.setLoading(true)
        state.setLocked(false)

        try {
          await jobServer.pauseOrResume(id, checked)
          state.localChecked = checked
          message.success(`作业 ${name} ${checked ? '开启' : '暂停'}回传`)
          // hack: 因为回传启停是异步的，加一个5秒的锁可以防止自动刷新重置开关状态
          lockHandler.current = setTimeout(() => state.setLocked(true), 5000)
        } finally {
          state.setLoading(false)
        }
      }}
    />
  )
})
