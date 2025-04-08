/* Copyright (C) 2016-present, Yuansuan.cn */

import { history } from '@/utils'

export const historyService = {
  'history.push': (_, ...args) => {
    Reflect.apply(history.push, history, args)
  },
  'history.replace': (_, ...args) => {
    Reflect.apply(history.replace, history, args)
  }
}
