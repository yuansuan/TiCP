import * as React from 'react'
import styled from 'styled-components'

interface IProps {
  align?: string
  width?: number
}

const LabelWrapper = styled.span<IProps>`
  display: inline-block;
  width: ${props => (props as any).width || '180'}px;
  text-align: ${props => (props as any).align || 'right'};
  .required {
    color: red;
  }
`

function Label({
  required = false,
  starBefore = true,
  align = 'right',
  width = 180,
  children
}) {
  return (
    <LabelWrapper align={align} width={width}>
      {required && starBefore ? <span className='required'> * </span> : ''}
      {children}
      {required && !starBefore ? <span className='required'> * </span> : ''}
    </LabelWrapper>
  )
}

export default Label
