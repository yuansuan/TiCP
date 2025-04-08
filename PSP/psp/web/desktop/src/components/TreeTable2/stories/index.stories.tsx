/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import TreeTable from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/TreeTable',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: TreeTable,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
