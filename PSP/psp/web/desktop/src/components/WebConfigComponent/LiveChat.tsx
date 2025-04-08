/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { currentUser, env, webConfig } from '@/domain'

export const LiveChatComponent = function LiveChatComponent() {
  useEffect(() => {
    const s = document.createElement('script')
    s.text = `(function(i,s,o,g,r,a,m){i["DaoVoiceObject"]=r;i[r]=i[r]||function(){(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)})(window,document,"script","https://widget.daovoice.io/widget/${webConfig.liveChatId}.js","daovoice");`
    document.getElementsByTagName('head')[0].appendChild(s)

    const call = document.createElement('script')
    call.text = `
      daovoice('init', {
        app_id: '${webConfig.liveChatId}',
        user_id: '${currentUser.id}', // 必填: 该用户在您系统上的唯一ID
        email: '${currentUser.email}', // 选填:  该用户在您系统上的主邮箱
        name: '${currentUser.displayName}', // 选填: 用户名
        phone: '${currentUser.phone}', // 选填: 电话
        company: {
          company_name: '${env.company?.name}', // 必填,公司
          company_id: '${env.company?.id}', // 可选，公司唯一标识
          company_monthly_spend: '', // 可选，公司月付费
          company_plan: '${env.company?.account_id}', // 可选，公司套餐
          company_created_at: 1450409868, // 可选公司创建时间
        }
      });
      daovoice('onShow', function() {
        window['__buryPoint__']({
          category: '悬浮按钮',
          label: '客服'
        })
      })
    `
    document.body.appendChild(call)
  }, [])

  useEffect(() => {
    const call = document.createElement('script')
    call.text = `
      daovoice('update', {
        app_id: '${webConfig.liveChatId}',
        user_id: '${currentUser.id}', // 必填: 该用户在您系统上的唯一ID
        email: '${currentUser.email}', // 选填:  该用户在您系统上的主邮箱
        name: '${currentUser.displayName}', // 选填: 用户名
        phone: '${currentUser.phone}', // 选填: 电话
        company: {
          company_name: '${env.company?.name}', // 必填,公司
          company_id: '${env.company?.id}', // 可选，公司唯一标识
          company_monthly_spend: '', // 可选，公司月付费
          company_plan: '${env.company?.account_id}', // 可选，公司套餐
          company_created_at: 1450409868, // 可选公司创建时间
        }
      });`
    document.body.appendChild(call)
  }, [currentUser.displayName, currentUser.phone, env.company?.id])

  return <></>
}
