/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useMemo, useState } from 'react'
import styled from 'styled-components'
import { Card, Empty, Pagination } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import debounce from 'lodash/debounce'

import { lastMessages, Messages as DomainMessages } from '@/domain'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;

  .pagination {
    margin: 20px auto;
  }

  .body {
    .card {
      margin: 10px;
      .cardBody {
        display: flex;
        .right {
          margin-left: auto;

          > button {
            margin: 0 4px;
          }

          .isRead {
            color: #ccc;
          }

          .notRead {
            cursor: pointer;
            color: #ff4d4f;
          }
        }
      }
    }
  }
`

interface IProps {
  searchKey: string
  visible: boolean
}

export const Messages = observer(function Messages(props: IProps) {
  const [page, setPage] = useState({ pageIndex: 1, pageSize: 10 })
  const store = useLocalStore(() => ({
    messages: new DomainMessages()
  }))

  const { list, page_ctx } = store.messages
  const { pageIndex, pageSize } = page

  const params = useMemo(
    () => ({
      page_index: page.pageIndex,
      page_size: page.pageSize,
      filter: { content: props.searchKey, state: 0 }
    }),
    [page.pageIndex, page.pageSize, props.searchKey]
  )

  const fetch = () => store.messages.fetch(params)


  useEffect(() => {
    const debouncedFetch = debounce(function () {
     fetch()
    }, 300)

   if(props.visible) debouncedFetch()
  }, [params])

  const onPageChange = (pageIndex, pageSize) => {
    setPage({
      ...page,
      pageIndex,
      pageSize
    })
  }

  const read = async item => {
    await item.read()
    lastMessages.fetchUnreadCount()
  }

  function description(item) {
    switch (item.keyType) {
      case 'Message_ComputeJobEvent': {
        return (
          // <Link to={`/job/${item.body.job_id}`}>
          <div>
            <span style={{ paddingRight: 4 }}>{item.body.job_name}</span>
            {item.message}
          </div>
          // </Link>
        )
      }

      default:
        return item.message
    }
  }

  return (
    <StyledLayout>
      <div className='body'>
        {list.length === 0 && <Empty description='暂无消息通知' />}

        {list.map(item => (
          <Card
            className='card'
            key={item.id}
            title={item.typeToString}
            extra={item.timeTitle}>
            <div className='cardBody'>
              <div className='content'>{description(item)}</div>
              <div className='right'>
                {item.state === 2 ? (
                  <span className='isRead'>已读</span>
                ) : (
                  <span className='notRead' onClick={() => read(item)}>
                    标为已读
                  </span>
                )}
              </div>
            </div>
          </Card>
        ))}
      </div>

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={pageSize}
        current={pageIndex}
        total={page_ctx.total}
        onChange={onPageChange}
      />
    </StyledLayout>
  )
})
