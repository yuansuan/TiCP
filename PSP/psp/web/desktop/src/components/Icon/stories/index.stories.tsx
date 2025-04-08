/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Icon from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Icon',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Icon,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Custom } from './Custom'
