/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Modal, Icon } from '@/components'
import { getImageNaturalSize, imageSize } from './utils'
import ImageSize from './ImageSize'

const StyledLayout = styled.div`
  width: 100%;
  height: 100%;
  background-color: #f6f8fa;
  display: flex;
  flex-direction: column;
  padding-bottom: 10px;

  > .image {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    > img {
      max-width: calc(100vw - 40px);
      max-height: calc(100vh - 40px - 55px - 22px - 10px);
      user-select: none;
    }
  }
`

function download(url, fileName) {
  const aEl = document.createElement('a')
  aEl.href = url
  if (fileName) {
    aEl.download = fileName
  }
  document.body.appendChild(aEl)
  aEl.click()
  document.body.removeChild(aEl)
}

interface ImageAreaPropType {
  imageSrc: string
  fileName: string
}

function ImageArea({ imageSrc, fileName }: ImageAreaPropType) {
  const [imageNaturalSize, setImageNaturalSize] = React.useState({
    width: 0,
    height: 0
  })
  const [zoom, setZoom] = React.useState(1)

  const img = React.useMemo(() => {
    const img = new Image()
    img.src = imageSrc
    return img
  }, [imageSrc])

  React.useEffect(() => {
    const imageOnload = () => {
      const imageNaturalSize: imageSize = getImageNaturalSize(img)
      setImageNaturalSize(imageNaturalSize)
    }
    // 图片加载完毕后，获取加载到的图片的尺寸，将其更新到 state
    img.addEventListener('load', imageOnload)
    // ummount 图片组件时，unsubscribe onload 事件，以避免对 unmounted 组件更新 state
    return () => {
      img.removeEventListener('load', imageOnload)
    }
  }, [img])

  function narrow() {
    setZoom(Math.max(0.1, zoom - 0.1))
  }
  function enlarge() {
    setZoom(zoom + 0.1)
  }
  function revert() {
    setZoom(1)
  }

  return (
    <StyledLayout>
      <Modal.Toolbar
        actions={[
          {
            tip: '缩小',
            slot: <Icon type='narrow' onClick={narrow} />
          },
          {
            tip: '放大',
            slot: <Icon type='enlarge' onClick={enlarge} />
          },
          {
            tip: '重置',
            slot: <Icon type='revert' onClick={revert} />
          },
          {
            tip: '下载',
            slot: (
              <Icon
                type='download'
                onClick={() => download(imageSrc, fileName)}
              />
            )
          }
        ]}
      />
      <ImageSize imageSize={imageNaturalSize}></ImageSize>
      <div className='image'>
        <img src={imageSrc} style={{ transform: `scale(${zoom})` }} />
      </div>
    </StyledLayout>
  )
}

interface PreviewImageType {
  fileName: string
  src: string
}

export function previewImage({ fileName, src }: PreviewImageType) {
  Modal.show({
    title: fileName,
    footer: null,
    content: <ImageArea imageSrc={src} fileName={fileName} />,
    centered: true,
    width: window.innerWidth - 40,
    bodyStyle: {
      height: 'calc(100vh - 55px - 40px)',
      padding: 0
    }
  })
}
