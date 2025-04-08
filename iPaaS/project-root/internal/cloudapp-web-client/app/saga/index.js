import { fork } from 'redux-saga/effects';
import appSaga from './AppSaga';
import WebRTCSaga from './WebRTCSaga'

export default function* rootSaga() {
    yield fork(appSaga);
    yield fork(WebRTCSaga);
}