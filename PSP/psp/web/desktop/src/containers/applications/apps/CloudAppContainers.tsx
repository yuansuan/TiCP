/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect, useCallback, useRef } from 'react'
import { useSelector } from 'react-redux'
import hotkeys from 'hotkeys-js'
import { ToolBar } from '@/utils/general'
import VisList from '@/pages/VisList/List'

export const CloudAppContainers = props => {
  const iframeRef = useRef(null)
  const wnapp = useSelector((state: any) => state.apps[props.id])
  const iframeKey = 'my-iframe'
  const iframeSrc = useCallback(() => {
    return wnapp?.url
  }, [wnapp?.url])

  useEffect(() => {
    const iframeElement = iframeRef.current

    const handleKeyDown = event => {
      // 检查按下的按键是否是Ctrl+S
      if (event.ctrlKey && event.key === 's') {
        event.preventDefault()
        iframeElement.contentWindow.postMessage('save', '*') // 向远程桌面发送保存命令
      }

      if (event.ctrlKey && event.key === 't') {
        event.preventDefault()
        iframeElement.contentWindow.postMessage('_blank', '*') // 向远程桌面发送打开新命令
      }
    }

    hotkeys('*', handleKeyDown)

    return () => {
      hotkeys.unbind('*')
    }
  }, [])

  return wnapp ? (
    <div
      className='calcApp floatTab dpShad'
      data-size={wnapp.size}
      id={wnapp.icon + 'App'}
      data-max={wnapp.max}
      style={{
        ...(wnapp.size == 'cstm' ? wnapp.dim : null),
        zIndex: wnapp.z
      }}
      data-hide={wnapp.hide}>
      <ToolBar
        app={wnapp.action}
        icon={wnapp.icon}
        size={wnapp.size}
        // name={wnapp.title}
        name='3D云应用'
      />
      <div
        className='windowScreen flex flex-col'
        data-dock='true'
        key={iframeKey}>
        {!wnapp.hide && wnapp?.className === 'CloudAppWrap_open' ? (
          <iframe
            id='remoteDesktop'
            ref={iframeRef}
            onLoad={() => {}}
            key='CloudAppWrap_open'
            src={iframeSrc()}
            style={{ width: '100%', height: '100%' }}
          />
        ) : (
          <VisList isRefresh={wnapp.hide} />
        )}
      </div>
    </div>
  ) : null
}
