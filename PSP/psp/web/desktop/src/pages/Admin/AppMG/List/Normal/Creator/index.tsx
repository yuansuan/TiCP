import * as React from 'react'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { message } from 'antd'

import { ModalFooter } from '@/components'
import { AppList, RemoteAppList } from '@/domain/Applications'
import { TemplateInfo } from '@/components'
import { Wrapper } from './style'

interface IProps {
  appList?: AppList
  remoteAppList?: RemoteAppList
  onCancel: () => void
  onOk: () => void
}

@observer
export default class TemplateCreator extends React.Component<IProps> {
  @observable name = ''
  @observable icon = ''
  @observable baseVersion = ''
  @observable newVersion = ''
  @observable description = ''
  @observable baseName = ''
  @observable baseState = ''
  @observable loading = false
  @observable newType = ''
  @observable binPath = []
  @observable schedulerParam = []
  @observable queues = []
  @observable licenses = []
  @observable cloudOutAppId = ''
  @observable image = ''
  @observable residualLogParser = ''
  @observable enableResidual = false
  @observable enableSnapshot = false

  @action
  updateName = name => (this.name = name)
  @action
  updateIcon = icon => (this.icon = icon)
  @action
  updateDesc = desc => (this.description = desc)
  @action
  updateBaseName = baseName => (this.baseName = baseName)
  @action
  updateBaseVersion = baseVersion => (this.baseVersion = baseVersion)
  @action
  updateBaseState = baseState => (this.baseState = baseState)
  @action
  updateLoading = loading => (this.loading = loading)
  @action
  updateVersion = version => (this.newVersion = version)
  @action
  updateType = type => (this.newType = type)
  @action
  updateImage = image => (this.image = image)
  @action
  updateCloudOutAppId = id => (this.cloudOutAppId = id)
  @action
  updateBinPath = binPath => (this.binPath = binPath)
  @action
  updateSchedulerParam = params => (this.schedulerParam = params)
  @action
  updateQueue = queue => (this.queues = queue)
  @action
  updateLicense = license => (this.licenses = license)
  @action
  updateResidualLogParser = name => (this.residualLogParser = name)
  @action
  updateEnableResidual = bool => (this.enableResidual = bool)
  @action
  updateEnableSnapshot = bool => (this.enableSnapshot = bool)

  async componentDidMount() {
    await this.props.appList.fetchZoneAreaList()
    await this.props.appList.fetchQueueList()
    await this.props.appList.fetchLicenseList()
  }
  private onCancel = () => {
    const { onCancel } = this.props

    onCancel()
  }

  private onOk = () => {
    let tip = ''
    if (!this.newType) {
      tip = '模版类型不能为空'
    } else if (!this.newVersion) {
      tip = '版本不能为空'
    } else if (!this.binPath.length && !this.image) {
      tip = '可执行文件路径或者镜像名必须填写一项！'
    }

    if (tip) {
      message.error(tip)
    } else {
      this.updateLoading(true)
      this.props.appList
        .add({
          newType: this.newType,
          baseName: this.baseName,
          queues: this.queues,
          licenses: this.licenses,
          baseVersion: this.baseVersion,
          icon: this.icon,
          bin_path: this.binPath,
          scheduler_param: this.schedulerParam,
          image: this.image,
          cloud_out_app_id: this.cloudOutAppId,
          newVersion: this.newVersion,
          description: this.description,
          enable_residual: this.enableResidual,
          enable_snapshot: this.enableSnapshot,
          residual_log_parser: this.residualLogParser
        })
        .then(() => {
          message.success('新建模版成功')
          this.props.onOk()
        })
        .finally(() => this.updateLoading(false))
    }
  }

  render() {
    const { onCancel, onOk } = this
    const { appList, remoteAppList } = this.props

    return (
      <Wrapper>
        <div className='form'>
          <TemplateInfo
            remoteAppList={remoteAppList}
            appList={appList}
            model={{
              name: this.name,
              icon: this.icon,
              description: this.description,
              newVersion: this.newVersion,
              newType: this.newType,
              binPath: this.binPath,
              schedulerParam: this.schedulerParam,
              queues: this.queues,
              licenses: this.licenses,
              image: this.image,
              cloudOutAppId: this.cloudOutAppId,
              enableResidual: this.enableResidual,
              enableSnapshot: this.enableSnapshot,
              residualLogParser: this.residualLogParser,
              updateResidualLogParser: this.updateResidualLogParser,
              updateEnableResidual: this.updateEnableResidual,
              updateEnableSnapshot: this.updateEnableSnapshot,
              updateBinPath: this.updateBinPath,
              updateSchedulerParam: this.updateSchedulerParam,
              updateQueue: this.updateQueue,
              updateLicense: this.updateLicense,
              updateImage: this.updateImage,
              updateType: this.updateType,
              updateVersion: this.updateVersion,
              updateName: this.updateName,
              updateIcon: this.updateIcon,
              updateDesc: this.updateDesc,
              updateBaseName: this.updateBaseName,
              updateBaseVersion: this.updateBaseVersion,
              updateBaseState: this.updateBaseState,
              updateCloudOutAppId: this.updateCloudOutAppId
            }}
          />
        </div>
        <ModalFooter className='footer' onCancel={onCancel} onOk={onOk} />
      </Wrapper>
    )
  }
}
