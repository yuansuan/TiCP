/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { Select } from 'antd'

type SelectorProps = {
  loading: boolean
  filters: { key: string; name: string }[]
  onChange: (values: any) => void
  onDropdownVisibleChange?: (values: any) => void
}
export const Selector = observer(function Selector({
  filters,
  ...props
}: SelectorProps) {
  return (
    <Select
      {...props}
      style={{
        width: 220
      }}
      mode='multiple'
      allowClear
      placeholder='全部'
      filterOption={(input, option) =>
        option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
      }>
      {filters?.map(item => (
        <Select.Option key={item.key} value={item.key}>
          {item.name}
        </Select.Option>
      ))}
    </Select>
  )
})
