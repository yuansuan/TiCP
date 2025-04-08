/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'
import { Empty } from 'antd'

import { FileSystem } from '@/components'
import { RootPoint } from '@/domain/FileSystem'

const Wrapper = styled.div`
  height: 800px;
  margin-left: -20px;
`
// jobDir 当前作业路径
function JobFileSystem({ jobDir }) {
  return (
    <Wrapper>
      {jobDir ? (
        <FileSystem
          points={[
            new RootPoint({
              pointId: jobDir,
              path: jobDir,
              name: jobDir.split(/[\\/]/).pop(),
            }),
          ]}
          defaultPath={jobDir}
        />
      ) : (
        <Empty description='作业目录不存在' />
      )}
    </Wrapper>
  )
}

export default JobFileSystem
