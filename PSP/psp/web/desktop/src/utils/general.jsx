/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState, useEffect } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { useSelector, useDispatch } from 'react-redux'
import './general.css'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import * as FaIcons from '@fortawesome/free-solid-svg-icons'
import * as FaRegIcons from '@fortawesome/free-regular-svg-icons'
import * as AllIcons from './icons'
import history from './history'
import { checkImgExists } from '.'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

const defaultAppIcon = require('@/assets/images/defaultApp.svg')
String.prototype.strip = function (c) {
  let i = 0,
    j = this.length - 1
  while (this[i] === c) i++
  while (this[j] === c) j--
  return this.slice(i, j + 1)
}

String.prototype.count = function (c) {
  let result = 0,
    i = 0
  for (i; i < this.length; i++) if (this[i] == c) result++
  return result
}

export const Icon = props => {
  const dispatch = useDispatch()
  const [iconSrc, setIconSrc] = useState('')
  let src = ''
  // 除了求解和cloud应用软件icon取接口，其他的本地桌面应用取本地资源
  if (props.ext != null || (props.src && props.src.includes('base64'))) {
    // 判断src是否有效
    checkImgExists(props.src).then(async ({ success }) => {
      if (success) {
        // 图片有效
        await setIconSrc(props.src)
      }
    })
  } else {
    src = `/img/icon/${props.ui != null ? 'ui/' : ''}${
      props.src?.includes('CLOUDAPP') ? '3dcloudApp' : props.src
    }.png`
  }

  let prtclk = ''
  if (props.src) {
    if (props.onClick != null || props.pr != null) {
      prtclk = 'prtclk'
    }
  }

  const clickDispatch = event => {
    let action = {
      type: event.currentTarget.dataset.action,
      payload: event.currentTarget.dataset.payload
    }

    if (action.payload === 'winRefresh') {
      dispatch(action)
    }

    if (action.payload === 'close') {
      window.history.pushState(null, null, '/')
      window.localStorage.removeItem('CURRENTROUTERPATH')
    }

    // 如果已经激活重新点开则最小化
    if (action.type) {
      // 如果不是点击应用，走原有的逻辑
      if (!('active' in props)) {
        dispatch(action)
      } else if (props.active) {
        // 如果点击的是应用，并且已经处于active中，最小化
        dispatch({
          type: event.currentTarget.dataset.action,
          payload: 'mnmz'
        })
      } else if ('active' in props && !props.active) {
        // 点击的是应用，但是没有激活，那就最大化
        dispatch({
          type: event.currentTarget.dataset.action,
          payload: 'full'
        })
      }
    }
  }

  if (props.fafa != null) {
    return (
      <div
        className={`uicon prtclk ${props.className || ''}`}
        onClick={props.onClick || (props.click && clickDispatch) || null}
        data-action={props.click}
        data-payload={props.payload}
        data-menu={props.menu}
        data-routerpath={props.routerPath}
        data-sessionid={props.sessionId}
        data-url={props.url}
        data-icon={props.iconSrc}
        style={{
          width: props.width || '100%',
          height: props.height || '100%',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center'
        }}>
        <FontAwesomeIcon
          data-flip={props.flip != null}
          data-invert={props.invert != null ? 'true' : 'false'}
          data-rounded={props.rounded != null ? 'true' : 'false'}
          style={{
            width: props.width,
            height: props.height || props.width,
            color: props.color || null,
            margin: props.margin || null
          }}
          spin={props.spin != null ? true : false}
          icon={
            props.reg == null ? FaIcons[props.fafa] : FaRegIcons[props.fafa]
          }
        />
      </div>
    )
  } else if (props.icon != null) {
    let CustomIcon = AllIcons[props.icon]
    return (
      <div
        className={`uicon prtclk ${props.className || ''}`}
        onClick={props.onClick || (props.click && clickDispatch) || null}
        data-action={props.click}
        data-payload={props.payload}
        data-menu={props.menu}
        data-routerpath={props.routerPath}
        data-sessionid={props.sessionId}
        data-url={props.url}
        data-icon={props.iconSrc}
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center'
        }}>
        <CustomIcon
          data-flip={props.flip != null}
          data-invert={props.invert != null ? 'true' : 'false'}
          data-rounded={props.rounded != null ? 'true' : 'false'}
          style={{
            width: props.width,
            height: props.height || props.width,
            fill: props.color || null,
            margin: props.margin || null
          }}
        />
      </div>
    )
  } else {
    return (
      <div
        className={`uicon ${props.className || ''} ${prtclk}`}
        data-open={props.open}
        data-action={props.click}
        data-active={props.active}
        data-payload={props.payload}
        data-routerpath={props.routerPath}
        data-sessionid={props.sessionId}
        data-url={props.url}
        data-icon={props.iconSrc}
        onClick={props.onClick || (props.pr && clickDispatch) || null}
        data-menu={props.menu}
        data-pr={props.pr}>
        {props.className == 'tsIcon' ? (
          <div
            onClick={props.click != null ? clickDispatch : null}
            style={{
              width: '100%',
              height: '100%',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center'
            }}
            data-action={props.click}
            data-payload={props.payload}
            data-click={props.click != null}
            data-flip={props.flip != null}
            data-invert={props.invert != null ? 'true' : 'false'}
            data-rounded={props.rounded != null ? 'true' : 'false'}>
            <img
              width={props.width}
              height={props.height}
              data-action={props.click}
              data-payload={props.payload}
              data-click={props.click != null}
              data-flip={props.flip != null}
              data-invert={props.invert != null ? 'true' : 'false'}
              data-rounded={props.rounded != null ? 'true' : 'false'}
              src={iconSrc || src}
              onError={e => {
                e.currentTarget.src = defaultAppIcon
              }}
              style={{
                margin: props.margin || null
              }}
              alt=''
            />
          </div>
        ) : (
          <img
            width={props.width}
            height={props.height}
            onClick={props.click != null ? clickDispatch : null}
            data-action={props.click}
            data-payload={props.payload}
            data-click={props.click != null}
            data-flip={props.flip != null}
            data-invert={props.invert != null ? 'true' : 'false'}
            data-rounded={props.rounded != null ? 'true' : 'false'}
            src={iconSrc || src}
            onError={e => {
              e.currentTarget.src = defaultAppIcon
            }}
            style={{
              margin: props.margin || null
            }}
            alt=''
          />
        )}
      </div>
    )
  }
}

