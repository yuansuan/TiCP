/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

export function runInTyping(formModel: any, fn: () => void) {
  formModel.isTyping = true
  fn()
  setTimeout(() => (formModel.isTyping = false), 0)
}
