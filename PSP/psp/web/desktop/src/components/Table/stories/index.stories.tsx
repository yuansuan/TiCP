/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Table from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Table',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Table,
  parameters: {
    docs: {
      page: mdx
    }
  }
}

export { Basic } from './Basic'
export { Colspan } from './Colspan'
export { Customizable } from './Customizable'
export { Filterable } from './Filterable'
export { FixedColumn } from './FixedColumn'
export { PercentWidth } from './PercentWidth'
export { Resizable } from './Resizable'
export { RowEvent } from './RowEvent'
export { Selectable } from './Selectable'
export { Sortable } from './Sortable'
export { Tree } from './Tree'
export { Expanded } from './Expanded'
export { Virtualized } from './Virtualized'
