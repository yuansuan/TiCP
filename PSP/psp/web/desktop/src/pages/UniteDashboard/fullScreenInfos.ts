import { observable, computed } from 'mobx'

class FullScreenInfo {
  @observable height = 0
  @observable width = 0
  @observable isFullScreen = false

  @computed
  get baseSize() {
    return this.width / 80
  }

  @computed
  get baseHeight() {
    return this.height - this.baseSize * 3
  }

  @computed
  get baseRow() {
    return this.height * 0.33 - this.baseSize
  }

  @computed
  get firstRow() {
    return this.height * 0.4 - this.baseSize
  }

  @computed
  get otherRow() {
    return this.height * 0.3 - this.baseSize
  }

  @computed
  get chartHeight() {
    return this.isFullScreen ? this.otherRow * 0.68 : 180
  }

  @computed
  get onlyLineChartHeight() {
    return this.isFullScreen ? this.otherRow * 0.18 : 50
  }
}

export const fullScreenInfo = new FullScreenInfo()
