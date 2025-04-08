import React, { useState, useCallback } from 'react'
import { Popover } from 'antd'
import { Http } from '@/utils'
import { ApproveTimeline } from '@/components/ApproveTimeline'

export default function ApproveStatus({
  targetId,
  targetType,
  title = '',
  data,
  callback,
  children,
}) {
  const [loading, setLoading] = useState(false)
  const [approve, setApprove] = useState({})

  const { approve_status } = data

  let content = (
    <div>
      {loading ? <span>加载中...</span> : <ApproveTimeline approve={approve} />}
    </div>
  )

  const getApproveInfo = useCallback(() => {
    ;(async () => {
      setLoading(true)
      try {
        const res = await Http.get(
          `/audit/ask/latest?targetId=${targetId}&type=${targetType}`
        )
        if (approve_status !== res.data.status) {
          callback && callback()
        }
        setApprove(res.data)
      } finally {
        setLoading(false)
      }
    })()
  }, [approve_status, targetId])

  const onChange = val => {
    if (val) getApproveInfo()
  }

  return (
    <>
      {approve_status === -1 ? (
        <>{children}</>
      ) : (
        <Popover
          placement='left'
          title={title || '最近一次申请信息'}
          content={content}
          trigger='click'
          onVisibleChange={onChange}>
          {children}
        </Popover>
      )}
    </>
  )
}
