import * as React from 'react'
import styled from 'styled-components'

const ListWrapper = styled.div`
  display: flex;
  flex-wrap: wrap;

  .item {
    border: 1px solid #10398b;
    border-radius: 15px;
    margin-right: 12px;
    margin-bottom: 15px;
    background-color: white;
    padding: 0 25px;
    color: #10398b;
    height: 30px;
    font-size: 16px;
    line-height: 30px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
`

interface IProps {
  itemList: string[]
}

export default function RadiusItem({ itemList }: IProps) {
  return (
    <ListWrapper>
      {itemList?.map(i => (
        <span title={i} key={i} className='item'>
          {i}
        </span>
      ))}
    </ListWrapper>
  )
}
