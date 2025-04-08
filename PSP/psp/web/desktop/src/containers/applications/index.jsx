/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { useSelector } from 'react-redux'
import './tabs.css'
import './tabs2.css'
import './wnapp.css'

// 导出应用
export * from './apps/calculator'
export * from './apps/fileManage'
export * from './apps/jobManage'
export * from './apps/newJobDetail'
export * from './apps/newJobSetDetail'
export * from './apps/cloudApp'
export * from './apps/starCcmPlusCalc'
export * from './apps/enterpriseManage'
export * from './apps/messages'
export * from './apps/fileExplorer'
export * from './apps/dashboard'
export * from './apps/userlog'
export * from './apps/report'
export * from './apps/projectMG'
export * from './apps/auditLog'
export * from './apps/securityApproval'

import { CalcAppContainers } from './apps/CalcAppContainers'
import { CloudAppContainers } from './apps/CloudAppContainers'

export const ScreenPreview = () => {
  const tasks = useSelector(state => state.taskbar)

  return (
    <div className='prevCont' style={{ left: tasks.prevPos + '%' }}>
      <div className='prevScreen' id='prevApp' data-show={tasks.prev && false}>
        <div id='prevsc'></div>
      </div>
    </div>
  )
}

export const AllCloudApps = () => {
  const apps = useSelector(state => state.apps)
  const collectCloudApp = {}
  for (let key in apps) {
    if (apps[key].renderType && apps[key].renderType === 'cloudApp') {
      collectCloudApp[key] = apps[key]
    }
  }

  return Object.keys(collectCloudApp).map(item => {
    const currentAppInfo = collectCloudApp[item]
    return <CloudAppContainers id={item} key={item} />
  })
}

export const AllCalcApps = () => {
  const apps = useSelector(state => state.apps)
  const collectCalcApp = {}
  for (let key in apps) {
    if (apps[key].renderType && apps[key].renderType === 'calcApp') {
      collectCalcApp[key] = apps[key]
    }
  }

  return Object.keys(collectCalcApp).map(item => {
    const currentAppInfo = collectCalcApp[item]
    return <CalcAppContainers id={item} key={item} />
  })
}
