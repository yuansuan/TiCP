import * as React from 'react'
import { Button } from '@/components'
import styled from 'styled-components'

const FooterWrapper = styled.div`
  position: absolute;
  display: flex;
  bottom: 0px;
  right: 0;
  width: 100%;
  line-height: 70px;
  height: 70px;
  background: white;

  .footerMain {
    margin-left: auto;

    button {
      width: 120px;
      height: 40px;
      margin: 0 20px;
    }
  }
`

interface IProps {
  onCancel: () => void
  onOk: () => void
  processing: boolean
}

export default function PendingFooter({ onCancel, onOk, processing }: IProps) {
  return (
    <FooterWrapper>
      <div className='footerMain'>
        <Button disabled={processing} type='primary' onClick={onOk}>
          {processing ? '请求中...' : '确认'}
        </Button>
        <Button onClick={onCancel}>取消</Button>
      </div>
    </FooterWrapper>
  )
}
