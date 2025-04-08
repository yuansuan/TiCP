/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { RouterType } from './typing'
import { getBreadCrumb } from './utils'

const StyledLayout = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  padding-left: 20px;

  > .breadcrumb {
    display: flex;
    list-style: none;
    margin: 0;
    padding-left: 4px;
    height: 32px;
    overflow: hidden;
    height: 100%;

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
  routers: RouterType[]
  history
}

export const Breadcrumb = observer(function Breadcrumb({
  routers,
  history,
}: Props) {
  const state = useLocalStore(() => ({
    breadcrumb: [],
    setBreadcrumb(arr) {
      this.breadcrumb = arr
    },
  }))

  const getName = (item: RouterType) => {
    if (typeof item.name === 'function') {
      return item.name()
    } else {
      return item.name
    }
  }

  useEffect(() => {
    if (history.listen) {
      state.setBreadcrumb(getBreadCrumb(history?.location?.pathname, routers))
      return history.listen(({ pathname }) => {
        state.setBreadcrumb(getBreadCrumb(pathname, routers))
      })
    }

    return undefined
  }, [])

  return (
    <StyledLayout>
      <ul className='breadcrumb'>
        {state.breadcrumb.map((item, index) => [
          index === 0 ? null : <li key={index * 2 + 1}>/</li>,
          item.path && index !== state.breadcrumb.length - 1 ? (
            <li
              key={index}
              className='link'
              onClick={() =>
                history.push &&
                history.push(item.path, history?.location?.state)
              }>
              {getName(item)}
            </li>
          ) : (
            <li key={index}>{getName(item)}</li>
          ),
        ])}
      </ul>
    </StyledLayout>
  )
})
