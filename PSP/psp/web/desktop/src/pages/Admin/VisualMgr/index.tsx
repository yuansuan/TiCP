/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import VisIBV from '@/pages/Admin/VisIBV'
const Wrapper = styled.div`
  width: 100%;
  .ant-breadcrumb {
    padding: 20px 50px;
  }

  .content {
    flex: 1;
  }

  .detail-content {
    padding: 0px 50px;
  }
`

const VisualMgr = observer(() => {
  return (
    <Wrapper>
      <VisIBV />
    </Wrapper>
  )
})

export default VisualMgr
