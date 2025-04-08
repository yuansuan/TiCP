/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

// 直接测试文件
export const downloadTestFile = async () => {
  const aEl = document.createElement('a')
  aEl.href = '/api/cos/download/Yuansuan_Demo.zip'
  document.body.appendChild(aEl)
  aEl.click()
  document.body.removeChild(aEl)
}
