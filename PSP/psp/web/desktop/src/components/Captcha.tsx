/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Button } from '@/components'
import { message, Input } from 'antd'
import { Http } from '@/utils'
import { env } from '@/domain'

type Props = {
  phone: string
  captcha: string
  setCaptcha: (code: string) => void
  skip_check?: boolean
}

export const Captcha = observer(function Captcha({
  phone,
  captcha,
  setCaptcha,
  skip_check,
}: Props) {
  const invalidPhone = useMemo(() => {
    if (!/^1\d{10}/.test(phone)) {
      return '手机号不合法'
    }

    return false
  }, [phone])

  const state = useLocalStore(() => ({
    timer: null,
    counter: 0,
    updateCounter(counter) {
      this.counter = counter
    },
    loading: false,
    updateLoading(loading) {
      this.loading = loading
    },
  }))

  const getCode = async e => {
    e.stopPropagation()

    await Http.get('/go/send_code', {
      params: {
        phone,
        pid: env.productId,
        skip_check: skip_check ? '1' : undefined,
      },
      disableErrorMessage: true,
    })
      .then(res => {
        message.success('验证码发送成功')
        state.updateCounter(60)
        clearTimeout(state.timer)
        state.timer = setTimeout(function count(this: any) {
          state.updateCounter(state.counter - 1)
          if (state.counter > 0) {
            state.timer = setTimeout(count.bind(this), 1000)
          }
        }, 1000)
      })
      .catch(err => {
        if ('' + err.code === '90027') {
          message.error('手机号已存在')
        } else {
          message.error('发送失败')
        }
      })
  }

  return (
    <Input
      type='text'
      placeholder='请输入验证码'
      value={captcha}
      onChange={e => setCaptcha(e.target.value)}
      addonAfter={
        state.counter > 0 ? (
          <span>{state.counter} s</span>
        ) : (
          <Button
            type='secondary'
            size='small'
            disabled={invalidPhone}
            onClick={getCode}>
            发送验证码
          </Button>
        )
      }
    />
  )
})
