import log from '@/infra/log'
import EventEmitter from 'eventemitter3'
import PeerConn from '@/domain/peer-conn'
import CustomCursor from '@/infra/custom-cursor'
import JSONRPC from '@/infra/jsonrpc'
import idGen from '@/infra/id-gen'
import { makeInputInterceptor } from '@/infra/input-intercept'
import SignalChannel from '@/domain/signal-channel'
const protocol = new JSONRPC(idGen)

class WebRTCPipeline extends EventEmitter {
  constructor(options) {
    super()
    this._videoPlayer = options.videoPlayer
    this._peerId = options.peerId
    this._roomId = options.roomId
    this._pc = null
    this._interceptor = null
    this._signaling = null
    this.initVideoPlayer()
    this._cursor = new CustomCursor(this._video)
  }
  initVideoPlayer() {
    this._videoPlayer.onclick = () => {
      this._videoPlayer.focus()
    }
    this._videoPlayer.addEventListener('blur', e => {
      this.sendFocus()
    })
    this._videoPlayer.addEventListener('mouseout', e => {
      this.sendFocus()
    })
    this._video = document.createElement('video')
    this._videoPlayer.appendChild(this._video)
    this._video.setAttribute('id', 'stream')
    this._video.muted = true
    this._video.autoplay = true
    this._video.addEventListener('loadeddata', () => {
      this._video.style.display = 'block'
      this.emit('app-started')
    })
    this._video.addEventListener('resize', e => {
    })
  }
  resetPlayerSize() {
    const width = this._video.videoWidth
    const height = this._video.videoHeight
    if (width > 100 && height > 100) {
      this._video.width = this._video.videoWidth
      this._video.height = this._video.videoHeight
      this._videoPlayer.style.width = this._video.videoWidth + 'px'
      this._videoPlayer.style.height = this._video.videoHeight + 'px'
    }
  }
  sendFocus() {
    const msg = protocol.encode('focus', { value: 0 })
    if (this._pc) {
      this._pc.send(msg)
    }
  }
  sendCAD() {
    const cadData = protocol.encode('cad', 'CAD')
    log.info('will send cad msg', cadData)
    this._pc.send(cadData)
  }
  sendBitrate(value) {
    const msg = protocol.encode('bitrate', { max_bitrate: (value * 1000000) })
    if (this._pc) {
      this._pc.send(msg)
    }
  }

  _send(msg) {
    this._signaling.send(msg)
  }

  _handleError(...msg) {
    log.error(msg)
  }

  /**
   * @param {string} msg
   */
  _onDataChMessage(rawMsg) {
    log.debug('datachannel message: ', rawMsg)
    const msg = protocol.decode(rawMsg)
    const method = msg[0]
    const infos = msg[1]

    if (method === 'cursor') {
      const idx = infos.index
      if (infos.hidden) {
        this._cursor.hideCursor(idx)
      } else if (!infos.data) {
        this._cursor.setImage(idx)
      } else {
        const hotspotX = infos.x
        const hotspotY = infos.y
        const { data } = infos
        this._cursor.setImage(idx, hotspotX, hotspotY, data)
      }
    } else if (method === 'copy') {
      const text = infos.msg
      log.info('will show copy text:', text)
      this.emit('copy-out-text', text)
    } else if (method === 'resolution') {
      // TODO 分辨率
      log.info('RECEIVED RESOLUTION', infos)
      // const { height, width } = infos
      // this.resetPlayerSize()
      // this._video.width = width
      // this._video.height = height
      // const videoPlayer = document.getElementById('video-player')
      // videoPlayer.style.width = width + 'px'
      // videoPlayer.style.height = height + 'px'
    }
  }

