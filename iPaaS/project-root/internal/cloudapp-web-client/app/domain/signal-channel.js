import EventEmitter from 'eventemitter3'

export const getQueryVariable = variable => {
  var query = window.location.search.substring(1)
  var vars = query.split('&')
  for (var i = 0; i < vars.length; i++) {
    var pair = vars[i].split('=')
    if (pair[0] == variable) {
      return decodeURIComponent(pair[1])
    }
  }
  return false
}

export default class SignalChannel extends EventEmitter {
  constructor() {
    super()
    var url = ''
    var roomid = getQueryVariable('room_id')
    var signal = getQueryVariable('signal')
    var ishttps = document.location.protocol === 'https:'

    url = ''
      .concat(ishttps ? 'wss' : 'ws', '://')
      .concat(signal, '/signal/')
      .concat(roomid)

    const ws = new WebSocket(url)
    this._ws = ws
    ws.onopen = event => {
      this.emit('onopen', event)
    }
    ws.onmessage = event => {
      this.emit('onmessage', event)
    }
    ws.onerror = event => {
      this.emit('onerror', event)
    }
    ws.onclose = event => {
      this.emit('onclose', event)
    }
  }

  send(data) {
    this._ws.send(data)
  }
  close() {
    this._ws.close()
  }
}