export const Image = props => {
  const dispatch = useDispatch()
  let src = `/img/${(props.dir ? props.dir + '/' : '') + props.src}.png`
  if (props.ext != null) {
    src = props.src
  }

  const errorHandler = e => {
    if (props.err) {
      e.currentTarget.src = props.err
    }
  }

  const clickDispatch = event => {
    let action = {
      type: event.currentTarget.dataset.action,
      payload: event.currentTarget.dataset.payload
    }

    if (action.type) {
      dispatch(action)
    }
  }

  return (
    <div
      className={`imageCont prtclk ${props.className || ''}`}
      id={props.id}
      style={{
        backgroundImage: props.back && `url(${src})`
      }}
      data-back={props.back != null}
      onClick={props.onClick || (props.click && clickDispatch)}
      data-action={props.click}
      data-payload={props.payload}
      data-var={props.var}>
      {!props.back ? (
        props.lazy ? (
          <LazyLoadImage
            width={props.w}
            height={props.h}
            data-free={props.free != null}
            data-var={props.var}
            loading={props.lazy ? 'lazy' : null}
            src={src}
            alt=''
            onError={errorHandler}
          />
        ) : (
          <img
            width={props.w}
            height={props.h}
            data-free={props.free != null}
            data-var={props.var}
            loading={props.lazy ? 'lazy' : null}
            src={src}
            alt=''
            onError={errorHandler}
          />
        )
      ) : null}
    </div>
  )
}

