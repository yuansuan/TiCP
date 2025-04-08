export default function MessageReducer(state = {}, action) {
  switch (action.type) {
      case "updateMessage":
          return Object.assign({}, state, action.message);
      default:
          return state;
  }
}