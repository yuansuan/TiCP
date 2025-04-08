/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { DatePicker, Pagination } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Table } from '@/components'
import { StyledLayout } from './style'
import moment from 'moment'
import { DetailList } from '../domain/AccountList/index'
import { env } from '@/domain'

const { RangePicker } = DatePicker

const AccountList = observer(function AccountList() {
  const state = useLocalStore(() => ({
    account: new DetailList(),
    start: moment().subtract(10, 'days').startOf('day').unix(),
    end: moment().endOf('day').unix(),
    setDate(date) {
      if (date !== null) {
        this.start = date[0].startOf('day').unix()
        this.end = date[1].endOf('day').unix()
      }
    },
    get date() {
      return {
        start: this.start,
        end: this.end,
      }
    },
    pageIndex: 1,
    setPageIndex(index) {
      this.pageIndex = index
    },
    pageSize: 10,
    setPageSize(size) {
      this.pageSize = size
    },
    get dataSource() {
      let data = [...(this.account && this.account.list)].map(item => ({
        ...item,
        trade_time: item.trade_time.toString(),
        trade_type: {
          1: '支付',
          2: '充值',
          3: '退款',
          4: '提现',
          5: '加款',
          6: '扣款',
        }[item.trade_type],
        bill_sign: {
          1: '收入',
          2: '支出',
        }[item.bill_sign],
        amount: (item.amount / 100000).toFixed(2),
        account_balance_contain_freezed: (
          item.account_balance_contain_freezed / 100000
        ).toFixed(2),
      }))
      return data
    },
  }))

  function onPageChange(index, size) {
    state.setPageIndex(index)
    state.setPageSize(size)
  }

  useEffect(() => {
    state.account.fetch({
      start_seconds: state.date.start,
      end_seconds: state.date.end,
      page_index: state.pageIndex,
      page_size: state.pageSize,
      account_id: env.accountId,
    })
  }, [state.date, state.pageIndex, state.pageSize])

  return (
    <StyledLayout>
      <div className='header'>
        <div className='title'>收支明细</div>
      </div>
      <div className='body'>
        <div className='date-picker'>
          <label>日期范围选择：</label>
          <RangePicker
            onChange={date => {
              state.setDate(date)
              state.setPageIndex(1)
            }}
            style={{ width: 400 }}
            defaultValue={[moment.unix(state.start), moment.unix(state.end)]}
          />
        </div>
        <div className='detail'>
          <Table
            props={{
              height: 500,
              data: state.dataSource,
            }}
            columns={[
              {
                header: '账单编号',
                dataKey: 'trade_id',
                props: {
                  colSpan: 2,
                  flexGrow: 1,
                },
              },
              {
                header: '消费时间',
                dataKey: 'trade_time',
                props: {
                  flexGrow: 1.3,
                },
              },
              {
                header: '收支类型',
                dataKey: 'bill_sign',
                props: {
                  flexGrow: 1,
                },
              },
              {
                header: '交易类型',
                dataKey: 'trade_type',
                props: {
                  flexGrow: 1,
                },
              },
              {
                header: '订单号',
                dataKey: 'out_trade_id',
                props: {
                  flexGrow: 1,
                },
              },
              {
                header: '交易备注',
                dataKey: 'remark',
                props: {
                  flexGrow: 1,
                },
              },
              {
                header: '收支金额(元)',
                dataKey: 'amount',
                props: {
                  flexGrow: 1,
                },
              },
              {
                header: '账户余额(元)',
                dataKey: 'account_balance_contain_freezed',
                props: {
                  flexGrow: 1,
                },
              },
            ]}
          />
        </div>
        <div className='Pagination'>
          <Pagination
            total={state.account.page_ctx.total}
            onChange={onPageChange}
            current={state.pageIndex}
            showSizeChanger
          />
        </div>
      </div>
    </StyledLayout>
  )
})

export default AccountList