export const SnapScreen = props => {
  const dispatch = useDispatch()
  const [delay, setDelay] = useState(false)
  const lays = useSelector(state => state.globals.lays)

  const vr = 'var(--radii)'

  const clickDispatch = event => {
    event.preventDefault()
    event.stopPropagation()
    let action = {
      type: event.currentTarget.dataset.action,
      payload: event.currentTarget.dataset.payload,
      dim: JSON.parse(event.currentTarget.dataset.dim)
    }

    if (action.dim && action.type) {
      dispatch(action)
      props.closeSnap()
    }
  }

  useEffect(() => {
    if (delay && props.snap) {
      setTimeout(() => {
        setDelay(false)
      }, 500)
    } else if (props.snap) {
      setDelay(true)
    }
  })

  return props.snap || delay ? (
    <div className='snapcont mdShad' data-dark={props.invert != null}>
      {lays.map((x, i) => {
        return (
          <div key={i} className='snapLay'>
            {x.map((y, j) => (
              <div
                key={j}
                className='snapper'
                style={{
                  borderTopLeftRadius: (y.br % 2 == 0) * 4,
                  borderTopRightRadius: (y.br % 3 == 0) * 4,
                  borderBottomRightRadius: (y.br % 5 == 0) * 4,
                  borderBottomLeftRadius: (y.br % 7 == 0) * 4
                }}
                onClick={clickDispatch}
                data-dim={JSON.stringify(y.dim)}
                data-action={props.app}
                data-payload='resize'></div>
            ))}
          </div>
        )
      })}
    </div>
  ) : null
}

