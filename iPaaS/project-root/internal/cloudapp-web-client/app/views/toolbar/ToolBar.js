import React from 'react'
import { connect } from 'react-redux'
import screenfull from 'screenfull'
import dayjs from 'dayjs'
import './toolbar.scss'
import PasteModal from './PasteModal'
import EndModal from './EndModal'
import InactiveTipModal from './InactiveTipModal'
import WorkTaskService from '@/services/WorkTaskService'
import 'rc-slider/assets/index.css';
import 'rc-tooltip/assets/bootstrap.css';
import Slider from 'rc-slider'
import Tooltip from 'rc-tooltip'

const Handle = Slider.Handle

class ToolBar extends React.Component {
  autoAdjustResolution = true
  resolutionIdSelected = 'resolution-auto-adjust'
  maxBitrate = 30

  state = {
    isFullscreen: false,
    minimized: false,
    showPasteModal: false,
    showEndModal: false,
    durationText: '--:--',
    isEditingMode: false
  }
  componentDidMount() {
    screenfull.onchange(() => {
      const isFullscreen = screenfull.isFullscreen
      this.setState({ isFullscreen, minimized: isFullscreen })
    })
    const search = new URLSearchParams(window.location.search)
    const isEditingMode = !!(search.get('editing') === 'true')
    this.setState({ isEditingMode })

    // default to auto adjust resolution
    document
      .getElementById('resolution-auto-adjust')
      .classList.add('resolution-selected')
    this.handleWindowResizeDebounce()
  }
  componentWillUnmount() {
    clearInterval(this.duraion_timer)
  }
  updateDurationText = workTaskId => {
    WorkTaskService.getWorkTaskDetail(workTaskId).then(res => {
      if (res.errorCode === 991100) {
        location.href = '/#/login'
      }
      const sec = dayjs().diff(dayjs(res.data.start_time), 'second', false)
      const text = this.formatDurationText(sec)
      this.setState({ durationText: text })
    })
  }
  formatDurationText = sec => {
    const mins = Math.ceil(sec / 60)
    let minText = `${mins % 60}`
    if (minText.length < 2) {
      minText = `0${minText}`
    }
    let hText = `${Math.floor(mins / 60)}`
    if (hText.length < 2) {
      hText = `0${hText}`
    }
    return ` ${hText}:${minText}`
  }
  hideToolBar = () => {
    this.setState({ minimized: true })
  }
  restoreToolBar = () => {
    this.setState({ minimized: false })
  }
  toggleFullscreen = () => {
    screenfull.toggle()
  }
  showPasteModal = () => {
    this.setState({ showPasteModal: true })
  }
  closePasteModal = () => {
    this.setState({ showPasteModal: false })
  }
  showEndModal = () => {
    this.setState({ showEndModal: true })
  }
  closeEndModal = () => {
    this.setState({ showEndModal: false })
  }
  handleWindowResizeDebounce(width, height) {
    let resizeTimer = null
    this.resizeHandler = () => {
      clearTimeout(resizeTimer)
      resizeTimer = setTimeout(() => {
        this.adjustResolution(width, height)
      }, 500)
    }

    window.addEventListener('resize', this.resizeHandler)
  }

  onAutoAdjustResolution = e => {
    this.onResulotionTargetSelected(e)
    this.autoAdjustResolution = !this.autoAdjustResolution

    if (this.autoAdjustResolution) {
      this.adjustResolution()
      this.handleWindowResizeDebounce()
      return
    }

    if (this.resizeHandler) {
      window.removeEventListener('resize', this.resizeHandler)
    }
  }
  onSpecifySpecificResolution = (e, width, height) => {
    if (e.id === this.resolutionIdSelected) {
      return
    }

    this.onResulotionTargetSelected(e)

    if (this.autoAdjustResolution) {
      if (this.resizeHandler) {
        window.removeEventListener('resize', this.resizeHandler)
      }
      this.autoAdjustResolution = false
    }
    this.adjustResolution(width, height)
  }

  onResulotionTargetSelected(targetEle) {
    const currentClickId = targetEle.id
    if (this.resolutionIdSelected === -1) {
      targetEle.classList.add('resolution-selected')
      this.resolutionIdSelected = currentClickId
      return
    }

    if (currentClickId === this.resolutionIdSelected) {
      targetEle.classList.remove('resolution-selected')
      this.resolutionIdSelected = -1
    } else {
      targetEle.classList.add('resolution-selected')
      document
        .getElementById(this.resolutionIdSelected)
        .classList.remove('resolution-selected')
      this.resolutionIdSelected = currentClickId
    }
  }

