/* Copyright (C) 2016-present, Yuansuan.cn */

import './style.less'

type Option = {
  id: string
}

export class Toolbar {
  options: Option[] = [
    {
      id: 'confirm'
    },
    {
      id: 'close'
    }
  ]

  container: HTMLElement
  self: HTMLDivElement = document.createElement('div')

  constructor(
    container: HTMLElement,
    props?: Partial<{
      onClick: (option: Option) => void
    }>
  ) {
    this.container = container

    const { options } = this

    this.self.id = 'toolPanel'
    for (let i = 0; i < options.length; i++) {
      const item = options[i]
      const itemPanel = document.createElement('div')
      itemPanel.setAttribute('class', item.id)

      itemPanel.addEventListener('click', () => {
        props?.onClick(item)
      })

      this.self.appendChild(itemPanel)
    }

    this.container.appendChild(this.self)
    this.hide()
  }

  hide = () => {
    this.self.style.display = 'none'
  }

  show = (props: {
    startX: number
    startY: number
    width: number
    height: number
  }) => {
    this.self.style.display = 'block'

    const { startX, startY, width, height } = props
    // 工具栏X轴坐标 = (裁剪框的宽度 - 工具栏的宽度) + 裁剪框距离左侧的距离
    const mouseX = (width - this.self.offsetWidth) / 2 + startX
    // 工具栏Y轴坐标
    let mouseY = startY + height + 10
    if (width < 0 && height < 0) {
      // 从右下角拖动时，工具条y轴的位置应该为position.startY + 10
      mouseY = startY + 10
    }

    this.self.style.left = mouseX + 'px'
    this.self.style.top = mouseY + 'px'
  }

  destroy = () => {
    if (this.self) {
      this.container.removeChild(this.self)
    }
  }
}
