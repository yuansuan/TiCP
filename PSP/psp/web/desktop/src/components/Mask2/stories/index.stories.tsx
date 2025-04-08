/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Mask from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Mask',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Mask,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Loading } from './Loading'
export { Theme } from './Theme'
