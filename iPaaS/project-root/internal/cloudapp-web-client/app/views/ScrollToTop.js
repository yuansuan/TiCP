import React from 'react';

class ScrollToTop extends React.Component {
  componentDidUpdate(prevProps) {
    if (this.props.history.action =='PUSH' && this.props.location !== prevProps.location) {
      window.scrollTo(0, 0);
    }
  }

  render() {
    return this.props.children
  }
}

export default ScrollToTop