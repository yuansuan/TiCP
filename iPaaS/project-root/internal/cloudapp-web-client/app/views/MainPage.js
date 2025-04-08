import React from 'react'

import WorkSpace from '@/views/workspace/WorkSpace'
import ToolBar from '@/views/toolbar/ToolBar'
export default class MainPage extends React.Component {
  render() {
    return (
      <div id="container">
        <WorkSpace />
        <ToolBar />
      </div>
    )
  }
}
