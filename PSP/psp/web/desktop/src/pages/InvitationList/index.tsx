/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useMemo, useEffect } from 'react'
import styled from 'styled-components'
import { Pagination, Empty } from 'antd'
import { InvitationList as DomainList, env } from '@/domain'
import { Invitation } from './Invitation'
import { Page } from '@/components'
import { useObserver } from 'mobx-react-lite'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;

  .pagination {
    margin: 20px 20px 20px auto;
  }
`

export default function InvitationList() {
  const invitations = useMemo(() => new DomainList(), [])
  const [pageIndex, setPageIndex] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  useEffect(() => {
    invitations.fetchByCompany({
      company_id: env.company?.id,
      page_index: pageIndex,
      page_size: pageSize,
    })
  }, [pageIndex, pageSize])

  function onPageChange(index, size) {
    setPageIndex(index)
    setPageSize(size)
  }

  return useObserver(() => {
    const {
      list,
      pageCtx: { total },
    } = invitations

    return (
      <Page header={null}>
        <StyledLayout>
          {total === 0 && (
            <Empty style={{ padding: 20 }} description='暂无成员邀请' />
          )}
          {list.map(item => (
            <Invitation key={item.id} item={item} />
          ))}

          <Pagination
            className='pagination'
            showQuickJumper
            showSizeChanger
            pageSize={pageSize}
            current={pageIndex}
            total={total}
            onChange={onPageChange}
          />
        </StyledLayout>
      </Page>
    )
  })
}
