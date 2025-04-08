import * as React from 'react'
import styled from 'styled-components'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { Icon } from '@/components'

const Wrapper = styled.div`
  h3 {
    margin: 0;
    color: #3a8dff;
    font-size: 12px;
    margin-top: 10px;
    margin-bottom: 6px;
    display: flex;
    align-items: center;
    cursor: pointer;

    .icon {
      margin-right: 8px;
      font-size: 14px;
    }

    .icon.caret {
      font-size: 12px;
      margin-left: 5px;
    }
  }

  .section-content {
    border: 1px solid #eee;
    padding: 10px;
    display: flex;
    flex-wrap: wrap;

    .item {
      width: 25%;
      padding: 3px;
    }
  }
`

const Children: any = styled.div`
  display: ${(props: any) => (props.collapse ? 'none' : 'block')};
`

interface IProps {
  title: string
  icon: string
}

@observer
export default class Section extends React.Component<IProps> {
  @observable collapse = false

  render() {
    const { icon, title, children } = this.props
    return (
      <Wrapper>
        <h3 onClick={() => (this.collapse = !this.collapse)}>
          <Icon type={icon} />
          {title}
          <Icon
            className='caret'
            type={this.collapse ? 'caret-down' : 'caret-up'}
          />
        </h3>
        <Children collapse={this.collapse}>{children}</Children>
      </Wrapper>
    )
  }
}
