/* Copyright (C) 2016-present, Yuansuan.cn */

import { PluginProps, ClipRect, Region } from './type'

export function Resizable({ core, utils, constants }: PluginProps) {
  const { noNegative } = utils
  let mouseStartX
  let mouseStartY
  let activeNode: Region
  let tempRect: ClipRect

  // 获取边框节点相关信息
  function getClipNodes(): Region[] {
    if (!core.clipRect) {
      return []
    }

    const borderSize = constants.BORDER_SIZE
    // 获取裁剪框位置信息
    const { startX, startY, width, height } = tempRect || core.clipRect
    const halfBorderSize = borderSize / 2

    return [
      //n
      {
        x: startX + halfBorderSize,
        y: startY,
        width: width - borderSize,
        height: halfBorderSize,
        style: {
          cursor: 'ns-resize'
        },
        state: {
          direction: 'n'
        }
      },
      {
        x: startX - halfBorderSize + width / 2,
        y: startY - halfBorderSize,
        width: borderSize,
        height: halfBorderSize,
        style: {
          cursor: 'ns-resize'
        },
        state: {
          direction: 'n'
        }
      },
      // s
      {
        x: startX + halfBorderSize,
        y: startY - halfBorderSize + height,
        width: width - borderSize,
        height: halfBorderSize,
        style: {
          cursor: 'ns-resize'
        },
        state: {
          direction: 's'
        }
      },
      {
        x: startX - halfBorderSize + width / 2,
        y: startY + height,
        width: borderSize,
        height: halfBorderSize,
        style: {
          cursor: 'ns-resize'
        },
        state: {
          direction: 's'
        }
      },
      // w
      {
        x: startX,
        y: startY + halfBorderSize,
        width: halfBorderSize,
        height: height - borderSize,
        style: {
          cursor: 'ew-resize'
        },
        state: {
          direction: 'w'
        }
      },
      {
        x: startX - halfBorderSize,
        y: startY - halfBorderSize + height / 2,
        width: halfBorderSize,
        height: borderSize,
        style: {
          cursor: 'ew-resize'
        },
        state: {
          direction: 'w'
        }
      },
      //e
      {
        x: startX - halfBorderSize + width,
        y: startY + halfBorderSize,
        width: halfBorderSize,
        height: height - borderSize,
        style: {
          cursor: 'ew-resize'
        },
        state: {
          direction: 'e'
        }
      },
      {
        x: startX + width,
        y: startY - halfBorderSize + height / 2,
        width: halfBorderSize,
        height: borderSize,
        style: {
          cursor: 'ew-resize'
        },
        state: {
          direction: 'e'
        }
      },
      // nw
      {
        x: startX - halfBorderSize,
        y: startY - halfBorderSize,
        width: borderSize,
        height: borderSize,
        style: {
          cursor: 'nwse-resize'
        },
        state: {
          direction: 'nw'
        }
      },
      // se
      {
        x: startX - halfBorderSize + width,
        y: startY - halfBorderSize + height,
        width: borderSize,
        height: borderSize,
        style: {
          cursor: 'nwse-resize'
        },
        state: {
          direction: 'se'
        }
      },
      //ne
      {
        x: startX - halfBorderSize + width,
        y: startY - halfBorderSize,
        width: borderSize,
        height: borderSize,
        style: {
          cursor: 'nesw-resize'
        },
        state: {
          direction: 'ne'
        }
      },
      // sw
      {
        x: startX - halfBorderSize,
        y: startY - halfBorderSize + height,
        width: borderSize,
        height: borderSize,
        style: {
          cursor: 'nesw-resize'
        },
        state: {
          direction: 'sw'
        }
      }
    ]
  }

  // 调整裁剪框尺寸
  function resize(props: {
    currentX: number
    currentY: number
    direction: string
  }) {
    const { startX, startY, width, height } = core.clipRect
    const { currentX, currentY, direction } = props

    let res = {
      ...core.clipRect
    }

    if (direction.includes('w')) {
      res.startX = currentX - (startX + width) > 0 ? startX + width : currentX
      res.width = noNegative(width - (currentX - startX))
    }
    if (direction.includes('n')) {
      res.startY = currentY - (startY + height) > 0 ? startY + height : currentY
      res.height = noNegative(height - (currentY - startY))
    }
    if (direction.includes('s')) {
      res.height = noNegative(currentY - startY)
    }
    if (direction.includes('e')) {
      res.width = noNegative(currentX - startX)
    }

    return res
  }

  function onMouseDown(event) {
    mouseStartX = noNegative(event.offsetX)
    mouseStartY = noNegative(event.offsetY)

    if (core.clipRect) {
      // 判断用户是否想要缩放裁剪框
      const node = core.getRegionByPos({
        x: mouseStartX,
        y: mouseStartY
      })

      if (node?.state?.direction) {
        activeNode = node
        core.emit('onStart')
      } else {
        activeNode = null
      }
    }
  }

  function onMouseMove(event) {
    const currentX = event.offsetX
    const currentY = event.offsetY

    // 重置鼠标样式
    const node = core.getRegionByPos({
      x: currentX,
      y: currentY
    })

    if (node) {
      core.container.style.cursor = node.cursor
    } else {
      core.container.style.cursor = 'default'
    }

    // 缩放裁剪框
    if (activeNode) {
      tempRect = resize({
        currentX,
        currentY,
        direction: activeNode.state.direction
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

    // 清理缩放中的裁剪框节点
    activeNode = undefined
  }

  core.registerRegions([getClipNodes])

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
