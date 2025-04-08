import React, { useEffect } from 'react'
import { Page } from '@/components/Page'
import styled from 'styled-components'
import { useParams } from 'react-router'
import { Pagination, Table } from 'antd'
import { lmUsageList } from '@/domain'
import { reaction } from 'mobx'
import { observer } from 'mobx-react-lite'

export const StyledLayout = styled.div`
  padding: 30px;
  overflow-y: hidden;

  .pagination {
    text-align: right;
    margin: 20px 20px 0;
  }
`

const ListWrapper = styled.div`
  padding-top: 20px;
`

export default observer(function LicenseUsage() {
  const { id } = useParams<{ id: string }>()

  function refresh() {
    lmUsageList.fetch(id)
  }

  useEffect(() => {
    lmUsageList?.onPageChange(1, 10)
    refresh()
  }, [])

  useEffect(() => {
    let disposer = reaction(
      () => ({
        index: lmUsageList.index,
        size: lmUsageList.size
      }),
      () => {
        lmUsageList.fetch(id)
      }
    )

    return () => {
      disposer()
    }
  }, [])

  const columns = [
    { title: '企业名称', dataIndex: 'company_name', key: 'company_name' },
    { title: '使用数', dataIndex: 'licenses', key: 'licenses' },
    { title: '作业ID', dataIndex: 'job_id', key: 'job_id' },
    { title: '作业名称', dataIndex: 'job_name', key: 'job_name' },
    {
      title: '应用名称',
      dataIndex: 'app_name',
      key: 'app_name'
    },
    {
      title: '开始时间',
      dataIndex: 'create_time',
      key: 'create_time',
      render: text => text.dateString
    }
  ]

  return (
    <Page header={null}>
      <StyledLayout>
        {/* 过滤条件后端没加，这一期暂时隐藏
        <Toolbar refresh={refresh} /> */}
        <ListWrapper>
          <Table
            rowKey={'id'}
            columns={columns}
            dataSource={lmUsageList?.list}
            pagination={false}
          />
          {lmUsageList?.total > 0 && (
            <Pagination
              className='pagination'
              showQuickJumper
              showSizeChanger
              pageSize={lmUsageList.size}
              current={lmUsageList.index}
              total={lmUsageList?.total}
              onShowSizeChange={lmUsageList.onSizeChange.bind(lmUsageList)}
              onChange={lmUsageList?.onPageChange.bind(lmUsageList)}
            />
          )}
        </ListWrapper>
      </StyledLayout>
    </Page>
  )
})
