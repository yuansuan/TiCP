/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import RichTextEditor from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/RichTextEditor',
  decorators: [storyFn => <div>{storyFn()}</div>],
  component: RichTextEditor,
  parameters: {
    docs: {
      page: mdx,
    },
  },
}

export { Basic } from './Basic'
export { Controlled } from './Controlled'
export { MaxHeight } from './MaxHeight'
