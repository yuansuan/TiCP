/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import screenfull from 'screenfull'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Tooltip, Checkbox } from 'antd'
import { Icon } from '@/components'

const StyledLayout = styled.div`
  display: flex;
  height: 100%;

  > .right {
    margin-left: auto;
    display: flex;
    height: 100%;
    .reverse {
      margin-left: 10px;
    }
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
  refresh: () => void
}

export const Toolbar = observer(function Toolbar({ editor, refresh }: Props) {
  const state = useLocalStore(() => ({
    fullscreen: false,
    setFullscreen(flag) {
      this.fullscreen = flag
    },
    isAutoRefrsh: false,
    setAutoRefresh(flag) {
      this.isAutoRefrsh = flag
    }
  }))
  let intervalId = null
  const { fullscreen, isAutoRefrsh } = state
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
      if(intervalId) {
        clearInterval(intervalId)
      }
    }
  }, [editor])

  function toggleFullscreen() {
    if (screenfull.isEnabled) {
      const $modal = document.getElementsByClassName('__fileEditor__')[0]
      screenfull.toggle($modal)
    }
  }
  const onAutoRefreshCheck = checked => {
    state.setAutoRefresh(checked)

    if (checked) {
      refresh()
      intervalId && clearInterval(intervalId)

      intervalId = setInterval(() => {
        refresh()
      }, 5 * 1000)
    } else {
      clearInterval(intervalId)
    }
  }
  return (
    <StyledLayout>
      <div className='right'>
        <Tooltip title='查找'>
          <div>
            <Icon type='search' onClick={find}  style={{ fontSize: 18 }}  />
          </div>
        </Tooltip>
        <Tooltip title='刷新'>
          <div>
            <Icon type='revert' onClick={refresh} />
            <Checkbox
              className='reverse'
              checked={isAutoRefrsh}
              onChange={e => onAutoRefreshCheck(e.target.checked)}>
              自动刷新
              <Tooltip title={'每5秒自动刷新内容'}>
                <Icon type={'help-circle'} style={{ fontSize: 24,paddingTop: 6 }} />
              </Tooltip>
            </Checkbox>
          </div>
        </Tooltip>

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
