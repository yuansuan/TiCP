/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { useDispatch, useSelector } from 'react-redux'
import './indexDesktop.css'
import ActMenu from './components/menu'
import {
  BandPane,
  CalnWid,
  DesktopApp,
  SidePane,
  StartMenu,
  WidPane
} from './components/start'
import Taskbar from './components/taskbar'
import { Background } from './containers/background'
import { Uploader } from '@/components/DrawerUploader'
import VerticalDragButton from '@/components/VerticalDragButton'
import * as Applications from './containers/applications'
import * as Drafts from './containers/applications/draft'
import Page500 from '@/pages/500'
import { fetchSoftware } from '@/domain/refreshDesktop'
import { nextTick } from '@/utils'
import { RouterTogg } from '@/constant'
import { extractPathAndParamsFromURL } from '@/utils'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { newBoxServer } from '@/server'
import SysConfig from '@/domain/SysConfig'
import { currentUser } from '@/domain'

const threeMemberApps = ['SecurityApproval', 'AuditLog']

const server = serverFactory(newBoxServer)

function ErrorFallback({ error, resetErrorBoundary }) {
  return (
    <div>
      <Page500 description={error.message} />
    </div>
  )
}

function App() {
  const apps = useSelector(state => state.apps)
  const dispatch = useDispatch()
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH') || ''
  const afterMath = event => {
    let ess = [
      ['START', 'STARTHID'],
      ['BAND', 'BANDHIDE'],
      ['PANE', 'PANEHIDE'],
      ['WIDG', 'WIDGHIDE'],
      ['CALN', 'CALNHIDE'],
      ['MENU', 'MENUHIDE']
    ]

    let actionType = ''
    try {
      actionType = event.target.dataset.action || ''
    } catch (err) {}

    let actionType0 = getComputedStyle(event.target).getPropertyValue(
      '--prefix'
    )

    ess.forEach((item, i) => {
      if (!actionType.startsWith(item[0]) && !actionType0.startsWith(item[0])) {
        dispatch({
          type: item[1]
        })
      }
    })
  }

  function removeFilesFocus() {
    document.querySelectorAll('.conticon').forEach(el => {
      if (el.getAttribute('data-focus') !== null) {
        el.dataset.focus = 'false'
        el.classList.remove('selected')
      }
    })
  }

  const updateMenuItemByPayload = (menuName, menuItemPayload, attrs) => {
    dispatch({
      type: 'MENUITEMUPDATE',
      payload: {
        menu: menuName,
        menuItemPayload: menuItemPayload,
        menuItemAttr: attrs
      }
    })
  }

  window.oncontextmenu = async e => {
    afterMath(e)
    e.preventDefault()
    let data = {
      top: e.clientY,
      left: e.clientX
    }

    if (e.target.dataset.menu != null) {
      data.menu = e.target.dataset.menu
      data.attr = e.target.attributes
      data.dataset = e.target.dataset
      // if(e.target.dataset.size !== null && e.target.dataset.focus !== null && e.target.dataset.path !== null){
      //   e.target.classList.add('selected')
      // }else {
      //   e.target.classList.remove('selected')
      // }

      const checkUserPackageStatus = async () => {
        try {
          const res = await server.getUserCompressStatus()
          if (res.data.length !== 0) {
            const hasRunning = res.data.some(l => l.Status === 2)
            if (hasRunning) {
              updateMenuItemByPayload(data.menu, 'compress', {
                disabled: true,
                busy: true,
                tooltip: '有正在执行的压缩任务'
              })
            } else {
              // 列表无论是成功的 3 还是失败的 1 都是任务结束状态
              updateMenuItemByPayload(data.menu, 'compress', {
                disabled: false,
                busy: false,
                tooltip: ''
              })
            }
          } else {
            updateMenuItemByPayload(data.menu, 'compress', {
              disabled: false,
              busy: false,
              tooltip: ''
            })
          }
        } catch (e) {
          console.log(e)
        }
      }

      if (
        data.menu === 'singleFile' ||
        data.menu === 'multipleFiles' ||
        data.menu === 'fileMenu'
      ) {
        await checkUserPackageStatus()
      }

      dispatch({
        type: 'MENUSHOW',
        payload: data
      })
    }
  }

  window.onclick = afterMath

  window.onload = e => {
    dispatch({ type: 'WALLBOOTED' })
  }

  useEffect(() => {
    if (!window.onstart) {
      // loadSettings()
      window.onstart = setTimeout(() => {
        dispatch({ type: 'WALLBOOTED' })
      }, 5000)
    }
  })

  useEffect(() => {
    fetchSoftware().finally(() => {
      nextTick(() => {
        // 获取刷新之后路由pathname，进行togg
        const formatPathParams = extractPathAndParamsFromURL(currentPath)
        const taskType = RouterTogg[formatPathParams?.path]
        if (formatPathParams?.path) {
          if (formatPathParams?.path === 'new-job-creator') {
            setTimeout(() => {
              dispatch({
                type:
                  formatPathParams?.appType + formatPathParams?.id || taskType,
                payload: 'togg'
              })
            }, 100)
          } else if (formatPathParams?.path === 'vis-session') {
            setTimeout(() => {
              dispatch({
                type: formatPathParams?.actionApp || taskType,
                payload: 'togg'
              })
            }, 100)
          } else if (formatPathParams?.path === 'new-jobs') {
            setTimeout(() => {
              dispatch({
                type: formatPathParams?.actionApp || taskType,
                payload: 'togg'
              })
            }, 100)
          } else if (
            (formatPathParams?.path === 'new-job' &&
              formatPathParams?.jobId &&
              formatPathParams?.isCloud !== null) ||
            (formatPathParams?.path === 'new-job-set' &&
              formatPathParams?.jobSetId &&
              formatPathParams?.isCloud !== null)
          ) {
            setTimeout(() => {
              dispatch({
                type: formatPathParams?.actionApp || taskType,
                payload: 'full'
              })
            }, 100)
          } else {
            dispatch({
              type: taskType,
              payload: 'togg'
            })
          }
        }
      })
    })
  }, [])

  return (
    <div className='App'>
      <ErrorBoundary FallbackComponent={ErrorFallback}>
        <div className='appwrap'>
          <Background />
          <div className='desktop' data-menu='desk'>
            <Uploader />
            <DesktopApp 
              changeApp={(app) => {
                // Hard Code for SecurityApproval
                if (app.name === 'SecurityApproval') {
                  app.title = currentUser.hasSecurityApprovalPerm ? app.title : '审批申请'
                  return app
                } else {
                  return app
                }
              }}
              filterApp={(appName) => (SysConfig.enableThreeMemberMgr ? true : !threeMemberApps.includes(appName))}/>
            <VerticalDragButton />
            {Object.keys(Applications).map((key, idx) => {
              let WinApp = Applications[key]
              return <WinApp key={idx} />
            })}
            {Object.keys(apps)
              .filter(x => x != 'hz')
              .map(key => apps[key])
              .map((app, i) => {
                if (app.pwa) {
                  let WinApp = Drafts[app.data.type]
                  return <WinApp key={i} icon={app.icon} {...app.data} />
                }
              })}
            <StartMenu />
            <BandPane />
            <SidePane />
            <WidPane />
            <CalnWid />
          </div>
          <Taskbar />
          <ActMenu />
        </div>
      </ErrorBoundary>
    </div>
  )
}

export default App
