/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState, useEffect, useRef, useCallback } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Tooltip, message, Empty } from 'antd'
import { Modal } from '@/components'
import { Icon } from '@/utils/general'
import { dispatchAction, handleFileOpen } from '@/actions'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { newBoxServer } from '@/server'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { showTextEditor } from '@/components'
import * as hooks from '@/utils/hooks'
import './assets/fileexpo.css'
import { GridView, ListView } from './components'
import { showDirSelector } from '@/components/NewFileMGT/DirSelector'
import SelectedFiles from '@/components/SelectItems'

import {moveTo} from '@/utils'
const server = serverFactory(newBoxServer)
const EDITABLE_SIZE = 3 * 1024 * 1024
const isEditable = size => {
  if (size > EDITABLE_SIZE) {
    return {
      editable: false,
      message: '文件大小超过 3M'
    }
  }

  return {
    editable: true
  }
}

const NavTitle = props => {
  let src = props.icon || 'folder'

  return (
    <div
      className='navtitle flex prtclk'
      data-action={props.action}
      data-payload={props.payload}
      onClick={dispatchAction}>
      <Icon
        className='mr-1'
        src={'win/' + src + '-sm'}
        width={props.isize || 16}
      />
      <span>{props.title}</span>
    </div>
  )
}

const FolderDrop = ({ dir }) => {
  const files = useSelector(state => state.files)
  const folder = files.data.getId(dir)

  return (
    <>
      {folder.data &&
        folder.data.map((item, i) => {
          if (item.type == 'folder') {
            return (
              <Dropdown
                key={i}
                icon={item.info && item.info.icon}
                title={item.name}
                notoggle={item.data.length == 0}
                dir={item.id}
              />
            )
          }
        })}
    </>
  )
}

const Dropdown = props => {
  const [open, setOpen] = useState(props.isDropped != null)
  const special = useSelector(state => state.files.data.special)
  const [fid] = useState(() => {
    if (props.spid) return special[props.spid]
    else return props.dir
  })
  const toggle = () => setOpen(!open)

  return (
    <div className='dropdownmenu'>
      <div className='droptitle'>
        {!props.notoggle ? (
          <Icon
            className='arrUi'
            fafa={open ? 'faChevronDown' : 'faChevronRight'}
            width={10}
            onClick={toggle}
            pr
          />
        ) : (
          <Icon className='arrUi opacity-0' fafa='faCircle' width={10} />
        )}
        <NavTitle
          icon={props.icon}
          title={props.title}
          isize={props.isize}
          action={props.action != '' ? props.action || 'FILEDIR' : null}
          payload={fid}
        />
        {props.pinned != null ? (
          <Icon className='pinUi' src='win/pinned' width={16} />
        ) : null}
      </div>
      {!props.notoggle ? (
        <div className='dropcontent'>
          {open ? props.children : null}
          {open && fid != null ? <FolderDrop dir={fid} /> : null}
        </div>
      ) : null}
    </div>
  )
}

