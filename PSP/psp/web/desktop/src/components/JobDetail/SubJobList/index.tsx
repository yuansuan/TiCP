/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'

import SubJobListCom from '@/components/JobList/SubJobList'

const Wrapper = styled.div`
  height: 800px;
`
function SubJobList({ jobId, isHistory, workspaceId }) {
  return (
    <Wrapper>
      <SubJobListCom
        jobId={jobId}
        isHistory={isHistory}
        workspaceId={workspaceId}
      />
    </Wrapper>
  )
}

export default SubJobList
