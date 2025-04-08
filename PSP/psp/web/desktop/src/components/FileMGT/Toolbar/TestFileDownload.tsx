/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { downloadTestFile } from '@/utils'

export const TestFileDownload = () => (
  <Button onClick={downloadTestFile}>下载示例文件</Button>
)
