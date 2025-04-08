/* Copyright (C) 2016-present, Yuansuan.cn */

import { PluginProps, ClipRect } from './type'

export function Moveable({ core, utils, constants }: PluginProps) {
  const { noNegative, confine } = utils
  let mouseStartX
  let mouseStartY
  let moving = false
  let tempRect: ClipRect

  // get region info for movement
  function getRegion() {
    if (!core.clipRect) {
      return undefined
    }

    const borderSize = constants.BORDER_SIZE
    // 获取裁剪框位置信息
    const { startX, startY, width, height } = tempRect || core.clipRect
    const halfBorderSize = borderSize / 2

    return {
      x: startX + halfBorderSize,
      y: startY + halfBorderSize,
      width: width - borderSize,
      height: height - borderSize,
      style: {
        cursor: 'move'
      },
      state: {
        id: 'regionForMove'
      }
    }
  }

  // 移动裁剪框
  function move(props: { currentX: number; currentY: number }) {
    const { startX, startY, width, height } = core.clipRect
    const { currentX, currentY } = props
    const { container } = core

    // 计算要移动到的x轴坐标
    const x = confine(currentX - (mouseStartX - startX), width, container.width)
    // 计算要移动到的y轴坐标
    const y = confine(
      currentY - (mouseStartY - startY),
      height,
      container.height
    )

    return {
      startX: x,
      startY: y,
      width,
      height
    }
  }

  function onMouseDown(event) {
    mouseStartX = noNegative(event.offsetX)
    mouseStartY = noNegative(event.offsetY)

    if (core.clipRect) {
      // 判断用户是否想要移动裁剪框
      const node = core.getRegionByPos({
        x: mouseStartX,
        y: mouseStartY
      })

      if (node.state.id === 'regionForMove') {
        moving = true
        core.emit('onStart')
      }
    }
  }

  function onMouseMove(event) {
    const currentX = event.offsetX
    const currentY = event.offsetY

    // 移动裁剪框
    if (moving) {
      tempRect = move({
        currentX,
        currentY
      })
      core.redraw(tempRect)
    }
  }

  function onMouseUp() {
    // 清理临时裁剪框
    if (tempRect) {
      core.clipRect = tempRect
      tempRect = undefined

      core.emit('onComplete')
    }

    // 清理移动标志
    moving = false
  }

  core.registerRegions([getRegion])

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
