export default function WebRTCReducer(state = { latency: 0 }, action) {
  switch (action.type) {
    case 'updateCopyOutText':
      return Object.assign({}, state, { copyOutText: action.text })
    case 'setLatency':
      return Object.assign({}, state, { latency: action.latency })
    case 'userInactive':
      return Object.assign({}, state, { inactive: action.inactive })
    case 'lastMessageSent':
      return Object.assign({}, state, {
        lastMessageSent: action.lastMessageSent,
      })
    default:
      return state
  }
}
