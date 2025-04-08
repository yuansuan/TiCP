import React from 'react'
import { observable, action } from 'mobx'
import { observer, inject } from 'mobx-react'
import { message, Tooltip } from 'antd'
import {
  FormOutlined,
  RollbackOutlined,
  SendOutlined,
  DeleteOutlined
} from '@ant-design/icons'

import {
  AppList,
  FavoriteList,
  RemoteAppList,
  RemoteFavoriteList
} from '@/domain/Applications'
import { history } from '@/utils'
import { Modal } from '@/components'
import { StyledOperators } from './style'

interface IProps {
  model: any
  isRemote?: boolean
  appList?: AppList
  favoriteList?: FavoriteList
  remoteAppList?: RemoteAppList
  remoteFavoriteList?: RemoteFavoriteList
}

@inject(({ appList, favoriteList, remoteAppList, remoteFavoriteList }) => ({
  appList,
  favoriteList,
  remoteAppList,
  remoteFavoriteList
}))
@observer
export default class Operators extends React.Component<IProps> {
  appList = this.props.isRemote ? this.props.remoteAppList : this.props.appList
  favoriteList = this.props.isRemote
    ? this.props.remoteFavoriteList
    : this.props.favoriteList

  @observable publishing = false
  @action
  updatePublishing = publishing => (this.publishing = publishing)

  render() {
    const {
      isRemote,
      model: { state }
    } = this.props

    return (
      <StyledOperators>
        {!isRemote &&
          (state === 'published' ? (
            <span className='item disabled'>
              <Tooltip title='编辑'>
                <FormOutlined />
              </Tooltip>
            </span>
          ) : (
            <span className='item' onClick={this.openEditor}>
              <Tooltip title='编辑'>
                <FormOutlined />
              </Tooltip>
            </span>
          ))}
        {isRemote &&
          (state === 'published' ? (
            <span className='item disabled'>
              <Tooltip title='编辑'>
                <FormOutlined />
              </Tooltip>
            </span>
          ) : (
            <span className='item' onClick={this.openEditor}>
              <Tooltip title='编辑'>
                <FormOutlined />
              </Tooltip>
            </span>
          ))}
        {state === 'published' ? (
          <span className='item' onClick={this.unpublish}>
            <Tooltip title='取消发布'>
              <RollbackOutlined />
            </Tooltip>
          </span>
        ) : (
          <span className='item' onClick={this.publish}>
            <Tooltip title='发布'>
              <SendOutlined />
            </Tooltip>
          </span>
        )}

        {!isRemote &&
          (state === 'published' ? (
            <span className='item disabled'>
              <Tooltip title='删除'>
                <DeleteOutlined />
              </Tooltip>
            </span>
          ) : (
            <span className='item' onClick={this.delete}>
              <Tooltip title='删除'>
                <DeleteOutlined />
              </Tooltip>
            </span>
          ))}
      </StyledOperators>
    )
  }

  private openEditor = e => {
    e.stopPropagation()

    const { model, isRemote } = this.props

    history.push({
      pathname: '/sys/template-edit',
      search: `?app=${model.name}&version=${model.version}&appId=${
        model.appId
      }&isRemote=${!!isRemote}`
    })
    window.localStorage.setItem(
      'CURRENTROUTERPATH',
      `/sys/template-edit?app=${model.name}&version=${model.version}&appId=${model.appId}&isRemote=${isRemote}`
    )
  }

  // publish template
  private publish = e => {
    e.stopPropagation()

    if (this.publishing) {
      return
    }

    const { model } = this.props

    this.updatePublishing(true)
    model
      .publish()
      .then(() => message.success('发布成功'))
      .finally(() => this.updatePublishing(false))
  }

  // unpublish template
  private unpublish = async e => {
    e.stopPropagation()

    if (this.publishing) {
      return
    }

    const { model } = this.props

    this.updatePublishing(true)
    model
      .unpublish()
      .then(() => message.success('取消发布成功'))
      .finally(() => this.updatePublishing(false))
  }

  // delete template
  private delete = e => {
    e.stopPropagation()

    const { model } = this.props

    Modal.showConfirm().then(() =>
      this.appList.delete(model.name).then(() => {
        message.success('模版删除成功')
      })
    )
  }
}
