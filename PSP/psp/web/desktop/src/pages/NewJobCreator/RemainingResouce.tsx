/* Copyright (C) 2016-present, Yuansuan.cn */

import React, { useRef, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useStore } from './store'
import { Table } from 'antd'
import { scList } from '@/domain'

export const RemainingResouce = observer(function RemainingResouce() {
  const store = useStore()
  const timeoutID = useRef(null)
  useEffect(() => {
    (async () => {
      await store?.data?.currentApp?.getRemainingResource(
        [store.scId]
      )
    })()
  }, [store.scId])
  useEffect(() => {
    if (!store.data.currentApp) return
    timeoutID.current && clearTimeout(timeoutID.current)
    ;(async function getRemainingResource() {
      await store.data.currentApp.getRemainingResource(
        [store.scId]
      )
      timeoutID.current = setTimeout(getRemainingResource, 10000)
    })()
    return () => {
      clearTimeout(timeoutID.current)
      timeoutID.current = null
    }
  }, [store.data.currentApp])

  return (
    store.data.currentApp?.remainingResource && (
      <div className='remaining-resource'>
        <Table
          dataSource={store.data.currentApp?.remainingResource}
          columns={[
            {
              title: '资源名称',
              dataIndex: 'sc_id',
              key: 'sc_id',
              render: sc_id => <div>{scList.getName(sc_id)}</div>
            },
            {
              title: '空闲资源',
              dataIndex: 'cores',
              key: 'cores',
              render: cores => cores && <div>{cores} 核</div>
            }
          ]}
          pagination={false}
        />
      </div>
    )
  )
})
