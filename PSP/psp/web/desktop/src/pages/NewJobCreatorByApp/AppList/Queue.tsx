/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect, useCallback } from 'react'
import styled from 'styled-components'
import ReactDOM from 'react-dom'
import { Select } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from '../store'
import { appList } from '@/domain'
import { Http, getComputeType } from '@/utils'

const Wrapper = styled.div`
  width: 300px;
  > .ant-select {
    width: 100%;
  }
`

interface IQueueProps {
  queues: Array<{ queue_name: string; cpu_number: number; select: boolean }>
}

export const QueueList = observer(function QueueList(props: IQueueProps) {
  const store = useStore()
  const ref = useRef(null)

  const { queues = [] } = props

  useEffect(() => {
    if (queues.length) {
      if (queues.length === 1) {
        // store.setQueue(queues[0].queue_name)
      }
    }
  }, [queues])

  const onSelectQueues = async values => {
    store.setJobQueue(values)
  }

  return (
    <Wrapper>
      <Select onChange={onSelectQueues} allowClear={true} placeholder='队列'>
        {queues.map(item => {
          return (
            <Select.Option
              title={`队列 ${item.queue_name} -- 可用 ${item.cpu_number} 核数`}
              key={item.queue_name}
              value={item.queue_name}>
              {`队列 ${item.queue_name} -- 可用 ${item.cpu_number} 核数`}
            </Select.Option>
          )
        })}
      </Select>
    </Wrapper>
  )
})
