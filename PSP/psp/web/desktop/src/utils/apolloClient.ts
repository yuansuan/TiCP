/* Copyright (C) 2016-present, Yuansuan.cn */

import { single } from '@/utils'
import { InMemoryCache, from } from '@apollo/client'
import { env } from '@/domain'
import { Modal } from '@/components'
import { createErrorMiddleware, createGraphqlClient, createHeaderMiddleware } from './graphqlClient'

export const apolloClient = createGraphqlClient({
  uri: '/graphql',
  cache: new InMemoryCache(),
  link: link =>
    from([
      createHeaderMiddleware(() => ({
        'X-Project-Id': env?.project?.id,
        'X-Company-Id': env?.company?.id,
        'X-Account-Id': env?.accountId
      })),
      createErrorMiddleware(exception => {
        if (exception?.status === 401) {
          window.location.replace(`/api/sso/login${window.location.hash}`)
        } else if (exception?.status === 409) {
          single('login-conflict-modal', () =>
            Modal.showConfirm({
              cancelButtonProps: { style: { display: 'none' } },
              closable: false,
              content: '用户登录冲突，请重新登录'
            }).then(() => {
              env.logout()
            })
          )
        }

        const res = exception?.response
        if (res) {
          const { errorCode } = res
          // not in project
          if (errorCode === 1001) {
            single('in-project-modal', () =>
              Modal.showConfirm({
                title: '工作空间不存在',
                content: '您已退出当前工作空间',
                closable: false,
                CancelButton: null
              }).then(async () => {
                location.replace('/')
              })
            )
          } else if (errorCode === 120006) {
            single('user-not-exist-modal', () =>
              Modal.showConfirm({
                title: '消息提示',
                content: '账号异常，点击确认重新登录',
                closable: false,
                CancelButton: null
              }).then(() => {
                env.logout()
              })
            )
          }
        }
      }),
      link
    ])
})
