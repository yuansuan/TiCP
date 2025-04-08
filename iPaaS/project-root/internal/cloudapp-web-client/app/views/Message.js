import React from 'react';

import { connect } from 'react-redux'
class Message extends React.Component {
  render() {
    const {type, text} = this.props.message
    const cls = `app-message ${type}`
    if(!text){
      return null
    }
    return (
      <div className={cls}>
          {text}
      </div>
      )
  }
}
function mapStateToProps(state) {
  return {
    message: state.message
  }
}
export default  connect(mapStateToProps)(Message)