/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createReducer, createAction } from '@/utils/reducer'

function init(): boolean {
  return true
}

const reducer = createReducer(init(), handleAction => [
  handleAction(
    createAction('TOGGLE_MENU')<boolean>(),
    (_, { payload }) => payload
  ),
])

export default {
  reducer,
  init,
}
