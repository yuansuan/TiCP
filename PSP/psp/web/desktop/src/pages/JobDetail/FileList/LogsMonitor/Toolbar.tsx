/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import screenfull from 'screenfull'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Tooltip, Popconfirm } from 'antd'
import { Icon } from '@/components'
import { ClearOutlined } from '@ant-design/icons'

const StyledLayout = styled.div`
  display: flex;
  height: 100%;

  > .right {
    margin-left: auto;
    display: flex;
    height: 100%;

    > div {
      display: flex;
      cursor: pointer;
      height: 100%;
      padding: 0 10px;
      align-items: center;

      &:hover {
        background-color: ${({ theme }) => theme.backgroundColorBase};
      }
    }
  }
`

type Props = {
  editor: any
  showRefreshAction: boolean
  refresh: () => void
  clearScreen: () => void
}

export const Toolbar = observer(function Toolbar({
  editor,
  refresh,
  clearScreen,
  showRefreshAction
}: Props) {
  const state = useLocalStore(() => ({
    fullscreen: false,
    setFullscreen(flag) {
      this.fullscreen = flag
    }
  }))
  const { fullscreen } = state

  function find() {
    editor.trigger('', 'actions.find')
  }

  useEffect(() => {
    function onFullscreenChange() {
      state.setFullscreen((screenfull as any).isFullscreen)

      const $modal = document.getElementsByClassName('__fileEditor__')[0]
      const $content: any = $modal.querySelector('.ant-modal-content')
      const $body: any = $modal.querySelector('.ant-modal-body')
      $content.style.height = state.fullscreen ? '100%' : 'auto'
      $body.style.height = state.fullscreen ? 'calc(100% - 55px)' : '600px'
      setTimeout(() => editor.layout(), 300)
    }

    if (screenfull.isEnabled) {
      screenfull.on('change', onFullscreenChange)
    }

    return () => {
      if (screenfull.isEnabled) {
        screenfull.off('change', onFullscreenChange)
      }
    }
  }, [editor])

  function toggleFullscreen() {
    if (screenfull.isEnabled) {
      const $modal = document.getElementsByClassName('__fileEditor__')[0]
      screenfull.toggle($modal)
    }
  }

  return (
    <StyledLayout>
      <div className='right'>
        <Tooltip title='查找'>
          <div>
            <Icon type='search' onClick={find} />
          </div>
        </Tooltip>
        <Tooltip title='清屏'>
          <Popconfirm
            title='确认要清除屏幕上日志吗?'
            placement='bottom'
            onConfirm={clearScreen}
            okText='确认'
            cancelText='取消'>
            <div>
              <ClearOutlined />
            </div>
          </Popconfirm>
        </Tooltip>
        {showRefreshAction && (
          <Tooltip title='重连'>
            <div>
              <img
                src={require('@/assets/images/connection.png')}
                alt='重连'
                onClick={refresh}
              />
            </div>
          </Tooltip>
        )}
        {!fullscreen && (
          <Tooltip title='全屏'>
            <div>
              <Icon type='full_screen' onClick={toggleFullscreen} />
            </div>
          </Tooltip>
        )}
        {fullscreen && (
          <Tooltip title='退出全屏'>
            <div>
              <Icon type='full_screen' onClick={toggleFullscreen} />
            </div>
          </Tooltip>
        )}
      </div>
    </StyledLayout>
  )
})
