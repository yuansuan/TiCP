/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import mdx from './doc.mdx'
import { Story, Meta } from '@storybook/react/types-6-0'
import { Tooltip as AntTooltip, message } from 'antd'
import EditableText, { EditableTextProps } from '..'
import { Icon, Button } from '../..'
import { useObserver } from 'mobx-react-lite'

export default {
  title: 'components/EditableText',
  component: EditableText,
  parameters: {
    docs: {
      page: mdx,
    },
  },
  argTypes: {
    EditIcon: { type: 'json' },
  },
} as Meta

const Template: Story<EditableTextProps> = function Template(props) {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      defaultValue={value}
      onConfirm={setValue}
      {...props}
    />
  )
}

export const Basic = Template.bind({})

export const DefaultEditing = Template.bind({})
DefaultEditing.args = {
  defaultEditing: true,
}

export const Hover = Template.bind({})
Hover.args = {
  defaultShowEdit: false,
}

export const Unit = Template.bind({})
Unit.args = {
  unit: 'MB/s',
}

export const Help = Template.bind({})
Help.args = {
  help: '帮助信息',
}

export const Filter = Template.bind({})
Filter.args = {
  filter: /[^a-zA-Z0-9]/g,
  help: '只允许输入字母和数字',
}

export function CustomEditIcon(props) {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      defaultValue={value}
      onConfirm={setValue}
      EditIcon={
        <AntTooltip title='重命名'>
          <Icon type='rename' />
        </AntTooltip>
      }
      {...props}
    />
  )
}

export const CustomText = function CustomText(props) {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      defaultValue={value}
      onConfirm={setValue}
      Text={value => (
        <AntTooltip title={value}>
          <span>{value}</span>
        </AntTooltip>
      )}
      {...props}
    />
  )
}

export const ErrorHandler: Story<EditableTextProps> = function ErrorHandler(
  props
) {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      defaultValue={value}
      onConfirm={setValue}
      defaultEditing={true}
      beforeConfirm={value =>
        value === 'hello world' ? true : Promise.reject('test error')
      }
      {...props}
    />
  )
}

export const Trigger = function Trigger() {
  const model = EditableText.useModel({
    defaultValue: 'hello world',
    defaultEditing: false,
  })

  return useObserver(() => (
    <>
      <EditableText
        style={{ width: 200, margin: 20, display: 'inline-block' }}
        showEdit={false}
        model={model}
      />
      {!model.editing && (
        <Button type='link' onClick={() => model.setEditing(true)}>
          编辑
        </Button>
      )}
    </>
  ))
}

export function Validation() {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      beforeConfirm={value => {
        const flag = /^[a-zA-Z0-9]*$/.test(value)

        if (!flag) {
          message.error('只允许输入字母和数字')
        }

        return flag
      }}
      help='只允许输入字母和数字'
      defaultEditing={true}
      defaultValue={value}
      onConfirm={setValue}
    />
  )
}
