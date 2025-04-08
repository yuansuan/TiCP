/* Copyright (C) 2016-present, Yuansuan.cn */

import { drawMask } from '../utils'
import { theme } from '@/constant'
import { ClipRect, RegionFactory } from './type'
import { BORDER_SIZE, default as constants } from './constants'
import { default as utils } from './utils'
import { Clipable } from './Clipable'
import { Moveable } from './Movable'
import { Resizable } from './Resizable'
import EventEmitter from 'eventemitter3'

export class Clipper extends EventEmitter {
  container: HTMLCanvasElement
  drawBackground: () => void = () => {}

  dispatchers = []

  // clip rect instance
  clipRect: ClipRect

  // regions
  private regions: Array<RegionFactory> = []

  constructor(props: {
    container: HTMLCanvasElement
    drawBackground: () => void
  }) {
    super()

    this.container = props.container
    this.drawBackground = props.drawBackground

    // register plugins
    this.dispatchers = [Clipable, Moveable, Resizable].map(plugin =>
      plugin(this.pluginProps)
    )
    this.container.addEventListener('mousemove', this.onMouseMove)
  }

  get pluginProps() {
    return {
      core: this,
      utils,
      constants
    }
  }

  get context() {
    return this.container?.getContext('2d')
  }

  redraw = (rect?: ClipRect) => {
    const ctx = this.container.getContext('2d')

    ctx.save()
    drawMask(this.container)

    if (rect) {
      const { startX, startY, width, height } = rect
      // 绘制裁剪框
      // 将蒙层凿开
      ctx.globalCompositeOperation = 'source-atop'
      // 裁剪选择框
      ctx.clearRect(startX, startY, width, height)
      // 绘制8个边框像素点并保存坐标信息以及事件参数
      ctx.globalCompositeOperation = 'source-over'
      ctx.fillStyle = theme.linkColor
      // 像素点大小
      const size = BORDER_SIZE
      // 绘制像素点
      ctx.fillRect(startX - size / 2, startY - size / 2, size, size)
      ctx.fillRect(startX - size / 2 + width / 2, startY - size / 2, size, size)
      ctx.fillRect(startX - size / 2 + width, startY - size / 2, size, size)
      ctx.fillRect(
        startX - size / 2,
        startY - size / 2 + height / 2,
        size,
        size
      )
      ctx.fillRect(
        startX - size / 2 + width,
        startY - size / 2 + height / 2,
        size,
        size
      )
      ctx.fillRect(startX - size / 2, startY - size / 2 + height, size, size)
      ctx.fillRect(
        startX - size / 2 + width / 2,
        startY - size / 2 + height,
        size,
        size
      )
      ctx.fillRect(
        startX - size / 2 + width,
        startY - size / 2 + height,
        size,
        size
      )
    }

    // 绘制结束
    ctx.restore()

    // 将图片绘制在蒙层下方
    ctx.save()
    ctx.globalCompositeOperation = 'destination-over'
    this.drawBackground()
  }

  clip = () => {
    // 获取裁剪区域位置信息
    const { startX, startY, width, height } = this.clipRect
    const borderSize = constants.BORDER_SIZE

    // 获取裁剪框区域图片信息
    const img = this.context.getImageData(
      startX + borderSize / 2.4,
      startY + borderSize / 2.4,
      width - borderSize * 2.4,
      height - borderSize * 2.4
    )
    // 创建canvas标签，用于存放裁剪区域的图片
    const canvas = document.createElement('canvas')
    canvas.width = width - borderSize * 2.4
    canvas.height = height - borderSize * 2.4
    // 获取裁剪框区域画布
    const imgContext = canvas.getContext('2d')
    if (imgContext) {
      // 将图片放进裁剪框内
      imgContext.putImageData(img, 0, 0)
      const a = document.createElement('a')
      // 获取图片
      a.href = canvas.toDataURL('png')
      // 下载图片
      a.download = `截图_${new Date().getTime()}.png`
      a.click()
    }
  }

  destroy = () => {
    this.dispatchers.forEach(dispatch => dispatch())

    this.removeAllListeners('onStart')
    this.removeAllListeners('onComplete')
    this.container.removeEventListener('mousemove', this.onMouseMove)
  }

  registerRegions(regions: Array<RegionFactory>) {
    this.regions = [...this.regions, ...regions].filter(Boolean)
  }

  getRegionByPos(pos: { x: number; y: number }) {
    const { regions } = this
    const nodes = regions.reduce((acc, getRegion) => {
      let nodes = getRegion()
      nodes = Array.isArray(nodes) ? nodes : [nodes]

      return [...acc, ...nodes].filter(Boolean)
    }, [])

    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d')
    // 设置裁剪框鼠标响应
    let node = null
    // 判断鼠标位置
    context.beginPath()
    for (let i = 0; i < nodes.length; i++) {
      const { x, y, width, height } = nodes[i]
      context.rect(x, y, width, height)
      if (context.isPointInPath(pos.x, pos.y)) {
        node = nodes[i]
        break
      }
    }
    context.closePath()

    return node
  }

  onMouseMove = event => {
    const currentX = event.offsetX
    const currentY = event.offsetY

    // 重置鼠标样式
    const node = this.getRegionByPos({
      x: currentX,
      y: currentY
    })

    if (node) {
      this.container.style.cursor = node.style.cursor
    } else {
      this.container.style.cursor = 'default'
    }
  }
}
