export default function MeReducer(state = {}, action) {
  switch (action.type) {
      case "updateAuthInfo":
          return Object.assign({}, state, { authInfo: action.authInfo });
      case "updateUserInfo":
          return Object.assign({}, state, { userInfo: action.userInfo });
      case "updatAppMap":
          let set = action.appMap;
          let map = {};
          set.forEach((app) => {
              map[app.appId] = app;
          });
          return Object.assign({}, state, { appMap: map, appList: set });
      default:
          return state;
  }
}