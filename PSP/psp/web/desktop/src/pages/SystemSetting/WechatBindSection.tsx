/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Modal, Button, Table } from '@/components'
import { Pagination, message } from 'antd'
import { useLocalStore } from 'mobx-react-lite'
import { Http } from '@/utils'
import { Timestamp } from '@/domain/common'
import IEditableText from '@/pages/SystemSetting/IEditableText'
import { showQRCodeModal } from '@/components'
import { runInAction } from 'mobx'
import { userServer } from '@/server'

const StyledDiv = styled.div`
  padding: 0 20px 20px 20px;
  font-family: PingFangSC-Regular;

  > h1 {
    margin-bottom: 16px;
    font-size: 16px;
    font-family: 'PingFangSC-Medium';
    color: #333333;

    span {
      font-size: 12px;
      color: #9b9b9b;
    }
  }

  > .desc-row-1 {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;

    font-size: 12px;
    color: #333333;
    margin-bottom: 16px;
  }

  > .desc-row-2 {
    font-size: 12px;
    color: #999999;
    letter-spacing: 0;
    margin-bottom: 16px;
  }

  > .bind-new {
    margin-top: 10px;
    margin-bottom: 16px;
  }

  .ant-pagination.pagination {
    padding: 20px;
    text-align: right;
  }
`

export const WechatBindSection = observer(function WechatBindSection() {
  const store = useLocalStore(() => ({
    loading: false,
    setLoading(bool) {
      this.loading = bool
    },
    get dataSource() {
      return this.data.map(v => ({
        ...v,
        activate_time: new Timestamp(v.activate_time).toString(),
      }))
    },
    data: [],
    setData(list) {
      this.data = list
    },
    pageCtx: { index: 1, size: 10, total: 0 },
    setPageCtx(v) {
      this.pageCtx = v
    },
    alarmBalance: undefined,
    setAlarmBalance(v) {
      this.alarmBalance = v
    },
    averageCost: undefined,
    setAverageCost(v) {
      this.averageCost = v
    },
    bound: undefined,
    setBound(bool) {
      this.bound = bool
    },
  }))

  useEffect(() => {
    fetchList()
  }, [store.pageCtx?.index, store.pageCtx?.size])

  useEffect(() => {
    getCostAverageIn3Days()
    fetchAlarmBalance()
  }, [])

  async function fetchList() {
    try {
      store.setLoading(true)
      const {
        data: { list, page_ctx },
      } = await Http.get('/company/wx/alarm_balance/list', {
        params: {
          page_index: store.pageCtx.index,
          page_size: store.pageCtx.size,
        },
      })

      runInAction(() => {
        store.setData(list)
        store.setPageCtx(page_ctx)
      })
    } finally {
      store.setLoading(false)
    }
  }

  async function getCostAverageIn3Days() {
    const {
      data: { average },
    } = await Http.get('/company/cost/average/3')

    runInAction(() => {
      store.setAverageCost(average / 100000)
    })
  }

  async function fetchAlarmBalance() {
    const {
      data: { alarmBalance },
    } = await Http.get('/company/config/list', {
      params: {
        keys: ['alarmBalance'],
      },
    })

    runInAction(() => {
      store.setAlarmBalance(alarmBalance)
    })
  }

  async function putAlarmBalance(alarmBalance) {
    await Http.put(
      '/company/config',
      {
        key: 'alarmBalance',
        value: alarmBalance,
      },
      {
        formatErrorMessage: () => '请输入大于 0 的整数',
      }
    )
  }

  async function unbindOne(wechat_openid, user_id) {
    await Http.delete(`/company/wx/alarm_balance/${wechat_openid}`, {
      params: { user_id },
    })
  }

  function onPageChange(index, size) {
    store.setPageCtx({
      ...store.pageCtx,
      index,
      size,
    })
  }

  return (
    <StyledDiv className='wx-bind-setting'>
      <h1>
        绑定微信<span>（绑定后可收到余额变动提醒）</span>
      </h1>
      <div className='desc-row-1'>
        可用余额小于该值时提醒：
        <IEditableText
          style={{
            width: 160,
            display: 'inline-block',
            height: 24,
            minHeight: 24,
          }}
          customHelp='请输入大于0的整数'
          customBeforeConfirm={(value: string) => {
            if (value === '') return '请输入大于 0 的整数'
            if (+(+value).toFixed(0) === 0) return '请输入大于 0 的整数'

            const flag = /^[+]{0,1}(\d+)$|^[+]{0,1}$/.test(value) // match自然数
            if (!flag) {
              return '请输入大于 0 的整数'
            } else {
              return true
            }
          }}
          unit='元'
          value={store.alarmBalance}
          setValue={async value => {
            store.setAlarmBalance(value)
            await putAlarmBalance(value).finally(() => {
              fetchAlarmBalance()
            })
            message.success('设置成功')
          }}
        />
      </div>
      <div className='desc-row-2'>
        提示：近三日您的日均消耗金额为 {store.averageCost?.toFixed(2)} 元
      </div>
      <Button
        type='primary'
        className='bind-new'
        onClick={() =>
          showQRCodeModal({
            descriptionNode: (
              <>
                <span>
                  使用微信扫描以下二维码，关注
                  <span className='bold'>“云仿真平台”</span>
                  公众号。
                </span>
                <span>
                  关注公众号后，您能够及时收到余额不足的通知，方便您更高效地使用云仿真平台。
                </span>
              </>
            ),
            validConfig: {},
            fetchQRCodeFunc: () => userServer.getWxCode('balance'),
            hideOk: true,
            afterCancel: fetchList,
          })
        }>
        新增绑定
      </Button>
      <Table
        props={{
          data: store.dataSource,
          rowKey: 'id',
          autoHeight: true,
          loading: store.loading,
        }}
        columns={[
          {
            header: '微信昵称',
            dataKey: 'wechat_nickname',
            props: {
              flexGrow: 1,
            },
          },
          {
            header: '绑定时间',
            dataKey: 'activate_time',
            props: {
              flexGrow: 1,
            },
          },
          {
            header: '操作',
            props: {
              width: 100,
            },
            cell: {
              props: {
                dataKey: 'wechat_openid',
              },
              render: ({ rowData, dataKey }) => (
                <Button
                  type='link'
                  onClick={async () => {
                    await Modal.showConfirm({
                      title: '提示',
                      content: (
                        <>
                          解绑后该微信号将无法及时收到余额变动等通知，确定要解除绑定
                          {rowData['wechat_nickname']} 吗?
                        </>
                      ),
                    })

                    await unbindOne(rowData[dataKey], rowData['user_id'])
                      .then(() => {
                        message.success('解绑成功')
                      })
                      .catch(() => message.error('解绑失败'))
                      .finally(() => fetchList())
                  }}>
                  解绑
                </Button>
              ),
            },
          },
        ]}
      />

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={store.pageCtx?.size}
        current={store.pageCtx?.index}
        total={store.pageCtx?.total}
        onChange={onPageChange}
      />
    </StyledDiv>
  )
})
