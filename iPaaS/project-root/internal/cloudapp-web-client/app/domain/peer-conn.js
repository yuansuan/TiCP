import EventEmitter from 'eventemitter3'
import UAParser from 'ua-parser-js'
import log from '@/infra/log'

export default class PeerConn extends EventEmitter {
  constructor(config) {
    super()
    this._pc = new RTCPeerConnection(config)
    this._dataChannel = null
    this.lastMessageTime = new Date()
    //  是数据通道调用时的事件处理器
    this._pc.ondatachannel = ({ channel }) => this._onDataChCreated(channel)
    this._pc.onicecandidate = event => {
      this.emit('on-ice-candidate', event)
    }
    // 收到addstream 事件调用时
    this._pc.onaddstream = event => {
      this.emit('on-add-stream', event)
    }
    this._pc.oniceconnectionstatechange = () => {
      this.emit('on-ice-connection-state-change', this._pc.iceConnectionState)
    }
    this._reportHandler = this.getReportHandler()
  }
  setConfiguration(config) {
    this._pc.setConfiguration(config)
  }
  close() {
    this._pc.ondatachannel = null
    this._pc.onicecandidate = null
    this._pc.onaddstream = null
    this._pc.oniceconnectionstatechange = null
    this._pc.close()
    this._pc = null
  }
  restartIce() {
    this._pc.restartIce()
  }

  _onDataChCreated(channel) {
    log.info('data channel created')
    this._dataChannel = channel
    this._dataChannel.onmessage = e => {
      this.emit('on-data-channel-message', e)
    }
    this.emit('on-data-channel')
  }

  async onGetDescription(desc) {
    await this._pc.setRemoteDescription(desc)
    const answerDesc = await this._pc.createAnswer()
    await this._pc.setLocalDescription(answerDesc)
    return answerDesc
  }

  handleIncomingICE(ice) {
    // log.info('got remote ice candidate')
    const candidate = new RTCIceCandidate(ice)
    this._pc.addIceCandidate(candidate)
  }

  getReportHandler() {
    var _packetsLost = 0
    var _packetsReceived = 0
    var _appDelay = 0
    return function (report) {
      let packetsLost = report.packetsLost
      let packetsReceived = report.packetsReceived
      let jitter = report.jitter

      if (
        typeof packetsLost !== 'undefined' &&
        typeof packetsReceived !== 'undefined'
      ) {
        const deltaPacketsReceived = packetsReceived - _packetsReceived

        if (!isNaN(deltaPacketsReceived)) {
          const deltaPacketsLost = packetsLost - _packetsLost
          const packetsLostRate =
            deltaPacketsLost / (deltaPacketsReceived + deltaPacketsLost)
          const appDelay = Math.floor(
            (_appDelay * 15 +
              jitter *
              (deltaPacketsReceived + 10) *
              (5 + packetsLostRate * 20)) /
            16
          )

          _appDelay = appDelay
          _packetsLost = packetsLost
          _packetsReceived = packetsReceived
          return appDelay
        }
      }
    }
  }
  async getLatencyFromStats() {
    const report = await this.getStats()
    const parser = new UAParser()
    const browser = parser.getBrowser().name
    if (browser === 'Firefox') {
      for (let type in report) {
        if (type === 'inbound-rtp') {
          return this._reportHandler(report[type])
        }
      }
    } else {
      for (let type in report) {
        for (let item in report[type]) {
          if (item === 'googCurrentDelayMs') {
            return report[type][item]
          }
        }
      }
    }
  }
  async getStats() {
    const parser = new UAParser()
    const browser = parser.getBrowser().name
    var report = {}
    if (browser === 'Firefox') {
      const stats = await this._pc.getStats()
      for (let v of stats.values()) {
        report[v.type] = v
      }
      return report
    } else {
      return new Promise(resolve => {
        this._pc.getStats(function (stats) {
          stats.result().forEach(function (res) {
            var item = {}
            res.names().forEach(function (name) {
              item[name] = res.stat(name)
            })
            item.id = res.id
            item.type = res.type
            item.timestamp = res.timestamp
            report[item.type] = item
          })
          resolve(report)
        })
      })
    }
  }
  /*
   * @param { string } msg send the message to dataChannel
   */
  send(msg) {
    if (this._dataChannel && this._dataChannel.readyState === 'open') {
      this._dataChannel.send(msg)
      this.emit('on-data-channel-message-sent')
      if (localStorage.getItem('debug_data_channel')) {
        console.log('#SEND# DATA CHANNEL', msg)
      }
    }
  }
}
