import * as React from 'react'

import { StyledLink } from './style'

interface IProps {
  name: string
  title?: string
  disabled: boolean
  active?: boolean
  onClick: () => void
}

export default class Link extends React.Component<IProps> {
  render() {
    const { name, title, disabled, active, onClick } = this.props

    return (
      <StyledLink
        className={active ? 'active' : ''}
        title={title}
        {...(disabled ? {} : { onClick: () => onClick() })}>
        {name}
      </StyledLink>
    )
  }
}
