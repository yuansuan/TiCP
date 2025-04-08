/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Scrollbars } from 'react-custom-scrollbars'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const StyledLayout = styled.div`
  border: 1px solid ${({ theme }) => theme.borderColorBase};
  border-radius: 2px;
  height: 100%;
  margin-right: 10px;
  display: flex;

  > .root {
    padding: 0;
    width: 68px;
    text-align: center;
    background-color: ${({ theme }) => theme.backgroundColorBase};
    cursor: pointer;
    color: rgba(0, 0, 0, 0.69);

    > .anticon {
      margin-right: 4px;
    }
  }

  .list {
    display: flex;
    list-style: none;
    margin: 0;
    padding-left: 4px;
    height: 32px;

    > li {
      margin: 0 2px;

      &.link {
        cursor: pointer;
        color: gray;

        &:hover {
          color: ${({ theme }) => theme.primaryColor};
        }
      }
    }
  }
`

type Props = {
  width: number
}

export const Bar = observer(function Bar({ width }: Props) {
  const store = useStore()
  const { history, dirTree } = store
  const { current } = history
  const path = current?.path || ''
  const list = path.split('/').filter(item => !!item)

  function jump(path) {
    const node = dirTree.filterFirstNode(item => item.path === path)
    if (node) {
      store.setNodeId(node.id)
    }
  }

  return (
    <StyledLayout>
      <div className='root' onClick={() => jump('/')}>
        {dirTree.children[0]?.name}
      </div>
      <Scrollbars autoHeight style={{ width: width - 68 }}>
        <ul className='list'>
          {list.map((item, index) => [
            <li key={index * 2 + 1}>/</li>,
            index === list.length - 1 ? (
              <li key={index}>{item}</li>
            ) : (
              <li
                key={index}
                className='link'
                onClick={() => jump(list.slice(0, index + 1).join('/'))}>
                {item}
              </li>
            ),
          ])}
        </ul>
      </Scrollbars>
    </StyledLayout>
  )
})
