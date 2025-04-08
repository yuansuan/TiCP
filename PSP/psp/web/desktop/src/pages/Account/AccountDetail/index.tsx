/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import { Empty } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { InfoItem, InfoBlock } from './style'
import { currentUser, env } from '@/domain'
import { account } from '@/domain'
import { ComboList } from '../domain/ComboList'
import { CHARGE_TYPE } from '@/constant'

const AccountDetail = observer(function AccountDetail() {
  const store = useLocalStore(() => ({
    combo: new ComboList(),
    get comboData() {
      return this.combo.list
    },
    get balance() {
      return (account.account_balance_contain_freezed / 100000).toFixed(2)
    },
    get amount() {
      return (account.freezed_amount / 100000).toFixed(2)
    },
    get quota() {
      return (account.credit_quota / 100000).toFixed(2)
    }
  }))

  useEffect(() => {
    account.fetch(env.accountId)
  }, [])

  useEffect(() => {
    if (!env.isPersonal) {
      store.combo.fetch(env.productId)
    }
  }, [])

  return (
    <InfoBlock>
      <div className='header'>
        <div className='title'>基本信息</div>
      </div>
      <div className='account'>
        {currentUser.id ? (
          <>
            <div className='info'>
              <InfoItem title={store.balance}>
                账户余额：{store.balance} 元
              </InfoItem>
              <InfoItem title={store.amount}>
                未结算金额：{store.amount} 元
              </InfoItem>
              {!env.isPersonal && (
                <InfoItem title={store.quota}>
                  授信额度：{store.quota} 元
                </InfoItem>
              )}
            </div>

            <div className='list'>
              {store.combo.list?.map(list => (
                <div className='combo'>
                  <h3>{list.combo_name}</h3>

                  {list.chargeType === CHARGE_TYPE.MONTHLY_TYPE ? (
                    <>
                      <InfoItem title={list.total_time + '个月'}>
                        包年包月：
                        {Math.round(list.total_time / (3600 * 24 * 30))} 个月
                      </InfoItem>
                      <InfoItem
                        title={
                          list.valid_begin_time.toString() +
                          '至' +
                          list.valid_end_time.toString()
                        }>
                        有效期：{list.valid_begin_time.toString()} 至<br />{' '}
                        {list.valid_end_time.toString()}
                      </InfoItem>
                    </>
                  ) : null}
                  {list.chargeType === CHARGE_TYPE.HOURLY_TYPE ? (
                    <>
                      <InfoItem title={list.total_time + '小时'}>
                        包小时：{Math.round(list.total_time / 3600)} 小时
                      </InfoItem>
                      <InfoItem title={list.used_time + '小时'}>
                        已使用：{(list.used_time / 3600).toFixed(2)} 小时
                      </InfoItem>
                      <InfoItem title={list.used_time + '小时'}>
                        未使用：{(list.remain_time / 3600).toFixed(2)} 小时
                      </InfoItem>
                      <InfoItem>
                        激活时间：{list.valid_begin_time.toString()}
                      </InfoItem>
                    </>
                  ) : null}
                </div>
              ))}
            </div>
          </>
        ) : (
          <Empty description='用户未登录' />
        )}
      </div>
    </InfoBlock>
  )
})

export default AccountDetail
