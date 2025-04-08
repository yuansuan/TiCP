import * as React from 'react'
import FormStepsFooter, { IBtn } from '../FormSteps/FormStepsFooter'

export default function UserFormStepsFooter({
  isEdit,
  current,
  onOk,
  onCancel,
  onNext,
  onPre,
  onClose,
  onCopy
}) {
  let btns = [
    { label: '取消', display: true, fn: onCancel, type: 'link' },
    { label: '下一步', display: true, fn: onNext, type: 'primary' }
  ]

  if (current === 1) {
    btns = [
      { label: '上一步', display: true, fn: onPre, type: 'link' },
      {
        label: isEdit ? '编辑' : '创建',
        display: true,
        fn: onOk,
        type: 'primary'
      }
    ]
  }

  if (current === 2) {
    btns = [
      !isEdit && {
        label: '复制用户信息',
        display: false,
        fn: onCopy,
        type: 'link'
      },
      { label: '确认', display: true, fn: onClose, type: 'primary' }
    ]
  }

  return <FormStepsFooter btns={btns as IBtn[]} />
}
