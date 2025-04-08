import { call, put, fork, takeLatest, throttle } from 'redux-saga/effects'

function* copyOutText(action) {
  yield put({ type: 'updateCopyOutText', text: action.text })
}
function* setLatency(action) {
  yield put({ type: 'setLatency', latency: action.latency })
}
function* userInactive(action) {
  yield put({ type: 'userInactive', inactive: action.inactive })
}
function* lastMessageSent(action) {
  yield put({
    type: 'lastMessageSent',
    lastMessageSent: action.lastMessageSent,
  })
}
function* WebRTCSaga() {
  yield throttle(2000, 'LAST_MESSAGE_SENT', lastMessageSent)
  yield takeLatest('USER_INACTIVE', userInactive)
  yield takeLatest('COPY_OUT_TEXT', copyOutText)
  yield takeLatest('SET_LATENCY', setLatency)
}

export default WebRTCSaga
