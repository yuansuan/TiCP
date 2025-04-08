import * as React from 'react'
import { FormSteps } from '..'

const steps = [
  {
    name: '填写用户信息'
  },
  {
    name: '角色'
  }
]

export default function UserFormSteps({ current }) {
  return (
    <div style={{ paddingBottom: 20 }}>
      <FormSteps steps={steps} current={steps[current].name} width='400px' />
    </div>
  )
}
