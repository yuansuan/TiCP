import { Icon } from '@/components'
import * as React from 'react'
import styled from 'styled-components'

const Wrapper = styled.div`
  position: relative;
  margin-bottom: 10px;
  min-height: 32px;

  .title {
    position: absolute;
    left: 20px;
    top: -16px;
    background: #fff;
    padding: 0 10px;
    line-height: 32px;
    font-size: 16px;
    margin-bottom: 6px;
    cursor: pointer;

    .text {
      display: inline-block;
      min-width: 70px;
    }

    .icon {
      padding-left: 10px;
    }
  }

  .content {
    width: 100%;
    overflow: hidden;
    padding: 30px 17px 0px 17px;
    margin-bottom: 40px;
    border: 1px solid #eeeeee;
    border-radius: 2px;
  }
`

interface InfoBlockProps {
  title: string | React.ReactNode
  icon?: string
  children?: any
}

function InfoBlock(props: InfoBlockProps) {
  const [visible, setVisible] = React.useState(true)
  const { icon, title, children } = props

  return (
    <Wrapper>
      <div className='title'>
        {icon ? <Icon type={icon} /> : null}
        <span onClick={() => setVisible(!visible)}>
          <span className='text'>{title}</span>
          <Icon className='icon' type={visible ? 'up' : 'down'} />
        </span>
      </div>
      <div className='content' style={{ display: visible ? 'block' : 'none' }}>
        {children}
      </div>
    </Wrapper>
  )
}
export default InfoBlock
