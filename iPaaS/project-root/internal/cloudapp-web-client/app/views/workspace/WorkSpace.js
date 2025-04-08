import React from 'react'
import 'webrtc-adapter'
import { connect } from 'react-redux'
import Pipeline from '@/domain/webrtc-pipeline'
import log from '@/infra/log'
import './workspace.scss'

class WorkSpace extends React.Component {
  constructor(props) {
    super(props)
    // 通过React.createRef()创建ref,挂载到组件上
    this.videoPlayer = React.createRef()
    this.videoContainer = React.createRef()
  }
  isEditingMode = false
  state = {
    fatalError: null,
    loading: true,
    signalRetry: 0
  }
  componentDidMount() {
    const search = new URLSearchParams(window.location.search)
    this.isEditingMode = !!(search.get('editing') === 'true')
    const userId = Number(search.get('user_id'))
    const workTaskId = Number(search.get('worktask_id'))
    const title = search.get('title')
    const roomId = search.get('room_id') || workTaskId
    const peerId = this.getPeerId(userId, workTaskId, roomId)

    if (title) {
      document.title = title
    }

    this.startApp(roomId, peerId, workTaskId)
  }
  getPeerId(userId, workTaskId, roomId) {
    function padZero(value, len) {
      const v = value.toString()
      if (v.length < len) {
        return '0'.repeat(len - v.length) + v
      }
      return v
    }
    if (userId && workTaskId) {
      return padZero(userId, 4) + padZero(workTaskId, 4)
    }
    return padZero(userId, 4) + padZero(roomId, 4)
  }
  startApp(roomId, peerId, workTaskId) {
    const pline = new Pipeline({
      videoPlayer: this.videoPlayer.current,
      peerId,
      roomId
    })
    window.webrtcPipeline = pline
    pline.start()
    pline.on('signal-close', () => {})
    pline.on('room-full-error', () => {
      this.setState({
        fatalError: {
          reason: 'room_full',
          message: '当前任务已经开启，请不要重复打开。'
        }
      })
    })
    pline.on('room-no-server', () => {
      const retryCount = this.state.signalRetry
      if (retryCount < 10) {
        setTimeout(() => {
          pline.startSingaling()
          this.setState({ signalRetry: this.state.signalRetry + 1 })
        }, 1000)
      } else {
        this.setState({
          fatalError: {
            reason: 'room_has_no_server',
            message: '远程服务已经关闭！'
          }
        })
        this.cleanWorkSpace()
      }
    })
    pline.on('app-started', () => {
      this.setState({ loading: false })
      // 初始化自适应
      this.onAppResize(pline)

      // ============================
      // this.startMeasureLatency(pline)
      if (!this.isEditingMode) {
        // this.startMonitorMessage()
      }
      // WorkTaskService.reportAppMetric({
      //   action_name: 'showVideo',
      //   success: true,
      //   task_id: workTaskId
      // })
    })
    pline.on('ice-status-changed', state => {
      if (state === 'disconnected' || state === 'failed') {
        this.props.updateMessage({ type: 'error', text: '网络不稳定！' })
        log.error(`ice connection state change - ${state}`)
        pline.restartIce()
      } else {
        this.props.updateMessage({ text: null })
        log.info(`ice connection state change - ${state}`)
      }
    })
    pline.on('copy-out-text', text => {
      this.props.copyOutText(text)
      // 从3D云应用内复制到外部系统中
      navigator.clipboard.writeText(text).catch(() => {})
    })
    pline.on('on-data-channel-message-sent', () => {
      this.props.updateLastMessageSent(new Date())
    })

    // 从外部复制到3D云应用中
    window.addEventListener('focus', () => {
      navigator.clipboard
        .readText()
        .then(text => {
          if (window.webrtcPipeline) {
            pline.sendPasteText(text)
          }
        })
        .catch(() => { })
    })
  }

  startMonitorMessage = () => {
    this._dh_timer = setInterval(() => {
      const duration =
        new Date().getTime() - this.props.webrtc.lastMessageSent.getTime()
      if (duration > 9 * 60 * 1000) {
        this.props.userInactive(true)
      }
    }, 10 * 1000)
  }
  startMeasureLatency = pline => {
    this._stats_timer = setInterval(async () => {
      const latentcy = await pline.getLatencyFromStats()
      this.props.setLatency(latentcy)
    }, 1500)
  }
  onAppResize(pline) {
    const container = this.videoContainer.current
    const params = {
      width: document.body.clientWidth,
      height: document.body.clientHeight
    }
    if (pline) {
      pline.sendResolution(params)
    }
  }
  cleanWorkSpace() {
    this.props.updateMessage({ text: null })
    clearInterval(this._stats_timer)
    clearInterval(this._dh_timer)
  }
  componentWillUnmount() {
    this.cleanWorkSpace()
  }
  render() {
    const { fatalError, loading } = this.state
    return (
      <div ref={this.videoContainer} id="video-container">
        {/* {loading && <Particle />} */}
        {!fatalError && (
          <div ref={this.videoPlayer} id="video-player" tabIndex="0" />
        )}
        {loading && !fatalError && (
          <div id="loading-placeholder">
            <img src="static/loading.svg" />
          </div>
        )}
        {fatalError && (
          <div className="fatal-error-tip">
            <h2>{fatalError.message}</h2>
          </div>
        )}
      </div>
    )
  }
}
function mapStateToProps(state) {
  return {
    webrtc: state.webrtc
  }
}

function mapDispatchToProps(dispatch) {
  return {
    updateLastMessageSent: lastMessageSent => {
      dispatch({ type: 'LAST_MESSAGE_SENT', lastMessageSent })
    },
    userInactive: inactive => {
      dispatch({ type: 'USER_INACTIVE', inactive })
    },

    copyOutText: text => {
      dispatch({ type: 'COPY_OUT_TEXT', text })
    },
    setLatency: latency => {
      dispatch({ type: 'SET_LATENCY', latency })
    },
    updateMessage: message => {
      dispatch({ type: 'updateMessage', message })
    }
  }
}
export default connect(mapStateToProps, mapDispatchToProps)(WorkSpace)
