/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

export function openVisualApp(url: string) {
  if ((<any>window).__devtron) {
    window.open(url)
  } else {
    const win = window.open()
    win.location.href = url
  }
}
