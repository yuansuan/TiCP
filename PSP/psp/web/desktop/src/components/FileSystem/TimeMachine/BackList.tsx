import * as React from 'react'

import Link from './Link'
import { StyledBackList } from './style'

interface IProps {
  onClick: (link: any) => void
  links: any[]
}

export default class BackList extends React.Component<IProps> {
  render() {
    const { links, onClick } = this.props

    return (
      <StyledBackList>
        <ul>
          {links.map((link, index) => (
            <li key={index}>
              <Link
                name={link.name}
                title={link.path}
                disabled={link.disabled}
                onClick={() => onClick(link)}
              />
            </li>
          ))}
        </ul>
      </StyledBackList>
    )
  }
}
