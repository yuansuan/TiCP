/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import SearchSelect from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/SearchSelect',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: SearchSelect,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { CaseSensitive } from './CaseSensitive'
export { CustomFilter } from './CustomFilter'
export { RenderOptions } from './RenderOptions'