const ContentArea = ({ searchtxt, zIndex, setSelectedFile }) => {
  const files = useSelector(state => state.files)
  const [selected, setSelect] = useState([])
  const [selectedName, setSelectName] = useState([])
  const [selectedObj, setSelectObj] = useState([])
  const fdata = files.data.getId(files.cdir)
  const dispatch = useDispatch()
  const container = useRef<HTMLDivElement>(null)
  const selection = useRef<HTMLDivElement>(null)
  const contentwrap = useRef<HTMLDivElement>(null)
  const [startX, setStartX] = useState(0)
  const [startY, setStartY] = useState(0)
  const [, setEndX] = useState(0)
  const [endY, setEndY] = useState(0)
  const [fileManageWrap, setFileManageWrap] = useState(null)
  const [showHiddenFile, setShowHiddenFile] = useState(true)
  const [changeListView, setChangeListView] = useState(false)
  const rect = hooks.useElementRect(container)
  const [loading, setLoading] = useState(false)
  const [compressPackageEnd, setCompressPackageEnd] = useState(false)

  useEffect(() => {
    const id = setInterval(() => {
       server.getUserCompressStatus().then(res => {
        if (res.data.length === 0) {
          clearInterval(id)
          setCompressPackageEnd(true)
        } else {
          // 失败报错，不用继续定时查询
          // TODO 目前是单任务，代码可以正常运行，支持多个任务，需要修改代码
          const hasFailed = res.data.some(l => l.Status === 1)
          if (hasFailed) {
            message.error('压缩任务失败')
          }
        }
       }).catch((e) => {
        console.log(e)
        clearInterval(id)
       })
     }, 5000)
 
     return () => {
       clearInterval(id)
     }
   }, [])

  useEffect(() => {
    if (rect) {
      const elements = document.querySelectorAll('.ys-filemanger-custom-list')
      if (selected.length > 0) {
        files.isSelectedFile = true
      } else {
        files.isSelectedFile = false
      }
      elements.forEach(element => {
        const wrap = element as HTMLDivElement
        const content = element.querySelector('div') as HTMLDivElement

        if (wrap && content) {
          wrap.setAttribute('data-menu', showMenu(selected))
          content.setAttribute('data-menu', showMenu(selected))
          content.setAttribute('data-parentcontent', 'parentcontent')
          content.setAttribute('data-selected', selected.join(','))
          content.setAttribute('data-selectedname', selectedName.join(','))
          content.setAttribute('data-selectedobj', JSON.stringify(selectedObj))
          wrap.setAttribute('data-parentcontent', 'parentcontent')
          wrap.setAttribute('data-selected', selected.join(','))
          wrap.setAttribute('data-selectedname', selectedName.join(','))
          wrap.setAttribute('data-selectedobj', JSON.stringify(selectedObj))
          if (files?.cdir) {
            wrap.setAttribute('data-currid', files.cdir)
            content.setAttribute('data-currid', files.cdir)
          }
        }
      })
    }
  }, [rect, selected, files?.cdir])



  useEffect(() => {
    // 遍历每个元素，检查是否与选择框相交
    document.querySelectorAll('.conticon').forEach(el => {
      if (selected.includes(el.getAttribute('id'))) {
        el.classList.add('selected')
        el.dataset.focus = 'true'
      }
    })
  },[selected])

  const newFileList =
    fdata?.data?.filter(item =>
      item?.name?.toLowerCase()?.includes(searchtxt?.toLowerCase())
    ) || []

  const handleClick = e => {
    if (!e.shiftKey) {
      e.stopPropagation()
      setSelect([e.target.dataset.id])
      setSelectName([e.target.dataset.name])
      setSelectObj([e.target.dataset])
    }
  }

  const handleDouble = e => {
    // 双击根据id获取
    if (!e.shiftKey) {
      e.stopPropagation()
      e.target.dataset.type !== 'folder' &&
        handleFileDoubleClick(e.target.dataset)
      handleFileOpen(e.target.dataset.id, server)
    }
  }

  const handleFileDoubleClick = async ({ type, path, size, name, ...file }) => {
    if (!isEditable(size).editable) {
      await Modal.showConfirm({
        title: '确认弹窗',
        content: `此文件无法直接预览，是否下载？`
      }).then(async () => {
        await server.download([path], [true], [size])
      })
      return
    }
    switch (true) {
      case /je?pg|png|gif|img|svg/gi.test(type):
        // const url = await server.getFileUrl([path], [true], [size], true)
        // previewImage({ fileName: name, src: url })
        await Modal.showConfirm({
          title: '确认弹窗',
          content: `非文本文件无法预览，是否下载？`
        }).then(async () => {
          await server.download([path], [true], [size])
        })
        return

      case /unknownfile|txt|dmg|mp4|pdf/gi.test(type):
        showTextEditor({
          path,
          fileInfo: {
            size,
            name,
            path,
            type
          },
          readonly: true,
          boxServerUtil: newBoxServer
        })
        return
    }
  }

  const emptyClick = e => {
    e.preventDefault()
    !endY && e.target.dataset?.parentcontent && setSelect([]) // 点击的是非app的空白地方
    if (e.target.dataset?.currid) {
      dispatch({
        type: 'FILEOPERATE',
        payload: 'rename',
        data: e.target.dataset.currid
      })
    }
  }
  const handleDoubleClick = e => {
    e.preventDefault()
  }

  useEffect(() => {
    if (!fileManageWrap) {
      const fileManageWrap = document.getElementById('explorerApp')
      setFileManageWrap(fileManageWrap)
      try {
        setLoading(true)
        server.fetch('.').then(res => {
          dispatch({
            type: 'generateFiles',
            payload: res._children
          })
        })
      } finally {
        setLoading(false)
      }
    }
  }, [])

  useEffect(() => {
    try {
      setLoading(true)
      if (fdata) {
        server.fetch(fdata.path || '.').then(res => {
          dispatch({ type: 'FILEDIR', payload: fdata.id, data: res._children })
        })
      }
    } finally {
      setLoading(false)
    }

    EE.on(EE_CUSTOM_EVENT.SHOW_HIDE_FILE, ({ show }) => {
      setShowHiddenFile(!show)
    })
    EE.on(EE_CUSTOM_EVENT.CHANGE_LIST_VIEW, ({ show }) => {
      setChangeListView(show)
    })
  }, [zIndex, files.cdir, showHiddenFile, changeListView, compressPackageEnd])

  useEffect(() => {
    setSelectedFile(selectedObj)
  }, [selectedObj])
  useEffect(() => {
    setSelect([])
    setSelectName([])
    setSelectObj([])
  }, [files.cdir, changeListView])

  const handleBlur = async (
    e,
    path,
    id,
    name,
    isMkdir,
    cdir,
    currentHistoryPath
  ) => {
    forbidDefaultEvent(e)
    try {
      dispatch({
        type: 'FILEOPERATE',
        payload: 'handleRename',
        data: {
          targetPath: e.target.value,
          originPath: path,
          id,
          name,
          isMkdir
        }
      })

      if (e.target.value) {
        setTimeout(() => {
          server.fetch(currentHistoryPath || '.').then(res => {
            dispatch({
              type: 'FILEDIR',
              payload: cdir,
              data: res._children
            })
          })
        }, 500)
      }
    } catch (err) {}
  }

  const handleMouseDown = e => {
    // 如果同时按下Shift键进行多选或反选
    if (e.shiftKey && !e.target.dataset?.parentcontent) {
      const target = e.target
      if (!target.dataset.focus || target.dataset.focus !== 'true') {
        target.classList.add('selected')
        target.dataset.focus = 'true'
        setSelect([...selected, target.dataset.id])
        setSelectName([...selectedName, target.dataset.name])
        setSelectObj([...selectedObj, target.dataset])
      } else if (target.dataset.focus === 'true') {
        // 反选
        target.classList.remove('selected')
        target.dataset.focus = 'false'
        setSelect(selected.filter(item => item !== target.dataset.id))
        setSelectName(selectedName.filter(item => item !== target.dataset.name))
        setSelectObj(selectedObj.filter(item => item.id !== target.dataset.id))
      }

      return
    } else if (e.button === 0) {
      let startX, startY
      const selectionItem = selection.current
      selectionItem.style.width && (selectionItem.style.width = '0')
      selectionItem.style.height && (selectionItem.style.height = '0')
      startX = e.clientX - 90
      startY = e.clientY - 60
      if (fileManageWrap.style.left || fileManageWrap.dataset.size === 'mini') {
        const left = parseInt(fileManageWrap.style.left) || 300
        startX = e.clientX - left
      }
      if (fileManageWrap.style.top || fileManageWrap.dataset.size === 'mini') {
        const top = parseInt(fileManageWrap.style.top) + 60 || 116
        startY = e.clientY - top
      }
      setStartX(startX)
      setStartY(startY)

      selectionItem.style.display = 'block'
      selectionItem.style.top = startY + 'px'
      selectionItem.style.left = startX + 'px'

      // 左键按下时清除已选择元素的标记
      document.querySelectorAll('.conticon.selected').forEach(el => {
        el.classList.remove('selected')
        el.dataset.focus = 'false'
      })
    }
  }

  const isIntersect = (rect1, rect2) => {
    return !(
      rect2.left > rect1.right ||
      rect2.right < rect1.left ||
      rect2.top > rect1.bottom ||
      rect2.bottom < rect1.top
    )
  }

  const handleMouseMove = e => {
    let endX, endY
    const selectionItem = selection.current
    if (selectionItem.style.display === 'block') {
      endX = e.clientX - 90
      endY = e.clientY - 60
      if (fileManageWrap.style.left || fileManageWrap.dataset.size === 'mini') {
        const left = parseInt(fileManageWrap.style.left) || 300
        endX = e.clientX - left
      }
      if (fileManageWrap.style.top || fileManageWrap.dataset.size === 'mini') {
        const top = parseInt(fileManageWrap.style.top) + 60 || 116
        endY = e.clientY - top
      }
      setEndX(endX)
      setEndY(endY)

      selectionItem.style.width = Math.abs(endX - startX) + 'px'
      selectionItem.style.height = Math.abs(endY - startY) + 'px'
      selectionItem.style.top = Math.min(endY, startY) + 'px'
      selectionItem.style.left = Math.min(endX, startX) + 'px'

      // 遍历每个元素，检查是否与选择框相交
      document.querySelectorAll('.conticon').forEach(el => {
        const rect = el.getBoundingClientRect()
        if (isIntersect(rect, selectionItem.getBoundingClientRect())) {
          el.classList.add('selected')
          el.dataset.focus = 'true'
        } else {
          el.dataset.focus = 'false'
          el.classList.remove('selected')
        }
      })
    }
  }

  const handleMouseUp = e => {
    // 收集选中的id
    const collectId = []
    const collectName = []
    const collectPath = []
    const collectObj = []
    document.querySelectorAll('.conticon').forEach(el => {
      if (el.dataset.focus === 'true') {
        collectId.push(el.dataset.id)
        collectName.push(el.dataset.name)
        // collectPath.push(el.dataset.path)
        collectObj.push(el.dataset)
      }
    })
    setSelect(collectId)
    setSelectName(collectName)
    // setSelectPath(collectPath)
    setSelectObj(collectObj)
    // gridshow.current.dataset.menu = 'rightMenu'
    // contentwrap.current.dataset.menu = 'rightMenu'
    selection.current.style.display = 'none'
  }

  // 阻止默认事件Function
  const forbidDefaultEvent = e => {
    e.preventDefault()
    e.stopPropagation()
  }
  const collectFiles = files => {
    let fileArray = []
    for (let i = 0; i < files.length; i++) {
      const entry = files[i].webkitGetAsEntry()
      if (entry.isDirectory) {
        const reader = entry.createReader()
        reader.readEntries(en => {
          en.forEach(item => {
            if (item.isDirectory) {
              collectFiles(item)
            } else {
              fileArray.push(item)
            }
          })
        })
      } else {
        fileArray.push(entry)
      }
    }
    return fileArray
  }

  const showMenu = useCallback(
    selected => {
      if (selected?.length === 0) {
        return 'noFile'
      } else if (selected?.length === 1 && selected[0] !== '') {
        return 'singleFile'
      } else {
        return 'multipleFiles'
      }
    },
    [selected]
  )

  return (
    <div
      className='contentarea'
      onClick={emptyClick}
      onDoubleClick={handleDoubleClick}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      data-parentcontent='parentcontent'
      ref={container}
      tabIndex={-1}
      data-menu={showMenu(selected)}>
      <div
        className='contentwrap win11Scroll'
        ref={contentwrap}
        data-selectedobj={JSON.stringify(selectedObj)}
        data-menu={showMenu(selected)}>
        {rect && (
          <div
            className={changeListView ? 'gridshow' : 'listshow'}
            data-size='lg'
            data-currid={files?.cdir}
            data-parentcontent='parentcontent'
            data-selected={selected.join(',')}
            data-selectedname={selectedName.join(',')}
            data-selectedobj={JSON.stringify(selectedObj)}
            data-menu={showMenu(selected)}>
            {changeListView ? (
              <GridView
                rect={rect}
                loading={loading}
                items={newFileList}
                selected={selected}
                selectedName={selectedName}
                selectedObj={selectedObj}
                handleClick={handleClick}
                handleBlur={handleBlur}
                handleDouble={handleDouble}
                changeListView={changeListView}
              />
            ) : (
              <ListView
                rect={rect}
                loading={loading}
                items={newFileList}
                selected={selected}
                selectedName={selectedName}
                selectedObj={selectedObj}
                handleClick={handleClick}
                handleBlur={handleBlur}
                handleDouble={handleDouble}
                changeListView={changeListView}
              />
            )}
          </div>
        )}
        {newFileList?.length == 0 ? (
          <span className='text-xs mx-auto emptyView'>
            <Empty />
          </span>
        ) : null}
      </div>
      <div id='selection' ref={selection}></div>
    </div>
  )
}

