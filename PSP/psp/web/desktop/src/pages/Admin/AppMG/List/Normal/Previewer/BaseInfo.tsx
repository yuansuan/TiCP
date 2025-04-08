import * as React from 'react'
import { observer } from 'mobx-react'

import { inject } from '@/pages/context'
import { TemplateInfo } from '@/components'

interface IProps {
  app?: any
  isRemote: boolean
}

@inject(({ app }) => ({ app }))
@observer
export default class BaseInfo extends React.Component<IProps> {
  render() {
    const { app, isRemote } = this.props

    return (
      <div
        style={{ padding: '20px', height: '100%', backgroundColor: '#f0f5fd' }}>
        <TemplateInfo
          isRemote={isRemote}
          model={{
            name: app.name,
            newType: app.type,
            icon: app.iconData,
            image: app.image,
            queues: app.queues,
            licenses: app.licenses,
            binPath: app.bin_path,
            schedulerParam: app.scheduler_param,
            cloudOutAppId: app.cloud_out_app_id,
            enableResidual: app.enable_residual,
            enableSnapshot: app.enable_snapshot,
            residualLogParser: app.residual_log_parser,
            newVersion: app.version,
            description: app.description
          }}
        />
      </div>
    )
  }
}
