import React from 'react'
import { connect } from 'react-redux'

class PasteModal extends React.Component {
  componentDidMount() {}

  onClose = () => {
    this.props.onClose && this.props.onClose()
  }
  onPasteIn = e => {
    const pasteText = e.clipboardData.getData('text/plain')
    this.sendPasteText(pasteText)

    setTimeout(() => {
      this.refs.pasteTextArea.value = ''
      this.onClose()
    }, 500)
  }
  onPasteKeyUp = e => {
    if (e.keyCode == 13) {
      const pasteText = this.refs.pasteTextArea.value
      this.sendPasteText(pasteText)
      setTimeout(() => {
        this.onClose()
      }, 500)
      this.refs.pasteTextArea.value = ''
    }
  }
  onChangeCopyOutText = e => {
    this.props.copyOutText(e.target.value)
  }
  sendPasteText = text => {
    const pline = window.webrtcPipeline
    pline && pline.sendPasteText(text)
  }
  render() {
    const { copyOutText } = this.props.webrtc
    return (
      <div className="modal is-active" id="set-paste-text-modal">
        <div className="modal-background" />
        <div className="modal-content">
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">剪切板</p>
              <span
                id="set-paste-text-modal-close"
                aria-label="close"
                onClick={this.onClose}>
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
              <p>拷贝到工作站</p>
              <textarea
                className="textarea"
                id="set-paste-textarea"
                ref="pasteTextArea"
                onPaste={this.onPasteIn}
                onKeyUp={this.onPasteKeyUp}
                placeholder="拷贝到工作站"
              />
              <p>拷贝到本地电脑</p>
              <textarea
                className="textarea"
                id="set-copy-textarea"
                value={copyOutText}
                ref="copyOutTextArea"
                onChange={this.onChangeCopyOutText}
                placeholder="拷贝到本地电脑"
              />
            </section>
          </div>
        </div>
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
    copyOutText: text => {
      dispatch({ type: 'COPY_OUT_TEXT', text })
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(PasteModal)
