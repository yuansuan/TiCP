/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import styled from 'styled-components'
import { Story, Meta } from '@storybook/react/types-6-0'
import mdx from './doc.mdx'
import Uploader, { UploaderProps } from '..'
import { Button } from '../..'

export default {
  title: 'components/Uploader',
  parameters: {
    docs: {
      page: mdx,
    },
  },
} as Meta

const StyledLayout = styled.div`
  padding: 10px 0;
`

const Template: Story<UploaderProps> = props => {
  const uploader = useMemo(() => new Uploader(props), [])

  return (
    <StyledLayout style={props['layoutStyle']}>
      <Button
        type='primary'
        {...props}
        onClick={() => uploader.upload({ origin: 'uploader-example' })}>
        上传
      </Button>
    </StyledLayout>
  )
}

export const upload = Template.bind({})
