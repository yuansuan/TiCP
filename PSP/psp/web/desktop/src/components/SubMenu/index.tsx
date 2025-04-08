import * as React from 'react'
import { Link, withRouter } from 'react-router-dom'
import { computed } from 'mobx'
import styled from 'styled-components'
import { Icon } from '@/components'
import SubPage from './SubPage'

interface IProps {
  title?: string
  menus: {
    name: string
    path: string
    icon?: string
  }[]
  children: any
  history?: any
}

const Wrapper = styled.div`
  background: #fff;
  .flex-box {
    display: flex;
  }

  .bar {
    font-size: 16px;
    height: 48px;
    line-height: 48px;
    border-bottom: 1px solid #e8e8e8;
    padding-left: 50px;
    position: relative;
  }

  .nav {
    width: 150px;
    border-right: 1px solid #e8e8e8;

    .menu-item {
      display: block;
      height: 40px;
      line-height: 40px;
      padding-left: 20px;
      cursor: pointer;
      opacity: 0.65;
      font-size: 14px;
      color: #4a4a4a;

      &.active {
        background: #b3cdf5;
        color: #000000;

        .anticon {
          margin-right: 10px;
          color: #000000;
        }
      }

      .anticon {
        margin-right: 10px;
        color: #4a4a4a;
      }
    }
  }

  .content {
    flex: 1;
    min-height: calc(100vh - 230px);

    .ant-spin {
      width: 100%;
    }
  }
`

@(withRouter as any)
export default class SubMenu extends React.Component<IProps> {
  static SubPage = SubPage

  @computed
  get active() {
    return this.props.history.location.pathname
  }

  render() {
    const { title, menus, children } = this.props
    return (
      <Wrapper>
        <div className='bar'>{title}</div>
        <div className='flex-box'>
          <div className='nav'>
            {menus.map((menu, index) => (
              <Link
                to={menu.path}
                className={`menu-item ${
                  menu.path === this.active ? 'active' : ''
                }`}
                key={index}>
                {menu.icon && <Icon type={menu.icon} />}
                {menu.name}
              </Link>
            ))}
          </div>
          <div className='content'>{children}</div>
        </div>
      </Wrapper>
    )
  }
}
