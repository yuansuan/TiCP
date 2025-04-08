/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect } from 'react'
import * as ReactDOM from 'react-dom'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Tooltip } from 'antd'
import { debounce } from 'lodash'

import Icon from '../Icon'
import { StyledPartialList } from './style'

type Item = string | number

interface Props {
  items: Array<Item>
  itemMapper?: (item: Item, index) => Item | React.ReactNode
  additionalItemMapper?: (item: Item, index) => Item | React.ReactNode
  maxWidth?: number
  style?: React.CSSProperties
  Additional?: React.ComponentType<{ visibleIndex?: number }>
}

function computeWidth(element) {
  const style = element.currentStyle || window.getComputedStyle(element),
    width = element.offsetWidth, // or use style.width
    margin = parseFloat(style.marginLeft) + parseFloat(style.marginRight),
    padding = parseFloat(style.paddingLeft) + parseFloat(style.paddingRight),
    border =
      parseFloat(style.borderLeftWidth) + parseFloat(style.borderRightWidth)

  return width + margin + padding + border
}

export default observer(function PartialList(originProps: Props) {
  const state = useLocalStore(() => ({
    props: originProps,
    setProps(props) {
      this.props = props
    },
    visibleIndex: 0,
    // use computing to avoid flicker
    computing: false,
    setVisibleIndex(index) {
      this.visibleIndex = index
    },
    setComputing(computing) {
      this.computing = computing
    },
    get visibleItems() {
      return this.props.items.slice(0, Math.max(this.visibleIndex, 0))
    },
    get additionalItems() {
      return this.props.items.slice(Math.max(this.visibleIndex, 0))
    },
    get showAdditional() {
      return this.visibleIndex < this.props.items.length
    },
    get autoResizable() {
      return (
        this.props.maxWidth === undefined || isNaN(Number(this.props.maxWidth))
      )
    },
  }))

  useEffect(() => {
    state.setProps(originProps)
  }, [originProps])

  const containerRef = useRef(null)
  const listRef = useRef(null)
  const additionalRef = useRef(null)

  const getVisibleIndex = debounce(
    () => {
      // render all items to compute width
      const { items } = state.props
      state.setComputing(true)
      state.setVisibleIndex(items.length)

      let { maxWidth } = state.props
      if (state.autoResizable) {
        const container = ReactDOM.findDOMNode(containerRef.current) as any
        const additional = ReactDOM.findDOMNode(additionalRef.current) as any

        maxWidth = container.clientWidth - computeWidth(additional)
      }
      maxWidth = Number(maxWidth)

      // wait list render
      setTimeout(() => {
        state.setComputing(false)

        let total = 0
        const list = ReactDOM.findDOMNode(listRef.current) as any
        ;[...list.children].every((item, index) => {
          total += computeWidth(item)
          if (total > maxWidth) {
            state.setVisibleIndex(index)
            return false
          }
          return true
        })
      }, 0)
    },
    300,
    { leading: true }
  )

  useEffect(() => {
    state.setVisibleIndex(state.props.items.length)
    if (state.autoResizable) {
      window.addEventListener('resize', getVisibleIndex)
    }

    return () => {
      if (state.autoResizable) {
        window.removeEventListener('resize', getVisibleIndex)
      }
    }
  }, [state.autoResizable])

  useEffect(() => {
    getVisibleIndex()
  }, [])

  useEffect(() => {
    if (!state.autoResizable) {
      getVisibleIndex()
    }
  }, [state.props.maxWidth])

  const { style, Additional, itemMapper, additionalItemMapper } = state.props

  return (
    <StyledPartialList
      ref={containerRef}
      style={style}
      onClick={e => e.stopPropagation()}>
      <div
        className='list'
        ref={listRef}
        style={state.computing ? { visibility: 'hidden' } : {}}>
        {state.visibleItems.map(
          itemMapper ||
            ((item, index) => (
              <span key={index}>
                {item}
                {index === state.visibleItems.length - 1 ? null : '；'}
              </span>
            ))
        )}
      </div>
      <div
        ref={additionalRef}
        style={!state.showAdditional ? { visibility: 'hidden' } : {}}
        className='additional'>
        {Additional ? (
          <Additional visibleIndex={state.visibleIndex} />
        ) : (
          <Tooltip
            title={state.additionalItems.map(
              additionalItemMapper ||
                ((item, index) => (
                  <span key={index}>
                    {item}
                    {index === state.additionalItems.length - 1 ? null : '；'}
                  </span>
                ))
            )}>
            <Icon type='more' />
          </Tooltip>
        )}
      </div>
    </StyledPartialList>
  )
})
