/* Copyright (C) 2016-present, Yuansuan.cn */
let wps = localStorage.getItem('wps') || 0
let locked = localStorage.getItem('locked')

const walls = [
  'default/img0.jpg',
  'dark/img0.jpg',
]

const themes = ['default']

const defState = {
  themes: themes,
  wps: wps,
  src: walls[wps],
  locked: !(locked == 'false'),
  booted: false,
  act: '',
  dir: 0,
}

const wallReducer = (state = defState, action) => {
  switch (action.type) {
    case 'WALLUNLOCK':
      localStorage.setItem('locked', false)
      return {
        ...state,
        locked: false,
        dir: 0,
      }
    case 'WALLNEXT':
      let twps = (state.wps + 1) % walls.length
      localStorage.setItem('wps', twps)
      return {
        ...state,
        wps: twps,
        src: walls[twps],
      }
    case 'WALLALOCK':
      return {
        ...state,
        locked: true,
        dir: -1,
      }
    case 'WALLBOOTED':
      return {
        ...state,
        booted: true,
        dir: 0,
        act: '',
      }
    case 'WALLRESTART':
      return {
        ...state,
        booted: false,
        dir: -1,
        locked: true,
        act: 'restart',
      }
    case 'WALLSHUTDN':
      return {
        ...state,
        booted: false,
        dir: -1,
        locked: true,
        act: 'shutdn',
      }
    case 'WALLSET':
      let isIndex = !Number.isNaN(parseInt(action.payload)),
        wps = 0,
        src = ''

      if (isIndex) {
        wps = action.payload
        src = walls[action.payload]
      } else {
        src = action.payload
        wps = walls.indexOf(action.payload)
      }

      return {
        ...state,
        wps: wps,
        src: src,
      }
    default:
      return state
  }
}

export default wallReducer
