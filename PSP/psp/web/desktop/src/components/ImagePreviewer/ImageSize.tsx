/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import React from 'react'
import styled from 'styled-components'

const ImageSizeContainer = styled.p`
  text-align: center;
  margin: 0;
  user-select: none;
`

function ImageSize({ imageSize }) {
  return (
    <ImageSizeContainer>{`${imageSize.width} px * ${imageSize.height} px`}</ImageSizeContainer>
  )
}

export default ImageSize
