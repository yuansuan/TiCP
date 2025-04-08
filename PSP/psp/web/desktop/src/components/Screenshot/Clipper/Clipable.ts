/* Copyright (C) 2016-present, Yuansuan.cn */

import { PluginProps, ClipRect } from './type'

export function Clipable({ core, utils }: PluginProps) {
  const { noNegative } = utils
  let mouseStartX
  let mouseStartY
  let tempRect: ClipRect

  function onMouseDown(event) {
    mouseStartX = noNegative(event.offsetX)
    mouseStartY = noNegative(event.offsetY)

    // 开始绘制裁剪框
    if (!core.clipRect) {
      tempRect = {
        startX: mouseStartX,
        startY: mouseStartY,
        width: 0,
        height: 0
      }

      core.emit('onStart')
    }
  }

  function onMouseMove(event) {
    const currentX = event.offsetX
    const currentY = event.offsetY

    // 绘制裁剪框
    if (!core.clipRect && tempRect) {
      const { startX, startY } = tempRect
      // 裁剪框临时宽高
      tempRect = {
        ...tempRect,
        width: currentX - startX,
        height: currentY - startY
      }

      core.redraw(tempRect)
    }
  }

  function onMouseUp() {
    // 清理临时裁剪框
    if (tempRect) {
      if (tempRect.width > 0 && tempRect.height > 0) {
        core.clipRect = tempRect
        core.emit('onComplete')
      }
      tempRect = undefined
    }
  }

  const { container } = core
  container.addEventListener('mousedown', onMouseDown)
  container.addEventListener('mousemove', onMouseMove)
  container.addEventListener('mouseup', onMouseUp)

  return () => {
    container.removeEventListener('mousedown', onMouseDown)
    container.removeEventListener('mousemove', onMouseMove)
    container.removeEventListener('mouseup', onMouseUp)
  }
}
