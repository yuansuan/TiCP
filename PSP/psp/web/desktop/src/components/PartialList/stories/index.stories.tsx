/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import PartialList from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/PartialList',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: PartialList,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Controlled } from './Controlled'
export { CustomAdditional } from './CustomAdditional'
export { CustomAdditionalItem } from './CustomAdditionalItem'
export { CustomItem } from './CustomItem'
