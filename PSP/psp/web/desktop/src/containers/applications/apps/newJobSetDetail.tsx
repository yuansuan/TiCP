/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '../../../utils/general'
import NewJobSetDetail from '@/pages/NewJobSetDetail'

export const JobSetDetail = () => {
  const wnapp = useSelector(state => state.apps.jobSetDetail)

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
        name={wnapp.title}
      />
      <div className='windowScreen flex flex-col' data-dock='true'>
        {!wnapp.hide && <NewJobSetDetail />}
      </div>
    </div>
  ) : null
}
