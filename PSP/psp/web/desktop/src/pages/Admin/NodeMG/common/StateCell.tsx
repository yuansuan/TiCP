import * as React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react'

const Wrapper = styled.div`
  margin-right: 5px;
  height: 100%;
  display: flex;
  justify-content: flex-start;
  align-items: center;

  .point {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: ${props => props.color};
    margin-right: 10px;
  }

  .state {
    margin-right: 5px;
  }
`

interface IProps {
  state: any
  color: string
  children?: React.ReactNode
}

function StateCell(props: IProps) {
  const { color, state, children } = props
  return (
    <Wrapper color={color}>
      <div className='point' />
      <div className='state'>{state}</div>
      {children ? <div>{children}</div> : null}
    </Wrapper>
  )
}

export default observer(StateCell)
