/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '@/utils/general'
import NewJobCreatorByApp from '@/pages/NewJobCreatorByApp'
import { extractPathAndParamsFromURL } from '@/utils'

interface ParamsProps {
  action?: string
  appType?: string
}
export const CalcAppContainers = props => {
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
  const params: ParamsProps = extractPathAndParamsFromURL(currentPath)
  const apps = useSelector((state: any) => state.apps)
  const isTop = apps.hz === apps[props.id].z
  console.log(
    'window: ',
    apps.hz,
    apps[props.id].z,
    apps.hz === apps[props.id].z
  )
  const wnapp = useSelector((state: any) => {
    return state.apps[props.id]
  })
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
        {!wnapp.hide && (
          <NewJobCreatorByApp action={params?.appType} isTop={isTop} />
        )}
      </div>
    </div>
  ) : null
}
