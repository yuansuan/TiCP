import * as React from 'react'
import { Link, withRouter } from 'react-router-dom'
import styled from 'styled-components'
import { Breadcrumb } from 'antd'

interface IProps {
  title?: string
  children: any
  breadcrumbs?: { name: string; path?: string }[]
  history?: any
}

const Wrapper = styled.div`
  .flex-box {
    display: flex;
  }

  .bar {
    font-size: 16px;
    height: 48px;
    line-height: 48px;
    border-bottom: 1px solid ${props => props.theme.borderColor};
    padding-left: 40px;
    position: relative;
    display: flex;
    align-items: center;
  }
`

@(withRouter as any)
export default class SubMenu extends React.Component<IProps> {
  render() {
    const { title, children, breadcrumbs } = this.props
    return (
      <Wrapper>
        <div className='bar'>
          <div title={title} onClick={() => this.props.history.go(-1)} />
          {title}
        </div>
        <Breadcrumb>
          {breadcrumbs &&
            breadcrumbs.map(br => (
              <Breadcrumb.Item key={br.name}>
                {br.path ? <Link to='/monitor/system'>系统监控</Link> : br.name}
              </Breadcrumb.Item>
            ))}

          <Breadcrumb.Item>{name}</Breadcrumb.Item>
        </Breadcrumb>
        <div className='detail-content'>{children}</div>
      </Wrapper>
    )
  }
}
