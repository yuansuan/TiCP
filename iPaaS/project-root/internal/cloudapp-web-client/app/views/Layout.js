import React from 'react'
import UAParser from 'ua-parser-js'
import MainPage from '@/views/MainPage'
import Message from '@/views/Message'
import ErrorPage from '@/views/ErrorPage'
class Layout extends React.Component {
  componentDidMount() {}

  render() {
    const browser = new UAParser().getBrowser()
    if (
      ['Chrome', 'Edge'].indexOf(browser.name) === -1 
    ) {
      return <ErrorPage />
    }
    return (
      <div className="app-layout">
        <MainPage />
        <Message />
      </div>
    )
  }
}

export default Layout
