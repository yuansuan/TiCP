/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import { PermList } from '@/domain/UserMG'
import { Icon } from '@/components'
import { Section, RadiusItem } from '..'
import { PermPreviewWrapper, AuthorityWrapper } from './style'
import { StatsBall } from '@/components'
import { sysConfig } from '@/domain'
import { INSTALL_TYPE } from '@/utils/const'
import { DeployMode } from '@/constant'

interface IProps {
  perms: PermList
  sysIcon?: string
  sfIcon?: string
}

@observer
export default class PermPreview extends React.Component<IProps> {
  @computed
  get systemPerms() {
    const sys = this.props.perms.systemPerms
    return sys.filter(p => p.has) || []
  }

  @computed
  get subAppPerms() {
    return this.props.perms.subAppPerms
  }

  @computed
  get remoteAppPerms() {
    return this.props.perms.remoteAppPerms
  }
  @computed
  get cloudAppPerms() {
    return this.props.perms.cloudAppPerms
  }

  get isAIO() {
    return sysConfig.installType === INSTALL_TYPE.aio
  }

  render() {
    const { sysIcon, sfIcon } = this.props

    return (
      <PermPreviewWrapper>
        <Section
          icon={sysIcon ? <Icon type={sysIcon} /> : null}
          className='Systems'
          title={sysIcon ? '系统权限' : '系统权限'}>
          <RadiusItem itemList={this.systemPerms.map(p => p.name)} />
        </Section>

        <Section
          icon={sfIcon ? <Icon type={sfIcon} /> : null}
          title={sfIcon ? '软件使用权限' : '软件使用权限'}>
          <AuthorityWrapper>
            <div className='appAuthority'>
              <StatsBall color='black'>本地应用</StatsBall>
            </div>
            {this.subAppPerms.map(i => (
              <span
                title={i.name}
                key={i.id}
                className={i.has ? 'enable bubble' : 'disable bubble'}>
                {i.name}
              </span>
            ))}
          </AuthorityWrapper>

          {sysConfig.globalConfig.enable_visual && (
            <AuthorityWrapper>
              <div className='appAuthority'>
                <StatsBall color='black'>3D可视化镜像</StatsBall>
              </div>
              {this.remoteAppPerms.map(i => (
                <span
                  title={i.name}
                  key={i.id}
                  className={i.has ? 'enable bubble' : 'disable bubble'}>
                  {i.name}
                </span>
              ))}
            </AuthorityWrapper>
          )}
        </Section>
      </PermPreviewWrapper>
    )
  }
}
