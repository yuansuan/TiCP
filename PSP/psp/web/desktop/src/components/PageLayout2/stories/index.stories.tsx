/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import PageLayout from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/PageLayout',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: PageLayout,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { State } from './State'
export { Footer } from './Footer'
