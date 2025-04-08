import React from 'react'
import { message } from 'antd'
import { Modal } from '@/components'
import SysConfig from '@/domain/SysConfig'
import { AUDIT_REQUEST_TYPE } from '@/constant'
import axios from 'axios'
import { Form } from './ApproverForm'
import { currentUser } from '@/domain'

const CancelToken = axios.CancelToken
let cancelRequests = new Map()
/**
 * 不需要全局的错误提示消息时，请求需传入disableErrorMessage: true
 * 需要根据不同错误码做不同处理时，自行进行try catch
 * example:
 *  try {
 *    await Http.post('/api/test/url', data, { disableErrorMessage: true })
 *  } catch(e) {
 *    // custom error handle
 * }
 */

const findAuditRequest = config => {
  const requestTypes = Object.values(AUDIT_REQUEST_TYPE)
  for (let i = 0; i < requestTypes.length; i++) {
    const curr = requestTypes[i]
    if (
      config.url.replace(/^\/api\/v1/, '') === curr.url &&
      config.method === curr.method.toLowerCase()
    ) {
      return curr
    }
  }

  return null
}

const requestParamsKey = {
  USER: 'user_approve_info_request',
  ROLE: 'role_approve_info_request'
}

let errorMessageVisible = false
let messageVisible = false
let networkErrorVisible = false
let conflictModalVisible = false

export default instance => {
  // 三员管理模式下，拦截相关请求，发起申请, 注意超级管理员的操作，不进行拦截
  instance.interceptors.request.use(
    async config => {
      if (SysConfig.enableThreeMemberMgr && !currentUser?.isSuperAdmin) {
        if (config.url.includes('/user/') || config.url.includes('/role/')) {
          const auditRequest = findAuditRequest(config)

          if (auditRequest) {
            const res = await Modal.show({
              title: '发起操作申请',
              width: 600,
              bodyStyle: { height: 250 },
              content: ({ onOk, onCancel }) => (
                <Form onOk={onOk} onCancel={onCancel} />
              ),
              footer: null
            })

            if (!res) {
              const cancelRequest = cancelRequests.get(config.cancelRequestKey)
              cancelRequest && cancelRequest()
              return null
            }

            const body = {
              approve_type: auditRequest.approve_type,
              approve_user_id: res.id, // 审批人 id
              approve_user_name: res.name // 审批人 name
            }
            const key = requestParamsKey[auditRequest.type]

            body[key] = {
              ...config.params,
              ...config.data
            }

            let success = true
            // 发请求
            try {
              const response = await fetch('/api/v1/approve/apply', {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json; charset=utf-8',
                  'x-userid': localStorage.getItem('userId')
                },
                body: JSON.stringify(body)
              })

              const res = await response.json()
              if (res.success) {
                message.success('发起操作申请成功')
              } else {
                message.error(res.message || '发起操作申请失败')
                success = false
              }
            } catch (e) {
              message.error('发起操作申请失败')
              success = false
            }

            // TODO: 取消请求，会失败, 出错，需查明原因
            // const cancelRequest = cancelRequests.get(config.cancelRequestKey)
            // cancelRequest && cancelRequest()
            config.fake = true
            config.url = '/fake/404'
            config.fakeSuccess = success

            return config
          } else {
            return config
          }
        }

        return config
      } else {
        return config
      }
    },
    error => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.request.use(
    config => {
      if (navigator.onLine) {
        return config
      } else {
        message.error('您没有联网，请检查网络连接')
        return Promise.reject('no network')
      }
    },
    error => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.request.use(
    config => {
      if (config.method && config.method.toLowerCase() === 'get') {
        config.params = config.params || {}
        config.params['__timestamp__'] = Date.now()
      }
      return config
    },
    error => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.request.use(
    config => {
      config.headers['x-csrf-token'] = ''
      config.headers['x-userid'] = localStorage.getItem('userId')
      const groupId = -1
      config.headers['x-groupid'] = groupId
      return config
    },
    error => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.request.use(
    config => {
      config.cancelRequestKey = config.url + '_' + config.method
      config.cancelToken = new CancelToken(function executor(c) {
        cancelRequests.set(config.cancelRequestKey, c)
      })
      return config
    },
    error => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.response.use(
    response => response,
    error => {
      // if (error.message === 'Network Error' && !networkErrorVisible) {
      //   networkErrorVisible = true
      //   message.error('网络异常').promise.finally(() => {
      //     networkErrorVisible = false
      //   })
      // }

      const { response } = error
      if (response) {
        const { formatErrorMessage, disableErrorMessage } = response.config

        switch (response.status) {
          case 502: {
            if (!errorMessageVisible) {
              errorMessageVisible = true
              message
                .error('服务异常，请联系系统管理员。', 5)
                .promise.finally(() => {
                  errorMessageVisible = false
                })
            }
            break
          }
          case 409: {
            if (!conflictModalVisible) {
              conflictModalVisible = true
              return Modal.showConfirm({
                cancelButtonProps: { style: { display: 'none' } },
                closable: false,
                content: '用户登录冲突，请重新登录'
              })
                .then(() => {
                  localStorage.setItem('needLogin', 'true')
                  location.reload()
                })
                .finally(() => {
                  conflictModalVisible = false
                })
            }
            break
          }
          case 401:
            if (!disableErrorMessage && !messageVisible) {
              messageVisible = true
              message.error('未登录').promise.finally(() => {
                messageVisible = false
              })
            }

            localStorage.setItem('needLogin', 'true')
            location.reload()
            break
          case 403:
            if (!disableErrorMessage && !messageVisible) {
              messageVisible = true
              message.error('没有权限').promise.finally(() => {
                messageVisible = false
              })
            }
            localStorage.setItem('needLogin', 'true')
            location.reload()
            break
          default:
            if (
              !disableErrorMessage &&
              response &&
              response.data &&
              response.data.message
            ) {
              let msg = response.data.message
              if (formatErrorMessage) {
                msg = formatErrorMessage(msg)
              }
              message.error(msg)
            }
            break
        }
      }

      return Promise.reject(error)
    }
  )

  instance.interceptors.response.use(
    response => {
      return response
    },
    error => {
      const { response } = error

      if (response.config.fake && response.status === 404) {
        return Promise.reject({
          fake: true,
          success: response.config.fakeSuccess
        })
      }
      return Promise.reject(error)
    }
  )
}
