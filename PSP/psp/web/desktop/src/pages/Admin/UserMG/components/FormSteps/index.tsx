import * as React from 'react'
import styled from 'styled-components'
import { ArrowRightOutlined } from '@ant-design/icons'
interface IStep {
  icon?: string
  name: string
  selected?: boolean
}

interface IFormStepsProps {
  steps: IStep[]
  current?: string
  width?: string | number
}

interface IWrapperProps {
  width: string | number
}

const FormStepsWrapper = styled.div<IWrapperProps>`
  display: flex;
  width: ${props => props.width};
  justify-content: space-around;
`

export default function FormSteps(props: IFormStepsProps) {
  const { steps, current, width } = props
  return (
    <FormStepsWrapper width={width || '100%'}>
      {steps.map(s => (
        <Step
          key={s.name}
          name={s.name}
          selected={current ? current === s.name : s.selected}
        />
      ))}
    </FormStepsWrapper>
  )
}

const StepWrapper = styled.div`
  color: #194e8b;
  .title {
    font-size: 14px;
    padding: 0 10px;
  }
`

function Step({ name, selected }) {
  return (
    <StepWrapper>
      {selected ? <ArrowRightOutlined /> : null}
      <span className='title'>{name}</span>
    </StepWrapper>
  )
}
