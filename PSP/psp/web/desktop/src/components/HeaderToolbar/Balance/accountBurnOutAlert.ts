/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { account, env } from '@/domain'
import { notification } from '@/components'
import { single } from '@/utils'

const temp = async () => {
  await account.fetch(env.accountId)

  if (!env.isPersonal) {
    // 企业用户
    if (account.account_balance > 0) {
      if (account.account_balance / 100000 < 1000) {
        await notification.info({
          message: '请注意！',
          description:
            '账户余额即将用完，请联系管理员尽快充值，以免影响正常使用流程。'
        })
        return
      }
    } else if (account.account_balance <= 0) {
      if (account.account_balance > account.credit_quota * -1) {
        /*
        await notification.warning({
          message: '请注意！',
          description:
            '账户余额已用完，正在使用授信额度，请尽快联系管理员充值，以免影响正常使用流程。'
        })
        */
      } else {
        await notification.warning({
          message: '请注意！',
          description:
            '账户余额和授信额度均已用完，请尽快联系管理员充值，以免影响正常使用流程。'
        })
      }
      return
    }
  } else {
    // 个人用户
    if (account.account_balance > 0) {
      if (account.account_balance / 100000 < 100) {
        await notification.info({
          message: '请注意！',
          description: '账户余额即将用完，请尽快充值，以免影响正常使用流程。'
        })
        return
      }
    } else if (account.account_balance < 0) {
      await notification.warning({
        message: '请注意！',
        description: '账户余额已用完，请尽快充值。'
      })
      return
    }
  }
}

export const accountBurnOutAlert = () => {
  single('balanceAlert', temp)
}
