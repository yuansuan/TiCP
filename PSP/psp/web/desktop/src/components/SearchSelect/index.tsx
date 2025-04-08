/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { Select } from 'antd'
import { SelectProps } from 'antd/es/select'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useDidUpdate } from '@/utils/hooks'

type Option = {
  key: string | number
  name: string
  [key: string]: any
}

type IProps = {
  options: Array<Option | [string | number, string]>
  caseSensitive?: boolean
  filter?: (key: string, item: Option) => boolean
  renderOptions?: (options: Array<Option>) => React.ReactNode
} & Omit<SelectProps<string | string[] | number | number[]>, 'options'>

export default observer(
  function SearchSelect(originalProps: IProps, ref: any) {
    const state = useLocalStore(() => ({
      props: {
        options: originalProps.options,
        caseSensitive: originalProps.caseSensitive
      },
      setProps(props) {
        this.props = props
      },
      searchKey: '',
      setSearchKey(key) {
        this.searchKey = key
      },
      get options() {
        const opts = this.props.options.map(item => {
          if (Array.isArray(item)) {
            return {
              key: item[0],
              name: item[1]
            }
          }

          return item
        })

        if (!this.searchKey) {
          return opts
        }

        const { caseSensitive = true } = this.props

        return opts.filter(
          filter
            ? item => originalProps.filter(this.searchKey, item)
            : item => {
                if (caseSensitive) {
                  return item.name.includes(this.searchKey)
                } else {
                  return item.name
                    .toLowerCase()
                    .includes(this.searchKey.toLowerCase())
                }
              }
        )
      }
    }))

    useDidUpdate(() => {
      state.setProps({
        options: originalProps.options,
        caseSensitive: originalProps.caseSensitive
      })
    }, [originalProps.caseSensitive, originalProps.options])

    function _onSelect(value, option) {
      const { onSelect } = originalProps
      onSelect && onSelect(value, option)

      state.setSearchKey('')
    }

    const {
      options,
      onSelect,
      filter,
      renderOptions,
      caseSensitive,
      ...rest
    } = originalProps

    const Options = useMemo(
      () =>
        renderOptions
          ? renderOptions(state.options)
          : state.options.map(item => (
              <Select.Option key={item.key} value={item.key}>
                {item.name}
              </Select.Option>
            )),
      [state.options]
    )

    return (
      <Select
        ref={ref}
        showSearch
        filterOption={false}
        onSearch={key => state.setSearchKey(key)}
        onSelect={_onSelect}
        {...rest}>
        {originalProps.children}
        {Options}
      </Select>
    )
  },
  {
    forwardRef: true
  }
)