  adjustResolution = (width, height) => {
    let params = {}
    if (height && width) {
      params = { width, height }
    } else {
      params = {
        width: document.body.clientWidth,
        height: document.body.clientHeight
      }
    }

    console.log(
      'adjust resolution, width: ',
      params.width,
      ', height: ',
      params.height
    )

    const pline = window.webrtcPipeline
    pline && pline.sendResolution(params)
  }
  showDesktop = () => {    
    const pline = window.webrtcPipeline
    pline && pline.sendShowDesktop()
  }

  closeInactiveTipModal = () => {
    this.props.userInactive(false)
    this.props.updateLastMessageSent(new Date())
    const search = new URLSearchParams(window.location.search)
    const workTaskId = Number(search.get('worktask_id'))
    WorkTaskService.activeApp(workTaskId)
  }
  getSrc(delay) {
    let src = 'static/network-'
    if (delay < 150) {
      src += '4'
    } else if (delay < 300) {
      src += '3'
    } else if (delay < 450) {
      src += '2'
    } else {
      src += '1'
    }
    src += '.svg'
    return src
  }
  formatLatency = delay => {
    const TIME_UNITS = ['ms', 's']
    if (delay > 5000) {
      return ` >5s`
    }
    let unitIdx = 0
    let displayDelay = delay
    if (delay > 1000) {
      displayDelay /= 1000
      unitIdx += 1
    }
    return ` ${displayDelay}${TIME_UNITS[unitIdx]}`
  }
  setMaxBitrate = props => {
    const { dragging, index, ...restProps } = props
    const value = this.maxBitrate || props.value
    return (
      <Tooltip
        prefixCls="rc-slider-tooltip"
        overlay={`最大速率${value}Mbps`}
        visible={dragging}
        placement="top"
        key={index}>
        <Handle value={value} {...restProps} />
      </Tooltip>
    )
  }
  onMaxBitrateChange = value => {
    const pline = window.webrtcPipeline
    pline && pline.sendBitrate(value)
    this.maxBitrate = value
  }
  render() {
    const {
      minimized,
      isFullscreen,
      showPasteModal,
      showEndModal,
      durationText,
      isEditingMode
    } = this.state
    const { latency, inactive } = this.props.webrtc
    const screenBtnTitle = isFullscreen ? '退出全屏' : '全屏模式'
    const latencyTitle = `网络延迟${this.formatLatency(latency)}`
    const latencySrc = this.getSrc(latency)

    var resolutionIcon = document.getElementById('resolution-icon')
    var resolutionItem = document.getElementById('resolution-item')
    if (resolutionItem) {
      resolutionItem.onmouseover = resolutionIcon.onmouseover = function () {
        resolutionItem.style.display = 'block'
      }
      resolutionItem.onmouseout = resolutionIcon.onmouseout = function () {
        resolutionItem.style.display = 'none'
      }
    }

    return (
      <div>
        {inactive && <InactiveTipModal onClose={this.closeInactiveTipModal} />}
        {showPasteModal && <PasteModal onClose={this.closePasteModal} />}
        {showEndModal && <EndModal onClose={this.closeEndModal} />}
        {minimized && (
          <div id="mini-toolbar" onClick={this.restoreToolBar}>
            <svg
              t="1550223448316"
              height="20"
              width="20"
              fill="#fff"
              viewBox="0 0 1024 1024"
              version="1.1"
              p-id="1458">
              <path
                d="M543.2 224.5l403.9 409.6c17.1 17.4 17.1 45.8 0 63.2l-38.1 38.6c-17.1 17.4-45.1 17.4-62.3 0l-334.6-339.5-334.7 339.5c-17.1 17.4-45.2 17.4-62.3 0l-38.1-38.6c-17.1-17.4-17.1-45.8 0-63.2l403.9-409.6c17.2-17.3 45.2-17.3 62.3 0z"
                p-id="1459"
              />
            </svg>
          </div>
        )}
        {!minimized && (
          <div>
            <ul id="tool-bar">
              <li>
                <ul>
                  <li
                    className={`tool-button ${isFullscreen ? 'fullscreen' : ''
                      }`}
                    title={screenBtnTitle}
                    id="full-screen"
                    onClick={this.toggleFullscreen}>
                    <img className="full" src="static/full-screen.svg" />
                    <img
                      className="not-full"
                      src="static/exit-full-screen.svg"
                    />
                  </li>
                  <li
                    className="tool-button"
                    title="剪切板"
                    id="set-paste-text"
                    onClick={this.showPasteModal}>
                    <img src="static/paste.svg" />
                  </li>
                  <li
                    className="tool-button"
                    title="分辨率调整"
                    id="resolution-icon">
                    <img src="static/resolution.svg" />
                  </li>
                  <li
                    className="tool-button"
                    title="显示桌面"
                    id="show-desktop"
                    onClick={this.showDesktop}>
                    <img src="static/desktop.svg" />
                  </li>

                  {/* 删除「关闭」按钮 */}
                  {!isEditingMode && false && (
                    <li
                      className="tool-button"
                      title="关闭"
                      id="close-worktask"
                      onClick={this.showEndModal}>
                      <img src="static/shut_down.svg" />
                    </li>
                  )}
                </ul>
              </li>
              <li>
                <ul>
                  <li className="status-card fluence-slider" title="最大带宽占用">
                    <Slider
                      min={2}
                      max={30}
                      step={1}
                      defaultValue={this.maxBitrate}
                      handle={this.setMaxBitrate}
                      onChange={this.onMaxBitrateChange}
                    />
                  </li>
                  <li className="status-card latency" title={latencyTitle}>
                    <img id="latency" src={latencySrc} />
                  </li>
                  {false && (
                    <li className="status-card" title="已运行">
                      <img src="static/duration.svg" />
                      <span id="duration">{durationText}</span>
                    </li>
                  )}
                  <li className="status-card">
                    <img src="static/logo.svg" />
                  </li>
                  <li
                    className="tool-button"
                    id="hide-tb-btn"
                    onClick={this.hideToolBar}>
                    <svg
                      t="1550220491976"
                      width="20"
                      height="20"
                      fill="#fff"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      p-id="5481">
                      <path
                        d="M481 736l-403.9-409.6c-17.1-17.4-17.1-45.8 0-63.2l38.1-38.6c17.1-17.4 45.1-17.4 62.3 0l334.6 339.4 334.7-339.5c17.1-17.4 45.2-17.4 62.3 0l38.1 38.6c17.1 17.4 17.1 45.8 0 63.2l-404 409.7c-17.1 17.3-45.1 17.3-62.2 0z"
                        p-id="5482"
                      />
                    </svg>
                  </li>
                </ul>
              </li>
            </ul>
            <div id="resolution-item">
              <em></em>
              <i></i>
              <div>
                <span
                  id="resolution-auto-adjust"
                  onClick={e => this.onAutoAdjustResolution(e.currentTarget)}>
                  自适应窗口大小
                </span>
                <div className="resolution-classification">4:3</div>
                <span
                    id="resolution-800x600"
                    onClick={e =>
                        this.onSpecifySpecificResolution(
                            e.currentTarget,
                            800,
                            600
                        )
                    }>
                  800 * 600
                </span>
                <span
                  id="resolution-1024x768"
                  onClick={e =>
                    this.onSpecifySpecificResolution(
                      e.currentTarget,
                      1024,
                      768
                    )
                  }>
                  1024 * 768
                </span>
                <span
                  id="resolution-1280x1024"
                  onClick={e =>
                    this.onSpecifySpecificResolution(
                      e.currentTarget,
                      1280,
                      1024
                    )
                  }>
                  1280 * 1024
                </span>
                <div className="resolution-classification">16:9</div>
                <span
                  id="resolution-1600x900"
                  onClick={e =>
                    this.onSpecifySpecificResolution(
                      e.currentTarget,
                      1600,
                      900
                    )
                  }>
                  1600 * 900
                </span>
                <span
                  id="resolution-1920x1080"
                  onClick={e =>
                    this.onSpecifySpecificResolution(e.currentTarget, 1920, 1080)
                  }>
                  1920 * 1080
                </span>
                <span
                  id="resolution-2560x1440"
                  onClick={e =>
                    this.onSpecifySpecificResolution(e.currentTarget, 2560, 1440)
                  }>
                  2560 * 1440
                </span>
                <span
                  id="resolution-3840x2160"
                  onClick={e =>
                    this.onSpecifySpecificResolution(e.currentTarget, 3840, 2160)
                  }>
                  3840 * 2160
                </span>
                <div className="resolution-classification">16:10</div>
                <span
                  id="resolution-1680x1050"
                  onClick={e =>
                    this.onSpecifySpecificResolution(e.currentTarget, 1680, 1050)
                  }>
                  1680 * 1050
                </span>
                <span
                    id="resolution-1920x1200"
                    onClick={e =>
                        this.onSpecifySpecificResolution(e.currentTarget, 1920, 1200)
                    }>
                  1920 * 1200
                </span>
                {/*下面两个经过测试仍存在问题，暂不开放*/}
                <span
                    hidden="true"
                    id="resolution-2560x1600"
                    onClick={e =>
                        this.onSpecifySpecificResolution(e.currentTarget, 2560, 1600)
                    }>
                  2560 * 1600
                </span>
                <span
                    hidden="true"
                    id="resolution-3840x2400"
                    onClick={e =>
                        this.onSpecifySpecificResolution(e.currentTarget, 3840, 2400)
                    }>
                  3840 * 2400
                </span>
              </div>
            </div>
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
    }
  }
}
export default connect(mapStateToProps, mapDispatchToProps)(ToolBar)
