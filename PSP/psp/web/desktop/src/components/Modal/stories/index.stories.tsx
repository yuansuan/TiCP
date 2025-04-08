/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Modal from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Modal',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: Modal,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Async } from './Async'
export { Controlled } from './Controlled'
export { Custom } from './Custom'
export { Data } from './Data'
export { Theme } from './Theme'
