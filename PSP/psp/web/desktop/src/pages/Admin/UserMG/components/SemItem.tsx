import * as React from 'react'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import styled from 'styled-components'

const ListWrapper = styled.div`
  display: flex;

  .content {
    width: 95%;
    margin-left: 6px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
`

interface IProps {
  itemList: string[]
}

@observer
export default class SemItem extends React.Component<IProps> {
  @computed
  get titleContent() {
    return this.props.itemList.join('; ')
  }

  render() {
    return (
      <ListWrapper>
        <span title={this.titleContent} className='content'>
          {this.titleContent}
        </span>
      </ListWrapper>
    )
  }
}
