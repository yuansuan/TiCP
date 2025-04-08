/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import PartialList from '..'

const StyledItem = styled.span`
  display: inline-block;
  width: 30px;
  height: 30px;
  margin: 3px;
  line-height: 30px;
  text-align: center;
  border-radius: 50%;
  border: 1px solid #ccc;
`

export const CustomItem = () => (
  <PartialList
    style={{ width: 200, margin: 20 }}
    itemMapper={(item, index) => <StyledItem key={index}>{item}</StyledItem>}
    items={[0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18]}
  />
)
