/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { EditableText } from '@/components'
import { EditableTextProps } from '@/components/EditableText'

const numberValidate = (value: string) => {
  if (value === '') return true
  if (+(+value).toFixed(2) === 0) return false
  const flag = /^[+]{0,1}(\d+)$|^[+]{0,1}(\d+\.\d+)$/.test(value) // 匹配 0、正浮点数、正整数
  if (!flag) {
    return '只允许输入大于 0 的数'
  } else {
    return true
  }
}

type IProps = EditableTextProps & {
  value: number
  setValue: (arg: number) => void
  customHelp?: string
  customBeforeConfirm?: (value: string) => string | boolean
  style?: React.CSSProperties
}

const IEditableText = ({
  value,
  setValue,
  customHelp,
  customBeforeConfirm,
  style,
  ...props
}: IProps) => {
  const model = EditableText.useModel()

  React.useEffect(() => {
    model.value = value?.toString() || ''
  }, [value])

  return (
    <EditableText
      style={
        style || {
          width: 180,
          display: 'inline-block',
          height: 24,
          minHeight: 24,
        }
      }
      beforeConfirm={customBeforeConfirm || numberValidate}
      beforeCancel={() => {
        model.value = value?.toString() || ''
        return true
      }}
      model={model}
      onConfirm={(value: string) => {
        let newValue: number
        if (value === '') {
          newValue = null
        } else if (value !== '') {
          newValue = +(+value).toFixed(2)
        }
        setValue(newValue)
        model.value = newValue?.toString() || ''
      }}
      {...props}
    />
  )
}

export default IEditableText
