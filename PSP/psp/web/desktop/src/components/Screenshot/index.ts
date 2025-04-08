/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import html2canvas from 'html2canvas'
import './style.less'
import { Playground } from './Playground'

type Props = {
  getRoot: () => HTMLElement
}
export class Screenshot {
  private getRoot: () => HTMLElement
  playground: Playground

  constructor(props?: Partial<Props>) {
    const { getRoot = () => document.body } = props || {}
    this.getRoot = getRoot
  }

  destroy = () => {
    this.playground?.destroy()
    this.playground = null
  }

  capture = async () => {
    const root = this.getRoot()
    this.playground = new Playground(root)

    // set imageSrc for video
    const videoes = document.getElementsByTagName('video')
    if (videoes.length > 0) {
      const canvas = document.createElement('canvas')
      document.body.appendChild(canvas)
      try {
        Array.from(videoes).forEach(video => {
          const ctx = canvas.getContext('2d')
          ctx.drawImage(video, 0, 0, video.videoWidth, video.videoHeight)
          video.style.backgroundImage = `url(${canvas.toDataURL()})`
          video.style.backgroundSize = 'cover'
        })
      } catch (err) {
        console.error(err)
      }
      document.body.removeChild(canvas)
    }

    const canvas = await html2canvas(root)
    this.playground.setImageSource(canvas)
    return new Promise(resolve => {
      this.playground.addListener('onDestroy', () => {
        resolve(null)
      })
    })
  }
}
