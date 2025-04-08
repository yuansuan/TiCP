type WebSocketEventListener = (event: WebSocketEvent) => void

export enum WebSocketReadyState {
  Connecting = 0,
  Open = 1,
  Closing = 2,
  Closed = 3
}

export enum WebSocketEventType {
  Open = 'open',
  Close = 'close',
  Message = 'message',
  Error = 'error'
}

export interface WebSocketEvent {
  type: WebSocketEventType
  data?: any
}

class WebSocketClient {
  private url: string
  private socket: WebSocket | null
  private listeners: Record<WebSocketEventType, WebSocketEventListener[]>
  private reconnectInterval: number
  private reconnectAttempts: number
  private currentReconnectAttempts: number

  constructor(
    url: string,
    reconnectInterval: number = 5000,
    reconnectAttempts: number = 5
  ) {
    this.url = url
    this.socket = null
    this.listeners = {
      open: [],
      close: [],
      message: [],
      error: []
    }
    this.reconnectInterval = reconnectInterval
    this.reconnectAttempts = reconnectAttempts
    this.currentReconnectAttempts = 0
  }

  connect(): void {
    if (!this.url) return
    this.socket = new WebSocket(this.url)
    this.socket.onopen = this.handleOpen.bind(this)
    this.socket.onclose = this.handleClose.bind(this)
    this.socket.onmessage = this.handleMessage.bind(this)
    this.socket.onerror = this.handleError.bind(this)
  }

  disconnect(): void {
    if (this.socket) {
      this.socket.close()
      this.socket = null
      this.currentReconnectAttempts = 0 // Reset reconnect attempts
    }
  }

  send(data: any): void {
    if (this.socket && this.socket.readyState === WebSocketReadyState.Open) {
      this.socket.send(JSON.stringify(data))
    }
  }

  on(event: WebSocketEventType, listener: WebSocketEventListener): void {
    this.listeners[event].push(listener)
  }

  off(event: WebSocketEventType, listener: WebSocketEventListener): void {
    this.listeners[event] = this.listeners[event].filter(l => l !== listener)
  }

  private handleOpen(event: Event): void {
    this.emitEvent(WebSocketEventType.Open)
    this.currentReconnectAttempts = 0 // Reset reconnect attempts on successful connection
  }

  private handleClose(event: CloseEvent): void {
    this.emitEvent(WebSocketEventType.Close)

    if (this.currentReconnectAttempts < this.reconnectAttempts) {
      // Attempt to reconnect
      setTimeout(() => {
        this.connect()
        this.currentReconnectAttempts++
      }, this.reconnectInterval)
    }
  }

  private handleMessage(event: MessageEvent): void {
    const data = JSON.parse(event.data)
    this.emitEvent(WebSocketEventType.Message, data)
  }

  private handleError(event: Event): void {
    this.emitEvent(WebSocketEventType.Error)
  }

  private emitEvent(type: WebSocketEventType, data?: any): void {
    const event: WebSocketEvent = { type, data }
    this.listeners[type].forEach(listener => listener(event))
  }
}

export default WebSocketClient
