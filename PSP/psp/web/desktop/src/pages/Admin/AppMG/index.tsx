import { action, observable } from 'mobx'
import { Provider, observer } from 'mobx-react'
import * as React from 'react'

import List from './List'
import {
  AppList,
  FavoriteList,
  RemoteAppList,
  RemoteFavoriteList
} from '@/domain/Applications'

@observer
export default class AppMG extends React.Component {
  @observable appList = new AppList()
  @observable favoriteList = new FavoriteList()
  @observable remoteAppList = new RemoteAppList()
  @observable remoteFavoriteList = new RemoteFavoriteList()

  @observable public app = null
  @action
  public updateApp = app => (this.app = app)

  public render() {
    return (
      <Provider
        appList={this.appList}
        favoriteList={this.favoriteList}
        remoteAppList={this.remoteAppList}
        remoteFavoriteList={this.remoteFavoriteList}>
        <List />
      </Provider>
    )
  }
}
