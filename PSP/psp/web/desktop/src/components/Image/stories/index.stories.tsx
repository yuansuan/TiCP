/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Image from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Image',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Image,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Empty } from './Empty'
export { NotFound } from './NotFound'
