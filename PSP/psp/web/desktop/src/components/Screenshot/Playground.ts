/* Copyright (C) 2016-present, Yuansuan.cn */

import { Clipper } from './Clipper'
import { Toolbar } from './Toolbar'
import EventEmitter from 'eventemitter3'

export class Playground extends EventEmitter {
  container: HTMLCanvasElement = document.createElement('canvas')
  clipper: Clipper
  toolbar: Toolbar
  imageSource: CanvasImageSource
  private root: HTMLElement

  constructor(root: HTMLElement) {
    super()

    this.root = root
    this.container.id = 'screenshot-playground'
    this.container.width = this.root.clientWidth
    this.container.height = this.root.clientHeight

    root.appendChild(this.container)
    this.clipper = new Clipper({
      container: this.container,
      drawBackground: this.draw
    })
    this.toolbar = new Toolbar(this.root, {
      onClick: ({ id }) => {
        if (id === 'close') {
          this.destroy()
        } else if (id === 'confirm') {
          this.clipper.clip()
          this.destroy()
        }
      }
    })

    // toggle toolbar
    this.clipper.addListener('onStart', () => {
      this.toolbar.hide()
    })
    this.clipper.addListener('onComplete', () => {
      this.toolbar.show(this.clipper.clipRect)
    })

    window.addEventListener('keydown', this.onEsc)
  }

  get context() {
    return this.container.getContext('2d')
  }

  setImageSource(imageSource: CanvasImageSource) {
    this.imageSource = imageSource
    this.clipper.redraw()
  }

  draw = () => {
    if (!this.imageSource) {
      return
    }

    const { width, height } = this.container.getBoundingClientRect()
    const { context } = this

    context.drawImage(this.imageSource, 0, 0, width, height)
  }

  destroy = () => {
    if (this.container) {
      this.root.removeChild(this.container)
    }

    this.clipper.destroy()
    this.toolbar.destroy()

    window.removeEventListener('keydown', this.onEsc)

    this.emit('onDestroy')
  }

  onEsc = e => {
    if (e.key === 'Escape') {
      this.destroy()
    }
  }
}
