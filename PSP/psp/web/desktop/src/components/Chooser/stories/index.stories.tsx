/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Chooser from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Chooser',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Chooser,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Controlled } from './Controlled'
export { Custom } from './Custom'
export { Filter } from './Filter'
