/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useLayoutEffect } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Modal } from '@/components'
import { showMask } from '@/components'
import { Button } from 'antd'
import { env } from '@/domain'
import { lightStore } from './Light'

const StyledDiv = styled.div`
  > .container {
    margin-top: 20px;
    margin-bottom: 20px;
    min-height: 200px;
    position: relative;
    display: flex;
    flex-flow: column wrap;
    justify-content: center;
    align-items: center;
    cursor: pointer;

    > img {
      width: 100%;
    }

    > .play {
      position: absolute;
    }

    &:hover {
      > .play {
        transform: scale(1.1);
      }
    }
  }
`

const StyledContent = styled.div`
  width: 100vw;
  height: 100vh;
  position: relative;
  padding: 60px 24px 98px;

  > iframe {
    width: 100%;
    height: 100%;
  }

  > video {
    width: 100%;
    height: 100%;
  }

  > div {
    position: absolute;
    margin-top: 7px;
    margin-bottom: 7px;
    width: 100%;
    text-align: center;

    > .close-btn {
      display: inline;
      cursor: pointer;
    }
  }
`

type Props = {
  nextBtnText?: string
  onCancel: any
  onOk: () => void
}

const VideoContent = ({ onClose }) => (
  <StyledContent>
    <video
      controls
      autoPlay
      onClick={e => {
        e.stopPropagation()
        e.nativeEvent?.stopImmediatePropagation()
      }}>
      <source
        src='https://euc-platform-bucket.yuansuan.cn/videos/job_creator.mov'
        type='video/mp4'
      />
    </video>
    <div>
      <div className='close-btn' onClick={onClose}>
        <img src={require('@/assets/images/watch_mask_close.svg')} alt='关闭' />
      </div>
    </div>
  </StyledContent>
)

const Content = observer(function Content({ onOk, nextBtnText }: Props) {
  function playVideo() {
    showMask({
      content: ({ onClose }) => <VideoContent onClose={onClose} />
    })
  }

  useLayoutEffect(() => {
    lightStore.setLightExpanded(true)
  }, [])

  return (
    <StyledDiv>
      <div>快速入门 {'>'} 作业提交</div>
      <div className='container' onClick={playVideo}>
        <img src={require('@/assets/images/video_bg.png')} alt='video-bg' />
        <div className='play'>
          <img src={require('@/assets/images/watch.svg')} alt='watch' />
        </div>
      </div>
      <Modal.Footer
        className='footer'
        OkButton={
          <Button
            type='primary'
            style={{ marginRight: 0 }}
            onClick={() => {
              onOk()
              lightStore.setLightExpanded(false)
            }}>
            {nextBtnText || '继续'}
          </Button>
        }
        CancelButton={null}
      />
    </StyledDiv>
  )
})

export async function showNewComerModal(
  props?: Omit<Props, 'onCancel' | 'onOk'>
) {
  let title
  if (env.isFuture) {
    title = '智算未来'
  } else if (env.isKaiwu) {
    title = '开物平台'
  } else if (env.custom.id) {
    title = env.custom.title
  } else {
    title = '云仿真平台'
  }
  return await Modal.show({
    width: 600,
    title: `使用${title}强大的集群算力`,
    content: ({ onCancel, onOk }) => (
      <Content onCancel={onCancel} onOk={onOk} {...props} />
    ),
    footer: null
  })
}
