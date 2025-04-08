/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useMemo, useState } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Empty, Pagination } from 'antd'
import { lastInvitations } from '@/domain'
import { InvitationList } from '@/domain/InvitationList'
import { Invitation as DomainInvitation } from '@/domain/InvitationList/Invitation'
import { Share } from './Share'
import { recordList } from '@/domain'

import { useTranslation } from 'react-i18next'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;

  .pagination {
    margin: 20px auto;
  }
`

interface IProps {
  visible: boolean
  searchKey: string
}

export const ShareList = observer(function Invitations(props: IProps) {
  const [page, setPage] = useState({ pageIndex: 1, pageSize: 10 })


  const params = useMemo(
    () => ({
      index: page.pageIndex,
      size: page.pageSize,
    }),
    [page.pageIndex, page.pageSize]
  )

  useEffect(() => {
    if (props.visible) fetch()
  }, [params, props.visible])

  const fetch = () => recordList.fetch(params)

  const onPageChange = (pageIndex, pageSize) => {
    setPage({
      ...page,
      pageIndex,
      pageSize,
    })
  }


  const readShare = async item => {
    await item.readShare(item)
    recordList.fetchLast()
  }
  

  return (
    <StyledLayout>
      <div className='body'>
        {recordList.list.length === 0 && (
          <Empty description='暂无分享内容' />
        )}
        {recordList.list.map(item => (
          <Share
            key={item.id}
            item={item}
            readShare={() => readShare(item)}
          />
        ))}
      </div>

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={page.pageSize}
        current={page.pageIndex}
        total={recordList.pageCtx.total}
        onChange={onPageChange}
      />
    </StyledLayout>
  )
})
