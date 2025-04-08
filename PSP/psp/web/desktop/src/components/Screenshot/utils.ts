/* Copyright (C) 2016-present, Yuansuan.cn */

export function drawMask(canvas: HTMLCanvasElement) {
  const ctx = canvas.getContext('2d')
  const { width, height } = canvas.getBoundingClientRect()

  // 清除画布
  ctx.clearRect(0, 0, width, height)
  // 绘制蒙层
  ctx.save()
  ctx.fillStyle = 'rgba(0, 0, 0, .6)'
  ctx.fillRect(0, 0, width, height)
}
