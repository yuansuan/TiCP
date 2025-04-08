/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { Modal, message, Tooltip } from 'antd'
import { useSelector, useDispatch } from 'react-redux'
import { Icon } from '../../utils/general'
import './menu.css'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { Http } from '@/utils'

import * as Actions from '../../actions'
import { history } from '@/utils'
import { vis, env, sysConfig } from '@/domain'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { newBoxServer } from '@/server'
import SelectedFiles from '@/components/SelectItems'
const server = serverFactory(newBoxServer)

export const ActMenu = () => {
  const menu = useSelector(state => state.menus)
  const menudata = menu.data[menu.opts]
  const dataSet = menu.dataset
  const selected = dataSet?.selected
    ? dataSet?.selected?.split(',')
    : [dataSet?.id]
  const selectedName = dataSet?.selectedname?.split(',') || [dataSet?.name]

  const checkReadOnly = async (action, dataSet) => {
    const readOnlyFiles = JSON.parse(dataSet?.selectedobj || '[]').filter(
      item => item.onlyread === 'true'
      )

    // 处理不同的操作
    let actionText = ''
    let actionWarning = ''
    if (action === 'move') {
      actionText = '移动'
      actionWarning = '移动文件'
    } else if (action === 'rename') {
      actionText = '重命名'
      actionWarning = '重命名文件'
    } else if (action === 'share') {
      actionText = '分享'
      actionWarning = '分享文件'
    } else if (action === 'send') {
      actionText = '发送'
      actionWarning = '发送文件'
    } else if (action === 'remove') {
      actionText = '删除'
      actionWarning = '文件删除'
    } else if (action === 'download') {
      actionText = '下载'
      actionWarning = '文件下载'
    }

    if (dataSet.onlyread === 'true' || readOnlyFiles.length > 0) {
      Modal.warn({
        title: actionWarning,
        content:
          selected.length <= 1 ? (
            ` ${dataSet?.name} 是系统内置目录，不能被${actionText}！`
          ) : (
            <div style={{maxHeight: '500px',overflowY: 'auto'}}>
              <p>以下目录是系统内置目录，不能被{actionText}！</p>
              <ul style={{ marginLeft: 20 }}>
                {readOnlyFiles.map((item, index) => (
                  <li key={index}>{item.name}</li>
                ))}
              </ul>
            </div>
          ),
        okText: '关闭'
      })

      return false
    }

    return true
  }

  const { abpos, isLeft } = useSelector(state => {
    let acount = state.menus.menus[state.menus.opts].length
    let tmpos = {
        top: state.menus.top,
        left: state.menus.left
      },
      tmpleft = false

    let wnwidth = window.innerWidth,
      wnheight = window.innerHeight

    let ewidth = 312,
      eheight = acount * 28

    tmpleft = wnwidth - tmpos.left > 504
    if (wnwidth - tmpos.left < ewidth) {
      tmpos.left = wnwidth - ewidth
    }

    if (wnheight - tmpos.top < eheight) {
      tmpos.bottom = wnheight - tmpos.top
      tmpos.top = null
    }

    return {
      abpos: tmpos,
      isLeft: tmpleft
    }
  })

  const dispatch = useDispatch()

  const uploadCallBack = async (path, id) => {
    const fileList = await server.fetch(path)
    dispatch({ type: 'FILEDIR', payload: id, data: fileList._children })
  }

  const fileCompressionCallBack = async (path, id) => {
    const fileList = await server.fetch(path)
    dispatch({ type: 'FILEDIR', payload: id, data: fileList._children })
    
    // 更新一下 压缩 菜单的状态
    dispatch({
      type: 'MENUITEMUPDATE',
      payload: {
        menu: 'singleFile',
        menuItemPayload: 'compress',
        menuItemAttr: {
          disabled: false,
          busy: false,
          tooltip: undefind
        }
      }
    })

    dispatch({
      type: 'MENUITEMUPDATE',
      payload: {
        menu: 'multipleFiles',
        menuItemPayload: 'compress',
        menuItemAttr: {
          disabled: false,
          busy: false,
          tooltip: undefind
        }
      }
    })
  }

  const clickDispatch = async event => {
    if (event.target.dataset.disabled === 'true') return 
    event.stopPropagation()
    event.preventDefault()
    let action = {
      type: event.target.dataset.action,
      payload: event.target.dataset.payload
    }
    if (dataSet.routerpath) {
      // history.push(dataSet.routerpath)
      window.localStorage.setItem('CURRENTROUTERPATH', dataSet.routerpath)
    }

    if (action.type) {
      if (action.payload === 'delshort') {
        Modal.confirm({
          title: '关闭3D云应用',
          content: '确认关闭！',
          okText: '确认',
          cancelText: '取消',
          onOk: () => {
            // 3D云应用删除逻辑
            if (dataSet.sessionid) {
              vis
                .closeSession(dataSet.sessionid)
                .then(() => {
                  message.success('删除会话成功，会话即将删除...')
                  dispatch({
                    type: 'desktop',
                    data: {
                      sessionId: dataSet.sessionid,
                      action: dataSet.icon
                    },
                    payload: 'closeSession'
                  })
                })
                .catch(() => {
                  message.error('删除失败！')
                })
            } else {
              if (action.type != action.type.toUpperCase()) {
                Actions[action.type](action.payload, menu)
              } else {
                dispatch(action)
              }
            }
          }
        })
      } else {
        if (action.payload === 'createFile') {
          EE.emit(EE_CUSTOM_EVENT.SETGRIDLISTSCROLLTOEND, true)
          dispatch({
            type: 'CREATEFILEDIR',
            payload: 'createFile',
            data: dataSet.currid
          })
        } else if (action.payload === 'rename') {
          const readOnlyCheck = await checkReadOnly('rename', dataSet)
          if (readOnlyCheck) {
            dispatch({
              type: 'FILEOPERATE',
              payload: 'rename',
              data: dataSet.id
            })
          }
        } else if (action.payload === 'remove') {
          const readOnlyCheck = await checkReadOnly('remove', dataSet)
          if (readOnlyCheck) {
            Modal.confirm({
              title: '文件删除',
              content: <SelectedFiles selectedName={selectedName}/>,
              okText: '确认',
              cancelText: '取消',
              onOk: () => {
                dispatch({
                  type: 'FILEREMOVE',
                  payload: 'remove',
                  data: { path: dataSet.path, id: dataSet.id, selected }
                })
              }
            })
          }
        } else if(action.payload === 'compress') {
          const selected = dataSet.selected ? dataSet.selected?.split(',') : [dataSet.id]
          dispatch({
            type: 'FILECOMPRESSION',
            payload: 'compress',
            data: { path: dataSet.path, id: dataSet.id, selected },
            callback: fileCompressionCallBack,
          })
        } else if (action.payload === 'uploadFile') {
          dispatch({
            type: 'UPLOAD',
            payload: false,
            callback: uploadCallBack
          })
        } else if (action.payload === 'uploadFiles') {
          dispatch({
            type: 'UPLOAD',
            payload: true,
            callback: uploadCallBack
          })
        } else if (action.payload === 'download') {
          const readOnlyCheck = await checkReadOnly('download', dataSet)
          if (readOnlyCheck) {
            const action = {
              type: 'DOWNLOAD',
              payload: dataSet.id || dataSet.currid
            }
            if (dataSet.currid) {
              const selectObj = JSON.parse(dataSet?.selectedobj || '[]')
              action.data = selectObj
            }
            dispatch(action)
          }
        } else if (action.payload === 'move') {
          const readOnlyCheck = await checkReadOnly('move', dataSet)
          if (readOnlyCheck) {
            dispatch({
              type: 'FIlESMOVE',
              payload: dataSet.currid
                ? JSON.parse(dataSet?.selectedobj || '[]')
                : [dataSet],
              callback: uploadCallBack
            })
          }
        } else if (action.payload === 'share') {
          const readOnlyCheck = await checkReadOnly('share', dataSet)
          if (readOnlyCheck) {
            dispatch({
              type: 'FIlESSHARE',
              payload: dataSet.currid
                ? JSON.parse(dataSet?.selectedobj || '[]')
                : [dataSet],
              actType: 'share'
            })
          }
        } else if (action.payload === 'send') {
          const readOnlyCheck = await checkReadOnly('send', dataSet)
          if (readOnlyCheck) {
            dispatch({
              type: 'FIlESSHARE',
              payload: dataSet.currid
                ? JSON.parse(dataSet?.selectedobj || '[]')
                : [dataSet],
              actType: 'send'
            })
          }
        } else if (action.type != action.type.toUpperCase()) {
          Actions[action.type](action.payload, menu)
        } else {
          dispatch(action)
        }
      }
      dispatch({ type: 'MENUHIDE' })
    }
  }

  const menuobj = data => {
    let mnode = []
    data.map((opt, i) => {
      if (opt.type == 'hr') {
        mnode.push(<div key={i} className='menuhr'></div>)
      } else {
        const div = (<div
            key={i}
            className='menuopt'
            data-dsb={opt.dsb}
            data-disabled={opt?.disabled}
            onClick={clickDispatch}
            data-action={opt.action}
            data-payload={opt.payload}>
            {menudata.ispace != false ? (
              <div className='spcont'>
                {opt.icon && opt.type == 'svg' ? (
                  <Icon icon={opt.icon} width={16} />
                ) : null}
                {opt.icon && opt.type == 'fa' ? (
                  <Icon fafa={opt.icon} width={16} />
                ) : null}
                {opt.icon && opt.type == null ? (
                  <Icon src={opt.icon} width={16} />
                ) : null}
              </div>
            ) : null}
            <div className='nopt'>{opt.name}</div>
            {opt.opts ? (
              <Icon
                className='micon rightIcon'
                fafa='faChevronRight'
                width={10}
                color='#999'
              />
            ) : null}
            {opt.dot ? (
              <Icon
                className='micon dotIcon'
                fafa='faCircle'
                width={4}
                height={4}
              />
            ) : null}
            {opt.check ? (
              <Icon
                className='micon checkIcon'
                fafa='faCheck'
                width={8}
                height={8}
              />
            ) : null}
            {opt.busy ? (
              <Icon
                className='micon checkIcon'
                fafa='faSpinner'
                spin
                width={12}
                height={12}
              />
            ) : null}
            {opt.opts ? (
              <div
                className='minimenu'
                style={{
                  minWidth: menudata.secwid
                }}>
                {menuobj(opt.opts)}
              </div>
            ) : null}
          </div>)

        const tipsDiv = <Tooltip title={opt.tooltip || ''}>{div}</Tooltip>

        mnode.push(opt.tooltip ? tipsDiv : div)
      }
    })

    return mnode
  }

  return (
    <div
      className='actmenu'
      id='actmenu'
      style={{
        ...abpos,
        '--prefix': 'MENU',
        width: menudata.width,
        backgroundColor: menudata.backgroundColor
      }}
      data-hide={menu.hide}
      data-left={isLeft}>
      {menuobj(menu.menus[menu.opts])}
    </div>
  )
}

export default ActMenu
