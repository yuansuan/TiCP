import React from 'react'

export default class ErrorPage extends React.Component {
  render() {
    return (
      <div className="error-page">
        <div>
          <span>不支持当前浏览器，请使用</span>
          <a href="https://www.google.com/intl/zh-CN/chrome/">Chrome</a>
          <span>(推荐)或者</span>
          <a href="https://www.microsoft.com/en-us/edge">Edge</a>
          <span>浏览器</span>
        </div>
      </div>
    )
  }
}
