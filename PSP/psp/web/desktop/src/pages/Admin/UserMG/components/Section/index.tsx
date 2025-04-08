import * as React from 'react'

import { SectionWrapper } from './style'

interface IProps {
  className?: string
  title?: string | React.ReactNode
  icon?: React.ReactNode
}

export default class Section extends React.Component<IProps> {
  render() {
    const { icon, title, className, children } = this.props

    return (
      <SectionWrapper className={className}>
        <div className='header'>
          {icon ? <span className='icon'>{icon}</span> : null}
          {title}
        </div>
        <div className='body'>{children}</div>
      </SectionWrapper>
    )
  }
}
