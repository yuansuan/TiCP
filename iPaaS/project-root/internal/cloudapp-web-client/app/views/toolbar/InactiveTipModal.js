import React from 'react'
import WorkTaskService from '@/services/WorkTaskService'

class InactiveTipModal extends React.Component {
  state = {
    countDown: 60,
  }
  componentDidMount() {
    this.timer = setInterval(() => {
      const cd = this.state.countDown
      if (cd <= 0) {
        this.exitApp()
      } else {
        this.setState({ countDown: cd - 1 })
      }
    }, 1000)
  }
  componentWillUnmount() {
    clearInterval(this.timer)
  }

  onClose = () => {
    this.props.onClose && this.props.onClose()
  }
  exitApp = () => {
    const search = new URLSearchParams(window.location.search)
    const userId = Number(search.get('user_id'))
    const workTaskId = Number(search.get('worktask_id'))
    WorkTaskService.closeWorkTask(userId, workTaskId).then(() => {
      setTimeout(() => {
        window.close()
      }, 200)
    })
  }

  render() {
    return (
      <div className="modal is-active" id="inactive-tip-confirm-modal">
        <div className="modal-background" />
        <div className="modal-content">
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">提醒</p>
              <span
                onClick={this.onClose}
                id="inactive-tip-modal-close"
                aria-label="close">
                <svg
                  t="1564628203999"
                  className="icon"
                  viewBox="0 0 1024 1024"
                  version="1.1"
                  xmlns="http://www.w3.org/2000/svg"
                  p-id="7464"
                  data-spm-anchor-id="a313x.7781069.0.i3"
                  width="200"
                  height="200">
                  <path
                    d="M510.65692138 571.33441162L246.2640686 837.06787109c-16.43417359 16.51657104-43.14495849 16.58331299-59.66152954 0.14996338-16.51657104-16.43334961-16.58331299-43.14495849-0.14996338-59.66235351l265.55136109-266.89691162-265.40057373-264.06573485c-16.51657104-16.43334961-16.58413697-43.14495849-0.15161133-59.66152956 16.43334961-16.51657104 43.14495849-16.58496095 59.66152954-0.15078735L513.34225463 452.66558838l264.39450074-265.73510742c16.43417359-16.51657104 43.14495849-16.58331299 59.66152955-0.14996338 16.51657104 16.43334961 16.58331299 43.14495849 0.14996337 59.66235351L571.99688721 513.34060669l265.40057373 264.06573487c16.51657104 16.43334961 16.58413697 43.14495849 0.15161133 59.66152953-16.43334961 16.51657104-43.14495849 16.58496095-59.66152954 0.15078735L510.65774537 571.33441162z"
                    fill="#5B5B5B"
                    p-id="7465"
                  />
                </svg>
              </span>
            </header>
            <section className="modal-card-body">
              长时间没有使用，应用将在<span className="count-down">
                {this.state.countDown}
              </span>秒后关闭！
            </section>
            <footer className="modal-card-foot">
              <button
                className="button is-success"
                onClick={this.onClose}
                id="close-worktask-confirm-modal-confirm">
                继续工作
              </button>
              <button
                className="button"
                onClick={this.exitApp}
                id="close-worktask-confirm-modal-cancel">
                退出应用
              </button>
            </footer>
          </div>
        </div>
      </div>
    )
  }
}

export default InactiveTipModal
