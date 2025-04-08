/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import styled from 'styled-components'
import Chooser from '..'
import Button from '../../Button'

const StyledButtonType = styled.div`
  margin: 5px 5px;

  > button {
    margin: 5px;
  }
`

export const Basic = () => {
  const [searchable, setSearchable] = useState(false)
  const [styled, setStyle] = useState(false)
  const [selectedKeys, setSelectedKeys] = useState(null)

  return (
    <StyledButtonType>
      <Button
        type={searchable ? 'default' : 'primary'}
        onClick={() => {
          setSearchable(!searchable)
        }}>
        {searchable ? '关闭' : '打开'}搜索框
      </Button>
      <Button
        type={styled ? 'default' : 'primary'}
        onClick={() => {
          setStyle(!styled)
        }}>
        样式
      </Button>
      <Chooser
        items={[
          ['1-apple', 'Apple'],
          ['2-pear', 'Pear'],
          ['3-orange', 'Orange'],
        ]}
        searchable={searchable}
        style={styled ? { lineHeight: '40px', color: 'green' } : {}}
        onChange={selectedKeys => {
          setSelectedKeys(selectedKeys)
        }}
      />
      <span>selectedKeys: {selectedKeys && selectedKeys.join(',')}</span>
    </StyledButtonType>
  )
}
