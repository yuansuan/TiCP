/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

type imageSize = {
  width: number
  height: number
}

function getImageNaturalSize(img: HTMLImageElement): imageSize {
  return {
    width: img.naturalWidth,
    height: img.naturalHeight,
  }
}

export { imageSize, getImageNaturalSize }
