/* Copyright (C) 2016-present, Yuansuan.cn */
import { allApps } from '../utils'

const defState = {}
for (let i = 0; i < allApps.length; i++) {
  defState[allApps[i].icon] = allApps[i]
  defState[allApps[i].icon].size = 'full'
  defState[allApps[i].icon].hide = true
  defState[allApps[i].icon].max = null
  defState[allApps[i].icon].z = 0
}

defState.hz = 2

function defStateFun(apps) {
  const defState = {}
  for (let i = 0; i < apps.length; i++) {
    defState[apps[i].icon] = apps[i]
    defState[apps[i].icon].size = 'full'
    defState[apps[i].icon].hide = true
    defState[apps[i].icon].max = null
    defState[apps[i].icon].z = 0
  }

  defState.hz = 2
  return defState
}
const appReducer = (state = defState, action) => {
  action.payload === 'generateApp' &&
    action.type === 'apps' &&
    (state = defStateFun(action?.data))

  let tmpState = {
    ...state
  }
  if (action.type == 'EDGELINK') {
    let obj = {
      ...tmpState['edge']
    }
    if (action.payload && action.payload.startsWith('http')) {
      obj.url = action.payload
    } else if (action.payload && action.payload.length != 0) {
      obj.url = 'https://www.bing.com/search?q=' + action.payload
    } else {
      obj.url = null
    }

    obj.size = 'full'
    obj.hide = false
    obj.max = true
    tmpState.hz += 1
    obj.z = tmpState.hz
    tmpState['edge'] = obj
    return tmpState
  } else if (action.type == 'NEWJODETAIL') {
    let obj = {
      ...tmpState['jobDetail']
    }

    if (action.payload == 'mnmz') {
      obj.max = false
      obj.hide = false
      if (obj.z == tmpState.hz) {
        tmpState.hz -= 1
      }
      obj.z = -1
    }
    if (action.payload == 'close') {
      obj.hide = true
      obj.max = null
      obj.z = -1
      tmpState.hz -= 1
    } else if (action.payload == 'full') {
      obj.size = 'full'
      obj.hide = false
      obj.max = true
      tmpState.hz += 2
      obj.z = tmpState.hz
    }
    tmpState['jobDetail'] = obj
    return tmpState
  } else if (action.type == 'NEWJOBSETDETAIL') {
    let obj = {
      ...tmpState['jobSetDetail']
    }

    if (action.payload == 'mnmz') {
      obj.max = false
      obj.hide = false
      if (obj.z == tmpState.hz) {
        tmpState.hz -= 1
      }
      obj.z = -1
    }
    if (action.payload == 'close') {
      obj.hide = true
      obj.max = null
      obj.z = -1
      tmpState.hz -= 1
    } else if (action.payload == 'full') {
      obj.size = 'full'
      obj.hide = false
      obj.max = true
      tmpState.hz += 2
      obj.z = tmpState.hz
    }
    tmpState['jobSetDetail'] = obj
    return tmpState
  } else if (action.type == 'SHOWDSK') {
    let keys = Object.keys(tmpState)

    for (let i = 0; i < keys.length; i++) {
      let obj = tmpState[keys[i]]
      if (obj.hide == false) {
        obj.max = false
        if (obj.z == tmpState.hz) {
          tmpState.hz -= 1
        }
        obj.z = -1
        tmpState[keys[i]] = obj
      }
    }

    return tmpState
  } else if (action.type == 'EXTERNAL') {
    window.open(action.payload, '_blank')
  } else if (action.type == 'OPENTERM') {
    let obj = {
      ...tmpState['terminal']
    }
    obj.dir = action.payload

    obj.size = 'full'
    obj.hide = false
    obj.max = true
    tmpState.hz += 1
    obj.z = tmpState.hz
    tmpState['terminal'] = obj
    return tmpState
  } else if (action.type == 'ADDAPP') {
    tmpState[action.payload.icon] = action.payload
    tmpState[action.payload.icon].size = 'full'
    tmpState[action.payload.icon].hide = true
    tmpState[action.payload.icon].max = null
    tmpState[action.payload.icon].z = 0

    return tmpState
  } else if (action.type == 'DELAPP') {
    delete tmpState[action.payload]
    return tmpState
  } else {
    let keys = Object.keys(state)
    for (let i = 0; i < keys.length; i++) {
      let obj = state[keys[i]]
      if (obj.action == action.type) {
        tmpState = {
          ...state
        }

        if (action.payload == 'full') {
          obj.size = 'full'
          obj.hide = false
          obj.max = true
          tmpState.hz += 1
          obj.z = tmpState.hz
        } else if (action.payload == 'close') {
          if (obj.z == 0) {
            return tmpState
          }
          obj.hide = true
          obj.max = null
          obj.z = -1
          tmpState.hz -= 1
        } else if (action.payload == 'mxmz') {
          obj.size = ['mini', 'full'][obj.size != 'full' ? 1 : 0]
          obj.hide = false
          obj.max = true
          tmpState.hz += 1
          obj.z = tmpState.hz
        } else if (action.payload == 'togg') {
          if (obj.z != tmpState.hz) {
            obj.hide = false
            if (!obj.max) {
              tmpState.hz += 1
              obj.z = tmpState.hz
              obj.max = true
            } else {
              obj.z = -1
              obj.max = false
            }
          } else {
            obj.max = !obj.max
            obj.hide = false
            if (obj.max) {
              tmpState.hz += 1
              obj.z = tmpState.hz
            } else {
              obj.z = -1
              tmpState.hz -= 1
            }
          }
        } else if (action.payload == 'winRefresh') {
          obj.winRefresh = true
        } else if (action.payload == 'mnmz') {
          obj.max = false
          obj.hide = false
          if (obj.z == tmpState.hz) {
            tmpState.hz -= 1
          }
          obj.z = -1
        } else if (action.payload == 'resize') {
          obj.size = 'cstm'
          obj.hide = false
          obj.max = true
          if (obj.z != tmpState.hz) tmpState.hz += 1
          obj.z = tmpState.hz
          obj.dim = action.dim
        } else if (action.payload == 'front') {
          obj.hide = false
          obj.max = true
          if (obj.z != tmpState.hz) {
            tmpState.hz += 1
            obj.z = tmpState.hz
          }
        }

        tmpState[keys[i]] = obj
        return tmpState
      }
    }
  }

  return state
}

export default appReducer
