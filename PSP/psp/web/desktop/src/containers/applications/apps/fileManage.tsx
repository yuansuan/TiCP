/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '@/utils/general'
import { Explorer } from './explorer'
export const FileManage = () => {
  const wnapp = useSelector(state => state.apps.explorer)

  return wnapp ? (
    <div
      className='notepad floatTab dpShad'
      data-size={wnapp.size}
      data-max={wnapp.max}
      style={{
        ...(wnapp.size == 'cstm' ? wnapp.dim : null),
        zIndex: wnapp.z
      }}
      data-hide={wnapp.hide}
      id={wnapp.icon + 'App'}>
      <ToolBar
        app={wnapp.action}
        icon={wnapp.icon}
        size={wnapp.size}
        name={wnapp.title}
      />
      <div className='windowScreen flex flex-col' data-dock='true'>
        <Explorer zIndex={wnapp.z} />
      </div>
    </div>
  ) : null
}