const NavPane = ({}) => {
  const files = useSelector(state => state.files)
  const special = useSelector(state => state.files.data.special)

  return (
    <div className='navpane win11Scroll'>
      <div className='extcont'>
        <Dropdown icon='star' title='Quick access' action='' isDropped>
          <Dropdown
            icon='down'
            title='Downloads'
            spid='%downloads%'
            notoggle
            pinned
          />
          <Dropdown icon='user' title='Blue' spid='%user%' notoggle pinned />
          <Dropdown
            icon='docs'
            title='Documents'
            spid='%documents%'
            notoggle
            pinned
          />
          <Dropdown title='Github' spid='%github%' notoggle />
          <Dropdown icon='pics' title='Pictures' spid='%pictures%' notoggle />
        </Dropdown>
        <Dropdown icon='onedrive' title='OneDrive' spid='%onedrive%' />
        <Dropdown icon='thispc' title='This PC' action='' isDropped>
          <Dropdown icon='desk' title='Desktop' spid='%desktop%' />
          <Dropdown icon='docs' title='Documents' spid='%documents%' />
          <Dropdown icon='down' title='Downloads' spid='%downloads%' />
          <Dropdown icon='music' title='Music' spid='%music%' />
          <Dropdown icon='pics' title='Pictures' spid='%pictures%' />
          <Dropdown icon='vid' title='Videos' spid='%videos%' />
          <Dropdown icon='disc' title='OS (C:)' spid='%cdrive%' />
          <Dropdown icon='disk' title='Blue (D:)' spid='%ddrive%' />
        </Dropdown>
      </div>
    </div>
  )
}

