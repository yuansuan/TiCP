/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

type ImageProps = React.ImgHTMLAttributes<HTMLImageElement>
const Image: React.SFC<ImageProps> & {
  NotFound?: () => JSX.Element
  Empty?: () => JSX.Element
} = function Image(props: ImageProps) {
  return <img {...props} />
}

Image.NotFound = (props?: ImageProps) => (
  <img src={require('../assets/images/notFound.png')} {...props} />
)

Image.Empty = (props?: ImageProps) => (
  <img src={require('../assets/images/nodata.png')} {...props} />
)

export default Image