  _onMessage(event) {
    log.debug('receive msg from server:', event.data)

    switch (event.data) {
      case 'HELLO':
        this._send(`ROOM ${this._roomId}`)
        return
      default:
        if (event.data.startsWith('ROOM_OK')) {
          log.info(
            'start remote desktop server msg is sent, waiting for call',
            event.data
          )
          return
        }
        if (event.data.startsWith('TURN_CRED_INFO')) {
          const turnInfo = event.data.trim().split(' ')
          let iceServers = []
          if (turnInfo.length >= 4) {
            for (let i = 1; i < turnInfo.length; i += 3) {
              let iserver = {
                urls: [turnInfo[i]],
                username: turnInfo[i + 1],
                credential: turnInfo[i + 2]
              }
              iceServers.push(iserver)
            }
          }
          this.startPeerConnection({ iceServers })
          return
        }
        if (event.data.startsWith('ROOM_PEER_JOINED')) {
          log.info('server join the room')
          return
        }
        if (event.data.startsWith('ROOM_PEER_LEFT')) {
          log.info('server left the room')
          this._handleError('server offline', event.data)
          this.emit('room-no-server')
          return
        }
        if (event.data.startsWith('WARNING')) {
          log.info('got warning from server', event.data)
          return
        }
        if (event.data.startsWith('ERROR')) {
          this._handleError('got error from server', event.data)
          if (event.data == 'ERROR room_full') {
            this.emit('room-full-error')
          } else if (event.data == 'ERROR room has no server yet') {
            this.emit('room-no-server')
          }
          return
        }
    }

    let msg
    // process webrtc messages
    try {
      msg = JSON.parse(event.data)
    } catch (e) {
      if (e instanceof SyntaxError) {
        this._handleError(`fail to parse incoming JSON, data=${event.data}`)
      } else {
        this._handleError(`Unknown error parsing response, data=${event.data}`)
      }
      return
    }

    if (msg.sdp) {
      this._pc.onGetDescription(msg.sdp).then(desc => {
        const sdp = { sdp: desc }
        this._send(JSON.stringify(sdp))
      })
    } else if (msg.ice) {
      this._pc.handleIncomingICE(msg.ice)
    } else {
      log.info('unknown incoming json', msg)
    }
  }

  _onError(event) {
    log.info('signal server error event', event)
    this.emit('signal-error')
  }

  _onClose(event) {
    log.info('disconnected from signal server', event)
    this.emit('signal-close')
  }

  _onOpen(event) {
    this._sendHello()
  }
  _sendHello() {
    this._send(`HELLO ${this._peerId} client`)
  }

  start() {
    this.startSingaling()
  }
  startPeerConnection(config) {
    this._pc = new PeerConn(config)
    this._pc.on('on-ice-candidate', event => {
      if (event.candidate === null) {
        return
      }
      this._send(JSON.stringify({ ice: event.candidate }))
    })

    this._pc.on('on-add-stream', event => {
      this._onRemoteStreamAdd(event)
    })
    this._pc.on('on-data-channel', () => {
      this._interceptor = makeInputInterceptor(this._pc, [this._videoPlayer])
      this._interceptor.start()
    })
    this._pc.on('on-data-channel-message', event => {
      this._onDataChMessage(event.data)
    })
    this._pc.on('on-data-channel-message-sent', event => {
      this.emit('on-data-channel-message-sent')
    })
    this._pc.on('on-ice-connection-state-change', state => {
      this.emit('ice-status-changed', state)
    })
  }
  startSingaling() {
    this._signaling = new SignalChannel()
    this._signaling.on('onopen', this._onOpen.bind(this))
    this._signaling.on('onclose', this._onClose.bind(this))
    this._signaling.on('onmessage', this._onMessage.bind(this))
    this._signaling.on('onerror', this._onError.bind(this))
  }
  restartIce() {
    this._pc.restartIce()
  }
  cleanPC() {
    this._pc.close()
    this._video.srcObject.stop()
    this._video.remove()
    if (this._interceptor) {
      this._interceptor.stop()
      this._interceptor = null
    }
  }

  _onRemoteStreamAdd(event) {
    const videoTracks = event.stream.getVideoTracks()
    const audioTracks = event.stream.getAudioTracks()
    if (videoTracks.length === 0) {
      this._handleError('stream with unknown tracks added')
      return
    }

    log.info(
      `incoming stream: ${videoTracks.length} video and` +
      `${audioTracks.length} audio tracks`
    )
    this._video.srcObject = event.stream
  }

  async getLatencyFromStats() {
    return await this._pc.getLatencyFromStats()
  }

  /**
   * @param {string} text
   */
  sendPasteText(text) {
    this._pc.send(protocol.encode('copy', { msg: text }))
  }

  /**
   * @param {Object} resolution
   */
  sendResolution(resolution) {
    // NOTE 告知后端改变分辨率
    if (this._pc) {
      const data = protocol.encode('resolution', resolution)
      log.info('REQ_CHANGE_RESOLUTION', data)
      this._pc.send(data)
    }
  }
  sendShowDesktop() {
    // NOTE 告知后端，显示桌面
    if (this._pc) {
      const data = protocol.encode('show_desktop', {})
      log.info('REQ_SHOW_DESKTOP', data)
      this._pc.send(data)
    }
  }
}

export default WebRTCPipeline
