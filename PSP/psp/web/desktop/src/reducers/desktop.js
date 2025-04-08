/* Copyright (C) 2016-present, Yuansuan.cn */
import { desktopApps } from '../utils'

const defState = {
  apps: [],
  hide: false,
  size: 1,
  sort: 'none',
  abOpen: false
}

const deskReducer = (state = defState, action) => {
  action.payload === 'generateApp' &&
    action.type === 'desktop' &&
    (state.apps = action?.data)
  if (action.payload === 'closeSession') {
    // state
    state.apps.forEach(item => {
      if (item?.id === action.data.sessionId) {
        item.className = 'CloudAppWrap_close'
        item.action = item.icon
        item.menu = 'cloudMenuClose'
        delete item.payload
      }
    })
  } else if (action.payload === 'openSession') {
    state.apps.forEach(item => {
      if (item?.id === action.data.sessionId) {
        item.className = 'CloudAppWrap_open'
        item.action = item.icon
        item.menu = 'cloudMenuOpen'
        delete item.payload
      }
    })
  }
  
  switch (action.type) {
    case 'DESKREM':
      let arr = state.apps.filter(x => x.name != action.payload)

      localStorage.setItem('desktop', JSON.stringify(arr.map(x => x.name)))
      return { ...state, apps: arr }
    case 'DESKADD':
      arr = [...state.apps]
      arr.push(action.payload)

      localStorage.setItem('desktop', JSON.stringify(arr.map(x => x.name)))
      return { ...state, apps: arr }
    case 'DESKHIDE':
      return {
        ...state,
        hide: true
      }
    case 'DESKSHOW':
      return {
        ...state,
        hide: false
      }
    case 'DESKTOGG':
      return {
        ...state,
        hide: !state.hide
      }
    case 'DESKSIZE':
      return {
        ...state,
        size: action.payload
      }
    case 'DESKSORT':
      return {
        ...state,
        sort: action.payload || 'none'
      }
    case 'DESKABOUT':
      return {
        ...state,
        abOpen: action.payload
      }
    default:
      return state
  }
}

export default deskReducer
