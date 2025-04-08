import * as React from 'react'

import { StyledSection } from './style'
import { InlineDiff } from './Diff'

const appPropMap = {
  name: '模版名称',
  state: '模版状态',
  description: '模版描述',
  type: '应用类型',
  ttl: '作业存放时间',
  image: '镜像名',
  cloud_out_app_id: '关联云应用ID',
  enable_residual: '残差图',
  enable_snapshot: '云图',
  residual_log_parser: '残差图日志解析器',
  scheduler_param: '调度器参数'
}

interface IProps {
  application: any
  iconData: any
}

export default class InfoDiff extends React.Component<IProps> {
  render() {
    const { application, iconData } = this.props
    const showkeys = Object.keys(application).filter(key => appPropMap[key])

    if (showkeys.length === 0 && !iconData) {
      return null
    }

    return (
      <StyledSection>
        <div>
          <span className='tag'>模版信息</span>
        </div>
        <div>
          {showkeys.map(key => {
            const prop = application[key]
            if (typeof prop.old === 'boolean') {
              prop.old = prop.old ? '启用' : '禁用'
            }
            if (typeof prop.new === 'boolean') {
              prop.new = prop.new ? '启用' : '禁用'
            }

            return (
              <InlineDiff
                key={key}
                name={appPropMap[key]}
                Old={prop.old}
                New={prop.new}
              />
            )
          })}
          {iconData && (
            <InlineDiff
              name='模版图标'
              Old={<img style={{ width: 64, height: 64 }} src={iconData.old} />}
              New={<img style={{ width: 64, height: 64 }} src={iconData.new} />}
            />
          )}
        </div>
      </StyledSection>
    )
  }
}
