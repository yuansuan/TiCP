/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { notification } from 'antd'
import { currentUser } from '@/domain'
import isDev from '@/utils/isDev'
export class WebSocketClient {
  private url: string
  private websocket: WebSocket
  private retryCount: number
  private maxRetries: number
  private heartbeatInterval: number
  private heartbeatTimer: number
  private retryTimer: number

  constructor(
    url: string,
    maxRetries: number = 3,
    heartbeatInterval: number = 30000
  ) {
    this.url = url
    this.websocket = null
    this.retryCount = 0
    this.maxRetries = maxRetries
    this.heartbeatInterval = heartbeatInterval
    this.heartbeatTimer = null
  }

  connect(): void {
    if (!this.url) return
    this.websocket = new WebSocket(this.url)

    this.websocket.onopen = () => {
      this.startHeartbeat()
    }

    this.websocket.onmessage = event => {
      try {
        const message = event.data

        if (message === 'Pong') {
          clearTimeout(this.retryTimer)
          return
        }
        const parseMsg = JSON.parse(message)
        // 处理收到的消息
        if (currentUser?.id === parseMsg?.user_id) {
          notification.info({
            message: '系统通知',
            description: parseMsg?.content,
            placement: 'bottomRight'
            // onClick: () => history.push(`/job/${job_id}`)
          })
        }
      } catch (err) {
        console.log('err: ', err)
      }
    }

    this.websocket.onclose = event => {
      this.stopHeartbeat()

      if (this.retryCount < this.maxRetries) {
        this.retryCount++
        this.retryTimer = setTimeout(() => {
          this.connect()
        }, 2000) // 2秒后重新连接
      } else {
        console.log('达到最大重试次数，停止连接尝试')
      }
    }

    this.websocket.onerror = error => {
      console.error('WebSocket连接发生错误:', error)
    }
  }

  send(message: string): void {
    if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
      this.websocket.send(message)
    } else {
      console.error('WebSocket连接未建立或已关闭，无法发送消息')
    }
  }

  private startHeartbeat(): void {
    this.heartbeatTimer = setInterval(() => {
      if (this.websocket.readyState === WebSocket.OPEN) {
        this.websocket.send('Ping')
      }
    }, this.heartbeatInterval)
  }

  private stopHeartbeat(): void {
    clearInterval(this.heartbeatTimer)
  }
}

const ishttps = window.location.protocol === 'https:'
// 10.0.7.146  192.168.111.239
const noticeUrl = `${ishttps ? 'wss' : 'ws'}://${
  window.location.hostname === 'localhost'
    ? '192.168.111.239/ws/v1/notice/consumer'
    : window.location.host + '/ws/v1/notice/consumer'
}`

const websocketClient = new WebSocketClient(noticeUrl, 500, 5000)
!isDev && websocketClient.connect()

// 发送消息
// websocketClient.send('Hello, WebSocket!')