export const ToolBar = props => {
  const dispatch = useDispatch()
  const [snap, setSnap] = useState(false)
  const [fullScreen, setFullScreen] = useState(false)
  const openSnap = () => {
    setSnap(true)
  }

  const closeSnap = () => {
    setSnap(false)
  }

  const toolClick = e => {
    e.stopPropagation()
    e.preventDefault()
    dispatch({
      type: props.app,
      payload: 'front'
    })
  }

  let posP = [0, 0],
    dimP = [0, 0],
    posM = [0, 0],
    wnapp = {},
    op = 0,
    vec = [0, 0]

  const toolDrag = e => {
    e.preventDefault()
    e.stopPropagation()
    e = e || window.event
    posM = [e.clientY, e.clientX]
    op = e.currentTarget.dataset.op

    if (op == 0) {
      wnapp =
        e.currentTarget.parentElement &&
        e.currentTarget.parentElement.parentElement
    } else {
      vec = e.currentTarget.dataset.vec.split(',')
      wnapp =
        e.currentTarget.parentElement &&
        e.currentTarget.parentElement.parentElement &&
        e.currentTarget.parentElement.parentElement.parentElement
    }

    if (wnapp) {
      wnapp.classList.add('notrans')
      wnapp.classList.add('z9900')
      posP = [wnapp.offsetTop, wnapp.offsetLeft]
      dimP = [
        parseFloat(getComputedStyle(wnapp).height.replaceAll('px', '')),
        parseFloat(getComputedStyle(wnapp).width.replaceAll('px', ''))
      ]
    }

    document.onmouseup = closeDrag
    document.onmousemove = eleDrag
  }

  const setPos = (pos0, pos1) => {
    wnapp.style.top = pos0 + 'px'
    wnapp.style.left = pos1 + 'px'
  }

  const setDim = (dim0, dim1) => {
    wnapp.style.height = dim0 + 'px'
    wnapp.style.width = dim1 + 'px'
  }

  const eleDrag = e => {
    e.preventDefault()
    e.stopPropagation()
    e = e || window.event

    let pos0 = posP[0] + e.clientY - posM[0],
      pos1 = posP[1] + e.clientX - posM[1],
      dim0 = dimP[0] + vec[0] * (e.clientY - posM[0]),
      dim1 = dimP[1] + vec[1] * (e.clientX - posM[1])

    if (op == 0) setPos(pos0, pos1)
    else {
      dim0 = Math.max(dim0, 320)
      dim1 = Math.max(dim1, 320)
      pos0 = posP[0] + Math.min(vec[0], 0) * (dim0 - dimP[0])
      pos1 = posP[1] + Math.min(vec[1], 0) * (dim1 - dimP[1])
      setPos(pos0, pos1)
      setDim(dim0, dim1)
    }
  }

  const closeDrag = () => {
    document.onmouseup = null
    document.onmousemove = null

    wnapp.classList.remove('notrans')
    wnapp.classList.remove('z9900')

    let action = {
      type: props.app,
      payload: 'resize',
      dim: {
        width: getComputedStyle(wnapp).width,
        height: getComputedStyle(wnapp).height,
        top: getComputedStyle(wnapp).top,
        left: getComputedStyle(wnapp).left
      }
    }

    dispatch(action)
  }
  const onClickFillScreen = () => {
    if (!document.fullscreenElement) {
      // 如果未进入全屏模式，则切换到全屏模式
      setFullScreen(true)
      document.documentElement.requestFullscreen()
    } else {
      // 如果已经进入全屏模式，则退出全屏模式
      if (document.exitFullscreen) {
        setFullScreen(false)
        document.exitFullscreen()
      }
    }
  }
  return (
    <>
      <div
        className='win11Toolbar'
        data-float={props.float != null}
        data-noinvert={props.noinvert != null}
        style={{
          background: props.bg
        }}>
        <div
          className='topInfo flex flex-grow items-center'
          // style={{ pointerEvents: 'none' }}
          data-float={props.float != null}
          onClick={toolClick}
          onMouseDown={toolDrag}
          data-op='0'>
          <Icon src={props.icon} width={14} />
          <div
            className='appFullName text-xss'
            data-white={props.invert != null}>
            {props.name}
          </div>
          {props.child && props.child}
        </div>
        <div className='actbtns flex items-center'>
          {props.hasRefresh && (
            <Icon
              invert={props.invert}
              click={props.app}
              payload='winRefresh'
              pr
              src='refresh'
              ui
              width={12}
            />
          )}
          <Icon
            invert={props.invert}
            click={props.app}
            payload='mnmz'
            pr
            src='minimize'
            ui
            width={12}
          />
          <div
            className='snapbox h-full'
            data-hv={snap}
            onClick={onClickFillScreen}
            // onMouseOver={openSnap}
            // onMouseLeave={closeSnap}
          >
            <Icon
              invert={props.invert}
              // click={props.app}
              ui
              pr
              width={12}
              payload='mxmz'
              src={fullScreen ? 'fullscreen-exit' : 'fullscreen'}
            />
            <SnapScreen
              invert={props.invert}
              app={props.app}
              snap={snap}
              // closeSnap={closeSnap}
            />
            {/* {snap?<SnapScreen app={props.app} closeSnap={closeSnap}/>:null} */}
          </div>
          <Icon
            className='closeBtn'
            invert={props.invert}
            click={props.app}
            payload='close'
            pr
            src='close'
            ui
            width={14}
          />
        </div>
      </div>
      <div className='resizecont topone'>
        <div className='flex'>
          <div
            className='conrsz cursor-nw-resize'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='-1,-1'></div>
          <div
            className='edgrsz cursor-n-resize wdws'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='-1,0'></div>
        </div>
      </div>
      <div className='resizecont leftone'>
        <div className='h-full'>
          <div
            className='edgrsz cursor-w-resize hdws'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='0,-1'></div>
        </div>
      </div>
      <div className='resizecont rightone'>
        <div className='h-full'>
          <div
            className='edgrsz cursor-w-resize hdws'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='0,1'></div>
        </div>
      </div>
      <div className='resizecont bottomone'>
        <div className='flex'>
          <div
            className='conrsz cursor-ne-resize'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='1,-1'></div>
          <div
            className='edgrsz cursor-n-resize wdws'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='1,0'></div>
          <div
            className='conrsz cursor-nw-resize'
            data-op='1'
            onMouseDown={toolDrag}
            data-vec='1,1'></div>
        </div>
      </div>
    </>
  )
}

export const LazyComponent = ({ show, children }) => {
  const [loaded, setLoad] = useState(false)

  useEffect(() => {
    if (show && !loaded) setLoad(true)
  }, [show])

  return show || loaded ? <>{children}</> : null
}
