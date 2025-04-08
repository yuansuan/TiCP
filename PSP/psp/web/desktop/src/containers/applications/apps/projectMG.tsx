/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '../../../utils/general'
import ProjectMGPage from '@/pages/Admin/ProjectMG'

export const ProjectMG = () => {
  const apps = useSelector(state => state.apps)
  const wnapp = useSelector(state => state.apps.projectMG)

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
        {!wnapp.hide && <ProjectMGPage entry='desktop' refresh={apps.hz === apps.projectMG.z} />}
      </div>
    </div>
  ) : null
}
