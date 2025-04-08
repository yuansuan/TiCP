/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Input, Checkbox } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { InputProps } from 'antd/lib/input'
import { StyledChooser } from './style'
import { useDidUpdate } from '@/utils/hooks'

const SearchInput = Input.Search

type Key = string | number

type Option = {
  [key: string]: any
  key: Key
  name: string
}

type Props = {
  items: Array<Option | [Key, string]>
  itemMapper?: (item: Option) => React.ReactNode
  style?: React.CSSProperties
  searchable?: boolean
  searchInputProps?: InputProps
  selectedKeys?: Key[]
  defaultSelectedKeys?: Key[]
  updateSelectedKeys?: (keys: Key[]) => void
  onChange?: (
    selectedKeys: Key[],
    info?: { node?: Option; checked: boolean }
  ) => void
  filter?: (item: Option, keywords: string) => boolean
}

export default observer(function Chooser(originProps: Props) {
  const state = useLocalStore(() => ({
    props: {
      items: originProps.items,
      selectedKeys: originProps.selectedKeys,
      searchable: originProps.searchable,
    },
    setProps(props) {
      this.props = props
    },
    _selectedKeys: originProps.defaultSelectedKeys || [],
    _setSelectedKeys(keys) {
      this._selectedKeys = keys
    },
    keywords: '',
    setKeywords(keywords) {
      this.keywords = keywords
    },
    setSelectedKeys(keys) {
      const { selectedKeys } = this.props

      if (selectedKeys) {
        selectedKeys && originProps.updateSelectedKeys(keys)
      } else {
        this._setSelectedKeys(keys)
      }
    },
    get selectedKeys() {
      return this.props.selectedKeys || this._selectedKeys
    },
    get visibleItems() {
      const items = this.props.items.map(item => {
        if (Array.isArray(item)) {
          return {
            key: item[0],
            name: item[1],
          }
        }
        return item
      })

      const { searchable = true } = state.props
      const { keywords } = this

      if (!searchable) {
        return items
      } else {
        const search =
          originProps.filter ||
          ((item: { key: Key; name: string }, keywords: string) => {
            return !keywords || item.name.includes(keywords)
          })

        return items.filter(item => search(item, keywords))
      }
    },
    get indeterminate() {
      return (
        !this.allSelected &&
        this.visibleItems.some(item => this.selectedKeys.includes(item.key))
      )
    },
    get allSelected() {
      return (
        this.visibleItems.length > 0 &&
        this.visibleItems.every(item => this.selectedKeys.includes(item.key))
      )
    },
  }))

  useDidUpdate(() => {
    state.setProps({
      items: originProps.items,
      selectedKeys: originProps.selectedKeys,
      searchable: originProps.searchable,
    })
  }, [originProps.items, originProps.selectedKeys, originProps.searchable])

  function selectAll(e) {
    const { checked } = e.target
    const { onChange } = originProps

    if (checked) {
      // select all visible items
      state.setSelectedKeys([
        ...new Set([
          ...state.selectedKeys,
          ...state.visibleItems.map(item => item.key),
        ]),
      ])
    } else {
      // cancel select all visible items
      const cancelKeys = state.visibleItems.map(item => item.key)
      state.setSelectedKeys([
        ...state.selectedKeys.filter(key => !cancelKeys.includes(key)),
      ])
    }

    // trigger onChange listener
    onChange && onChange([...state.selectedKeys])
  }

  function _onChange(checked, key) {
    const keys = state.selectedKeys

    if (checked) {
      keys.push(key)
    } else {
      const index = state.selectedKeys.findIndex(item => item === key)
      keys.splice(index, 1)
    }

    state.setSelectedKeys([...keys])
  }

  const {
    searchable = true,
    onChange,
    style,
    searchInputProps,
    itemMapper,
  } = originProps

  return (
    <StyledChooser style={{ minWidth: 150, ...style }}>
      {searchable ? (
        <SearchInput
          className='search'
          size='small'
          onChange={e => state.setKeywords(e.target.value)}
          allowClear
          {...searchInputProps}
        />
      ) : null}

      <div className='allSelector'>
        <Checkbox
          checked={state.allSelected}
          indeterminate={state.indeterminate}
          onChange={selectAll}>
          全选
        </Checkbox>
      </div>

      <div className='list'>
        {state.visibleItems.map(item => (
          <div key={item.key} className='item'>
            <Checkbox
              checked={state.selectedKeys.includes(item.key)}
              onChange={e => {
                const { checked } = e.target
                _onChange(checked, item.key)
                onChange &&
                  onChange([...state.selectedKeys], { node: item, checked })
              }}
            />
            {itemMapper ? (
              itemMapper(item)
            ) : (
              <span className='name' title={item.name}>
                {item.name}
              </span>
            )}
          </div>
        ))}
      </div>
    </StyledChooser>
  )
})
