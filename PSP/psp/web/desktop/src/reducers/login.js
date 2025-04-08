/* Copyright (C) 2016-present, Yuansuan.cn */

const defState = {
  needLogin: false
}

const loginReducer = (state = defState, action) => {
  switch (action.type) {
    case 'GOLOGIN':
      return {
        ...state,
        needLogin: true
      }
    case 'GOLOGOUT':
      return {
        ...state,
        needLogin: false
      }

    default:
      return state
  }
}

export default loginReducer
