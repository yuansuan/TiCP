import * as React from 'react'
import { Button } from 'antd'
import styled from 'styled-components'

const FooterWrapper = styled.div`
  position: absolute;
  display: flex;
  bottom: 0px;
  right: 0;
  width: 100%;
  padding: 20px;
  background: white;

  .footerMain {
    margin-left: auto;
    display: flex;

    button {
      min-width: 80px;
      height: 30px;
      margin: 0 10px;
    }
  }
`

export interface IBtn {
  fn: () => void
  display: boolean
  type: 'link' | 'default' | 'ghost' | 'primary' | 'dashed' | 'danger'
  label: string
}

interface IProps {
  btns: IBtn[]
}

export default function FormStepsFooter({ btns }: IProps) {
  return (
    <FooterWrapper>
      <div className='footerMain'>
        {btns.map(({ display, type, label, fn }, index) => (
          <Button
            key={index}
            type={type}
            style={{ display: display ? 'block' : 'none' }}
            onClick={fn}>
            {label}
          </Button>
        ))}
      </div>
    </FooterWrapper>
  )
}
