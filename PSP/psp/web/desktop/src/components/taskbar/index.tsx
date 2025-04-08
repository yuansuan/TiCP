/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Icon } from '../../utils/general'
import Battery from '../shared/Battery'
// import './taskbar.scss'
import './taskbar.css'
import { Select } from 'antd'
import { Tooltip } from 'antd'
import { BrowserRouter } from 'react-router-dom'
import { HeaderToolbar } from '@/components/HeaderToolbar'
import { env } from '@/domain'

import history from '@/utils/history'

const Taskbar = () => {
  const tasks = useSelector(state => {
    return state.taskbar
  })
  const apps = useSelector(state => {
    let tmpApps = { ...state.apps }
    for (let i = 0; i < state.taskbar.apps.length; i++) {
      tmpApps[state.taskbar.apps[i].icon].task = true
    }
    return tmpApps
  })
  const dispatch = useDispatch()

  const showPrev = event => {
    let ele = event.target
    while (ele && ele.getAttribute('value') == null) {
      ele = ele.parentElement
    }

    let appPrev = ele.getAttribute('value')
    let xpos = window.scrollX + ele.getBoundingClientRect().left

    let offsetx = Math.round((xpos * 10000) / window.innerWidth) / 100

    dispatch({
      type: 'TASKPSHOW',
      payload: {
        app: appPrev,
        pos: offsetx
      }
    })
  }

  const hidePrev = () => {
    dispatch({ type: 'TASKPHIDE' })
  }

  const clickDispatch = (event, path) => {
    let action = {
      type: event.target.dataset.action,
      payload: event.target.dataset.payload
    }
    path &&
      window.localStorage.setItem('CURRENTROUTERPATH', JSON.stringify(path))

    if (action.type) {
      dispatch(action)
    }
  }

  const [time, setTime] = useState(new Date())

  useEffect(() => {
    const interval = setInterval(() => {
      setTime(new Date())
    }, 1000)
    return () => clearInterval(interval)
  }, [])

  const onChange = (value: string) => {
    console.log(`selected ${value}`)
  }

  const onSearch = (value: string) => {
    console.log('search:', value)
  }

  return (
    <div className='taskbar'>
      <div className='taskcont'>
        <div className='tasksCont' data-menu='task' data-side={tasks.align}>
          <BrowserRouter>
            {/* <UserInfo type='inside' /> */}
            <HeaderToolbar />
          </BrowserRouter>
          <div className='tsbar' onMouseOut={hidePrev}>
            {/* <Icon className="tsIcon" src="home" width={24} click="STARTOGG" />
            {tasks.search ? (
              <Icon
                click="STARTSRC"
                className="tsIcon searchIcon"
                icon="taskSearch"
              />
            ) : null}
            {tasks.widgets ? (
              <Icon
                className="tsIcon widget"
                src="widget"
                width={24}
                click="WIDGTOGG"
              />
            ) : null} */}
            {tasks.apps.map((task, i) => {
              let isHidden = apps[task.icon].hide
              let isActive = apps[task.icon].z == apps.hz
              return (
                <div
                  key={i}
                  // onClick={e => clickDispatch(e, task.routerPath)}
                  onMouseOver={(!isActive && !isHidden && showPrev) || null}
                  value={task.icon}>
                  <Tooltip title={task.name}>
                    <Icon
                      className='tsIcon'
                      width={24}
                      routerPath={task.routerPath}
                      open={isHidden ? null : true}
                      click={task.action}
                      payload={task.payload || 'full'}
                      active={isActive}
                      src={task.icon}
                    />
                  </Tooltip>
                </div>
              )
            })}
            {Object.keys(apps).map((key, i) => {
              if (key != 'hz') {
                var isActive = apps[key].z == apps.hz
              }
              return key != 'hz' &&
                key != 'undefined' &&
                !apps[key].task &&
                !apps[key].hide ? (
                <div
                  key={i}
                  onMouseOver={(!isActive && showPrev) || null}
                  // onClick={e => clickDispatch(e, apps[key].routerPath)}
                  value={apps[key].icon}>
                  <Icon
                    className='tsIcon'
                    width={24}
                    active={isActive}
                    click={apps[key].action}
                    routerPath={apps[key].routerPath}
                    payload={apps[key].payload || 'full'}
                    open='true'
                    src={apps[key].icon}
                  />
                </div>
              ) : null
            })}
          </div>
        </div>
        <div className='taskright'>
          {/* <div
            className="px-2 prtclk handcr hvlight flex"
            onClick={clickDispatch}
            data-action="BANDTOGG"
          >
            <Icon fafa="faChevronUp" width={10} />
          </div> */}

          {/* <div
            className='prtclk handcr my-1 px-1 hvlight flex rounded'
            onClick={clickDispatch}
            data-action='PANETOGG'>
            <Icon className="taskIcon" src="wifi" ui width={16} /> 
            <Icon
              className='taskIcon'
              // src={"audio" + tasks.audio}
              src={'sun'}
              ui
              width={16}
            />
           <Battery />
          </div>  */}

          {/* <div
            className='taskDate m-1 handcr prtclk rounded hvlight'
            onClick={clickDispatch}
            data-action='CALNTOGG'>
            <div>
              {time.toLocaleTimeString('en-US', {
                hour: 'numeric',
                minute: 'numeric'
              })}
            </div>
            <div>
              {time.toLocaleDateString('en-US', {
                year: '2-digit',
                month: '2-digit',
                day: 'numeric'
              })}
            </div>
          </div> */}
          {/* <Icon className='graybd my-4' ui width={6} click='SHOWDSK' pr /> */}
          {/* <Select
            showSearch
            placeholder="默认空间_创物虚拟空间"
            optionFilterProp="children"
            onChange={onChange}
            onSearch={onSearch}
            filterOption={(input, option) =>
              (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
            }
            options={[
              {
                value: 'jack',
                label: 'Jack',
              },
              {
                value: 'lucy',
                label: 'Lucy',
              },
              {
                value: 'tom',
                label: 'Tom',
              },
            ]}
          /> */}
        </div>
      </div>
    </div>
  )
}

export default Taskbar