const Ribbon = ({ selectedFile = [], fdata }) => {
  const dispatch = useDispatch()
  const [toggle, setToggle] = useState(true)
  const [toggleView, setToggleView] = useState(true)
  EE.on(EE_CUSTOM_EVENT.SHOW_HIDE_FILE, ({ show }) => {
    setToggle(!show)
  })
  const freshFileList = async (path, id) => {
    const fileList = await server.fetch(path)
    dispatch({ type: 'FILEDIR', payload: id, data: fileList._children })
  }
  const upload = payload => {
    dispatch({
      type: 'UPLOAD',
      payload,
      callback: freshFileList
    })
  }
  const handleCreateFile = () => {
    EE.emit(EE_CUSTOM_EVENT.SETGRIDLISTSCROLLTOEND, true)

    dispatch({
      type: 'CREATEFILEDIR',
      payload: 'createFile'
    })
  }

  const checkReadOnly = async (action,selectedFile) => {
    if (!selectedFile.length) return
    const readOnlyFiles = selectedFile.filter(item => item.onlyread === 'true') || []
    // 处理不同的操作
    let actionText = ''
    let actionWarning = ''
    if (action === 'move') {
      actionText = '移动'
      actionWarning = '移动文件'
    } else if (action === 'remove') {
      actionText = '删除'
      actionWarning = '文件删除'
    } else if(action === 'download'){
      actionText = '下载',
      actionWarning = '文件下载'
    }
    if (readOnlyFiles.length > 0) {
      Modal.warn({
        title: actionWarning,
        content: readOnlyFiles.length <= 1 ? (
          ` ${
            readOnlyFiles[0].name
          } 是系统内置目录，不能被${actionText}！`
        ) : (
          <div>
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
  const deleteFiles =async () => {
    const check = await checkReadOnly('remove',selectedFile)
    const selectFileName = selectedFile.map(item => item.name)

    if(check){
      Modal.confirm({
        title: '文件删除',
        content: <SelectedFiles selectedName={selectFileName}/>,
        okText: '确认',
        cancelText: '取消',
        onOk: () => {
          const deletePath = selectedFile.map(item => item.path)
          server.delete(deletePath).then(res => {})
          freshFileList(fdata?.path || '.', fdata?.id)
          message.success('删除成功')
        }
      })
    }
  }

  const downloadFile = async () => {
    if (!selectedFile.length) return
    const path = selectedFile.map(item => item.path)
    const isFile = selectedFile.map(item => item.isfile === 'true')
    const size = selectedFile.map(item => item.size)

    const check = await checkReadOnly('download',selectedFile)
    if(check){
      server.download(path, isFile, size)
    }
  }


  const moveToFolder = async () => {
    const check = await checkReadOnly('move',selectedFile)
    if(check){
      showDirSelector({
        disabledPaths: selectedFile.map(item => item.path)
      }).then(async path => {
       await moveTo(path, selectedFile)
       await freshFileList(fdata?.path || '.', fdata?.id)
      })
    }
  }

  return (
    <div className='msribbon flex'>
      <div className='ribsec'>
        <Tooltip title='新建文件夹'>
          <div className='drdwcont flex' onClick={handleCreateFile}>
            <Icon src='mkdir' ui width={18} margin='0 6px' />
          </div>
        </Tooltip>
      </div>
      <div className='ribsec'>
        {/* <Icon src='cut' ui width={18} margin='0 6px' />
        <Icon src='copy' ui width={18} margin='0 6px' />
        <Icon src='paste' ui width={18} margin='0 6px' />
        <Icon src='rename' ui width={18} margin='0 6px' />
        <Icon src='share' ui width={18} margin='0 6px' /> */}
        <Tooltip title='上传文件'>
          <div className=' drdwcont flex upload_wrap' onClick={() => upload(false)}>
            <Icon src='uploadFile' ui width={18} margin='0 6px' />
          </div>
        </Tooltip>
        <Tooltip title='上传文件夹'>
          <div className=' drdwcont flex upload_wrap' onClick={() => upload(true)}>
            <Icon src='uploadFolder' ui width={18} margin='0 6px' />
          </div>
        </Tooltip>
      </div>
      <div className='ribsec'>
      <Tooltip title={toggle ? '显示所有文件' : '过滤隐藏文件'}>
      <div
          className='drdwcont flex'
          onClick={() => {
            setToggle(!toggle)
            EE.emit(EE_CUSTOM_EVENT.SHOW_HIDE_FILE, { show: toggle })
          }}>
          <Icon src='displayFiles' ui width={18} margin='0 6px' />
        </div>
      </Tooltip>
      </div>
      <div className='ribsec'>
        <Tooltip title={toggleView ? '大图标' : '列表视图'}>
        <div
          className='drdwcont flex'
          onClick={() => {
            setToggleView(!toggleView)
            EE.emit(EE_CUSTOM_EVENT.CHANGE_LIST_VIEW, { show: toggleView })
          }}>
          <Icon
            src={toggleView ? 'folderview' : 'listview'}
            ui
            width={18}
            margin='0 6px'
          />
        </div>
      </Tooltip>

      </div>
      {selectedFile?.length > 0 && (
        <div className='ribsec'>
          <Tooltip title='下载'>
            <div className='drdwcont flex' onClick={downloadFile}>
              <Icon src='download' ui width={18} margin='0 6px' />
            </div>
          </Tooltip>

          <Tooltip title='删除'>
            <div className='drdwcont flex' onClick={deleteFiles}>
              <Icon src='delete' ui width={18} margin='0 6px' />
            </div>
          </Tooltip>
          <Tooltip title='移动'>
            <div className='drdwcont flex' onClick={moveToFolder}>
              <Icon src='move' ui width={18} margin='0 6px' />
            </div>
          </Tooltip>
        </div>
      )}
    </div>
  )
}

export const Explorer = ({ zIndex }: { zIndex: number | string }) => {
  const files = useSelector(state => state.files)
  const fdata = files.data.getId(files.cdir)
  const [cpath, setPath] = useState(files.cpath)
  const [searchtxt, setShText] = useState('')
  const [selectedFile, setSelectedFile] = useState([])
  const dispatch = useDispatch()

  const handleChange = e => setPath(e.target.value)
  const handleSearchChange = e => setShText(e.target.value)

  const newFileList =
    fdata?.data?.filter(item =>
      item?.name?.toLowerCase()?.includes(searchtxt?.toLowerCase())
    ) || []

  const handleEnter = e => {
    if (e.key === 'Enter') {
      dispatch({ type: 'FILEPATH', payload: cpath })
    }
  }

  const DirCont = () => {
    let arr = [],
      curr = fdata,
      index = 0

    while (curr) {
      arr.push(
        <div key={index++} className='dirCont flex items-center'>
          <div
            className='dncont'
            onClick={dispatchAction}
            tabIndex={-1}
            data-action='FILEDIR'
            data-payload={curr.id}>
            {curr.name}
          </div>
          <Icon className='dirchev' fafa='faChevronRight' width={8} />
        </div>
      )

      curr = curr.host
    }

    arr.push(
      <div key={index++} className='dirCont flex items-center'>
        <Icon
          className='pr-1 pb-px'
          src={'win/' + fdata?.info?.icon + '-sm'}
          width={16}
        />
        <Icon className='dirchev' fafa='faChevronRight' width={8} />
      </div>
    )

    return (
      <div key={index++} className='dirfbox h-full flex'>
        {arr.reverse()}
      </div>
    )
  }

  useEffect(() => {
    setPath(files.cpath)
    setShText('')
  }, [files.cpath])

  return (
    <>
      <Ribbon selectedFile={selectedFile} fdata={fdata} />
      <div className='restWindow flex-grow flex flex-col'>
        <div className='sec1'>
          <Icon
            className={'navIcon hvtheme' + (files.hid == 0 ? ' disableIt' : '')}
            fafa='faArrowLeft'
            width={14}
            click='FILEPREV'
            pr
          />
          <Icon
            className={
              'navIcon hvtheme' +
              (files.hid + 1 == files.hist.length ? ' disableIt' : '')
            }
            fafa='faArrowRight'
            width={14}
            click='FILENEXT'
            pr
          />
          <Icon
            className='navIcon hvtheme'
            fafa='faArrowUp'
            width={14}
            click='FILEBACK'
            pr
          />
          <div className='path-bar noscroll' tabIndex={-1}>
            <input
              className='path-field'
              type='text'
              value={cpath}
              onChange={handleChange}
              onKeyDown={handleEnter}
            />
            <DirCont />
          </div>
          <div className='srchbar'>
            <Icon className='searchIcon' src='search' width={12} />
            <input
              type='text'
              onChange={handleSearchChange}
              value={searchtxt}
              placeholder='Search'
            />
          </div>
        </div>
        <div className='sec2'>
          {/* <NavPane /> */}
          <ContentArea
            searchtxt={searchtxt}
            zIndex={zIndex}
            setSelectedFile={setSelectedFile}
          />
        </div>
        <div className='sec3'>
          <div className='item-count text-xs'>
            {newFileList?.length || 0} items
          </div>
        </div>
      </div>
    </>
  )
}
