import * as React from 'react'
import { observer } from 'mobx-react'
import { AppList, RemoteAppList } from '@/domain/Applications'
import { inject } from '@/pages/context'
import { TemplateInfo } from '@/components'

interface IProps {
  app?: any
  isRemote: boolean
}

@inject(({ app }) => ({ app }))
@observer
export default class BaseInfo extends React.Component<IProps> {
  appList = new AppList()
  remoteAppList = new RemoteAppList()
  async componentDidMount() {
    await this.appList.fetchZoneAreaList()
    this.props.app.queues.length === 0 && (await this.appList.fetchQueueList())
    this.props.app.licenses.length === 0 && (await this.appList.fetchLicenseList())
    await this.remoteAppList.fetchTemplates()
  }
  render() {
    const { app } = this.props
    return (
      <div
        style={{
          padding: '20px',
          height: '100%',
          overflowY: 'auto',
          backgroundColor: '#f0f5fd'
        }}>
        <TemplateInfo
          appList={this.appList}
          remoteAppList={this.remoteAppList}
          isRemote={this.props.isRemote}
          model={{
            name: app.name,
            newType: app.type,
            icon: app.iconData,
            newVersion: app.version,
            description: app.description,
            image: app.image,
            queues: app.queues,
            licenses: app.licenses,
            binPath: app.bin_path,
            schedulerParam: app.scheduler_param,
            cloudOutAppId: app.cloud_out_app_id,
            residualLogParser: app.residual_log_parser,
            enableResidual: app.enable_residual,
            enableSnapshot: app.enable_snapshot,
            updateBinPath: binPath => (app.bin_path = binPath),
            updateSchedulerParam: params => (app.scheduler_param = params),
            updateImage: (name: string) => (app.image = name),
            updateIcon: icon => (app.iconData = icon),
            updateDesc: desc => (app.description = desc),
            updateQueue: queue => (app.queues = queue),
            updateLicense: license => (app.licenses = license),
            updateCloudOutAppId: id => (app.cloud_out_app_id = id),
            updateEnableResidual:  bool=> (app.enable_residual = bool),
            updateEnableSnapshot:  bool=> (app.enable_snapshot = bool),
            updateResidualLogParser: parse => (app.residual_log_parser = parse)
          }}
        />
      </div>
    )
  }
}
