/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Story, Meta } from '@storybook/react/types-6-0'
import Button, { ButtonProps } from '..'
import mdx from './doc.mdx'

export default {
  title: 'components/Button',
  component: Button,
  parameters: {
    docs: {
      page: mdx,
    },
  },
} as Meta

const StyledLayout = styled.div`
  padding: 10px 0;

  > * {
    margin: 5px;
  }
`
const Template: Story<ButtonProps> = props => (
  <StyledLayout style={props['layoutStyle']}>
    <Button type='primary' {...props}>
      主按钮
    </Button>
    <Button type='secondary' {...props}>
      次按钮
    </Button>
    <Button {...props}>默认按钮</Button>
    <Button type='dashed' {...props}>
      虚线按钮
    </Button>
    <Button danger {...props}>
      危险按钮
    </Button>
    <Button type='link' {...props}>
      链接按钮
    </Button>
    <Button type='cancel' {...props}>
      取消按钮
    </Button>
  </StyledLayout>
)

export const Type = Template.bind({})
Type.args = {}

export const Block = Template.bind({})
Block.args = {
  block: true,
}

export const Disabled = Template.bind({})
Disabled.args = {
  disabled: 'disabled tip',
}

export const Ghost = Template.bind({})
Ghost.args = {
  ghost: true,
  layoutStyle: {
    backgroundColor: 'pink',
  },
} as ButtonProps

export const Loading = Template.bind({})
Loading.args = {
  loading: true,
}

export const Size = Template.bind({})
Size.args = {
  size: 'large',
}

export const AsyncClick = Template.bind({})
AsyncClick.args = {
  onClick: () =>
    new Promise<void>((resolve, reject) => {
      setTimeout(() => {
        resolve()
      }, 3000)
    }),
}

export const ButtonGroup: Story<void> = () => (
  <StyledLayout>
    <Button.Group>
      <Button>Cancel</Button>
      <Button>OK</Button>
    </Button.Group>
    <Button.Group>
      <Button disabled>L</Button>
      <Button disabled>M</Button>
      <Button disabled>R</Button>
    </Button.Group>
    <Button.Group>
      <Button>L</Button>
      <Button>M</Button>
      <Button>R</Button>
    </Button.Group>

    <h4>With Icon</h4>
    <Button.Group>
      <Button type='primary' icon='define'>
        confirm
      </Button>
      <Button type='primary' icon='cancel'>
        cancel
      </Button>
    </Button.Group>
    <Button.Group>
      <Button type='primary' icon='define' />
      <Button type='primary' icon='cancel' />
    </Button.Group>
  </StyledLayout>
)
