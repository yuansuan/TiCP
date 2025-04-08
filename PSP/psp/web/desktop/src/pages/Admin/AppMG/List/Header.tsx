import * as React from 'react'
import styled from 'styled-components'
import { Button, message } from 'antd'
import { Subject } from 'rxjs'
import { observer, inject } from 'mobx-react'
import { observable, action } from 'mobx'
import {
  RollbackOutlined,
  SendOutlined,
  PlusSquareOutlined} from '@ant-design/icons'
import { Modal } from '@/components'
import { fromStream } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { RemoteAppList, AppList } from '@/domain/Applications'
import NormalCreator from './Normal/Creator'
import { Search } from '@/components'

const HeaderWrapper = styled.div`
  width: 100%;
  padding: 0 10px 10px 10px;

  .toolbar {
    display: flex;
    justify-content: space-between;

    .button-group {
      button {
        margin: 0 10px;
      }
    }

    .search-input {
      width: 150px;
      margin-left: auto;
      margin-right: 10px;
    }
  }
`

interface IProps {
  appList?: AppList
  remoteAppList?: RemoteAppList
  keyword: any
  selectedRowKeys: string[]
}

@inject(({ appList, remoteAppList }) => ({ appList, remoteAppList }))
@observer
export default class Header extends React.Component<IProps> {
  keyword$ = new Subject<string>()

  appList = this.props.appList
  Creator = NormalCreator

  @observable publishing = false
  @observable unpublishing = false
  @observable syncClouding = false
  @observable settingTemplate = false

  @action
  updatePublishing = publishing => (this.publishing = publishing)
  @action
  updateUnPublishing = unpublishing => (this.unpublishing = unpublishing)
  @action
  updateSyncClouding = syncClouding => (this.syncClouding = syncClouding)
  @action
  updateSettingTemplate = settingTemplate =>
    (this.settingTemplate = settingTemplate)

  componentDidMount() {
    fromStream(this.keyword$.pipe(untilDestroyed(this)), this.props.keyword)
  }

  private openEditor = async () => {
    const Creator = this.Creator as React.ElementType
    await this.props.remoteAppList.fetchTemplates()

    Modal.show({
      title: '新建模版',
      width: 800,
      bodyStyle: { height: 550, padding: 0 },
      footer: null,
      content: ({ onCancel, onOk }) => (
        <Creator
          appList={this.appList}
          remoteAppList={this.props.remoteAppList}
          onCancel={onCancel}
          onOk={onOk}
        />
      )
    })
  }

  // publish template
  private publish = async e => {
    e.stopPropagation()

    if (this.publishing) {
      return
    }

    const { selectedRowKeys } = this.props

    this.updatePublishing(true)
    try {
      await this.appList.publish(selectedRowKeys)
      message.success('发布成功')
      this.updatePublishing(false)
    } catch (err) {
      this.updatePublishing(false)
    }
  }

  // unpublish template
  private unpublish = async e => {
    e.stopPropagation()
    if (this.unpublishing) {
      return
    }

    const { selectedRowKeys } = this.props
    this.updateUnPublishing(true)
    try {
      await this.appList.unpublish(selectedRowKeys)
      message.success('取消发布成功')
      this.updateUnPublishing(false)
    } catch (err) {
      this.updateUnPublishing(false)
    }
  }

  render() {
    const { selectedRowKeys } = this.props

    return (
      <HeaderWrapper>
        <div className='toolbar'>
          <div className='button-group'>
            {(
              <Button
                ghost
                type='primary'
                icon={<PlusSquareOutlined />}
                onClick={this.openEditor}>
                新建
              </Button>
            )}

            <Button
              ghost
              type='primary'
              icon={<SendOutlined />}
              disabled={selectedRowKeys.length === 0}
              loading={this.publishing}
              onClick={this.publish}>
              发布
            </Button>

            <Button
              ghost
              type='primary'
              icon={<RollbackOutlined />}
              disabled={selectedRowKeys.length === 0}
              loading={this.unpublishing}
              onClick={this.unpublish}>
              取消发布
            </Button>
          </div>

          <Search
            className='search-input'
            placeholder='请输入模版名称'
            onSearch={value => this.keyword$.next(value)}
          />
        </div>
      </HeaderWrapper>
    )
  }
}
