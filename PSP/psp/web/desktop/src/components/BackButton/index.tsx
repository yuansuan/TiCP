import React from 'react'
import { LeftOutlined } from '@ant-design/icons'
import styled from 'styled-components'

const BackButtonWrapper = styled.div`
  display: flex;
  justify-content: flex-start;
  align-items: center;
  height: 30px;
  cursor: pointer;

  &:hover {
    color: #4193F7;
  }
`

export default function BackButton({ onClick, title, children, ...props }) {
  return (
    <BackButtonWrapper onClick={onClick} title={title} {...props}>
    <LeftOutlined />{children || ''}
    </BackButtonWrapper>
  )
}
