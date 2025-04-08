import { combineReducers } from 'redux';
import MeReducer from '@/reducers/MeReducer';
import WebRTCReducer from '@/reducers/WebRTCReducer';
import MessageReducer from '@/reducers/MessageReducer'

export default combineReducers({
    me:MeReducer,
    webrtc:WebRTCReducer,
    message:MessageReducer
})