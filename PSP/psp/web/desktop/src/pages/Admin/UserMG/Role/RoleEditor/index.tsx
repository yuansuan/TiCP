import * as React from 'react'
import { Input, message, Spin, Switch, Tabs, Checkbox } from 'antd'
import { action, computed, observable } from 'mobx'
import { observer } from 'mobx-react'
import { Role, RoleList } from '@/domain/UserMG'
import { Validator } from '@/utils'
import { Section, PendingFooter } from '../../components'
import { checkPerm } from '../../utils'
import { RoleEditorWrapper, RoleBasicInfoWrapper, SysWrapper } from './style'
import { sysConfig } from '@/domain'
import { INSTALL_TYPE } from '@/utils/const'
import { DeployMode } from '@/constant'

interface IProps {
  role: Role
  onCancel?: any
  onOk?: any
  canEdit?: any
  isAdd?: any
}

const { TabPane } = Tabs

@observer
export default class RoleEditor extends React.Component<IProps> {
  @observable errorMessage = ''
  @action
  public updateName = name => (this.props.role.name = name)
  @action
  public updateComment = comment => (this.props.role.comment = comment)

  @action
  public updatePermIds = (perms, permNames) => {
    this.props.role.permIds = perms
    this.props.role.permNames = permNames
  }

  @action
  public updateLoading = flag => (this.loading = flag)
  @computed
  get editMode() {
    return this.props.role.id > 0
  }

  get nameEditable() {
    // return this.props.role.type === RoleType.CUSTOM
    return true
  }

  @computed
  get permIds() {
    return this.props.role.permIds
  }

  @computed
  get permNames() {
    return this.props.role.permNames
  }

  @observable public loading = false

  public componentDidMount() {
    this.props.role.fetch()
  }
  private onBlurName = e => {
    const name = e.target.value
    if (!Validator.isValidInputName(name)) {
      this.errorMessage = '角色名称只能包含字母,汉字,数字和下划线'
    } else if (name.length > 64) {
      this.errorMessage = '角色名称长度不能大于 64 字符'
    } else {
      this.errorMessage = ''
    }
  }

  private onBlurComment = e => {
    const comment = e.target.value
    if (!Validator.isValidTextArea(comment)) {
      this.errorMessage = '角色描述只能包含字母,汉字,数字,下划线,空格和换行'
    } else if (comment.length > 255) {
      this.errorMessage = '角色描述长度不能大于 255 字符'
    } else {
      this.errorMessage = ''
    }
  }

  get isAIO() {
    return sysConfig.installType === INSTALL_TYPE.aio
  }

  @computed
  get tabMaps() {
    return [
      {
        ...({
              title: '本地应用',
              store: this.props.role.permList.subAppPerms
            })
      },
      {
        ...(sysConfig.globalConfig?.enable_visual
          ? {
              title: '3D可视化镜像',
              store: this.props.role.permList.remoteAppPerms
            }
          : null)
      }
    ]
  }

  public render() {
    const { loading } = this
    const { role } = this.props

    return (
      <RoleEditorWrapper>
        {loading && (
          <div className='loading'>
            <Spin />
          </div>
        )}

        <div className='body'>
          <RoleBasicInfoWrapper>
            <div className='module'>
              <span className='name'>
                <span className='warn'>*</span>角色名称：
              </span>
              {this.nameEditable ? (
                <div className='widget'>
                  <Input
                    autoFocus
                    maxLength={64}
                    onFocus={e => e.target.select()}
                    value={role.name}
                    onBlur={this.onBlurName}
                    onChange={this.onChangeName}
                  />
                </div>
              ) : (
                <span>{role.name}</span>
              )}
            </div>
            <div className='module module-bottom'>
              <span className='name'>角色描述：</span>
              {this.nameEditable ? (
                <div className='widget'>
                  <Input.TextArea
                    value={role.comment}
                    maxLength={255}
                    onChange={this.onChangeComment}
                    onBlur={this.onBlurComment}
                    placeholder='0/100'
                  />
                </div>
              ) : (
                <span>{role.comment}</span>
              )}
            </div>
          </RoleBasicInfoWrapper>

          {
            <SysWrapper>
              <span className='name'>系统权限：</span>
              <div className='body'>
                {role.permList.systemPerms
                  // .filter(i =>
                  //   sysConfig.enableThreeMembers
                  //     ? ![
                  //         15, //'系统管理员日志',
                  //         16, //'安全管理员日志',
                  //         14 // '普通用户日志',
                  //       ].includes(i.id)
                  //     : true
                  // )
                  .map((item, index) => (
                    <div key={index} className='sysperm'>
                      <Checkbox
                        checked={this.permIds.includes(item.id)}
                        onChange={e =>
                          this.onCheckPerm(e.target.checked, item.id)
                        }
                      />
                      <span title={item.name} className='permcheck'>
                        {item.name}
                      </span>
                    </div>
                  ))}
              </div>
            </SysWrapper>
          }

          {
            <Section title='软件使用权限：'>
              <Tabs defaultActiveKey={this.tabMaps[0].title}>
                {this.tabMaps
                  .filter(t => t.title)
                  .map(t => (
                    <TabPane tab={t.title} key={t.title}>
                      <ul className='Softwares'>
                        <li className='perm title' key='title'>
                          <label className='sf'>软件名称</label>
                          <label className='op'>操作</label>
                        </li>
                        {t.store.map(item => (
                          <li key={item.id} className='perm'>
                            <label title={item.name} className='sf'>
                              {item.name}
                            </label>
                            <Switch
                              className='switch'
                              checked={this.permIds.includes(item.id)}
                              onChange={checked =>
                                this.onCheckPerm(checked, item.id)
                              }
                            />
                          </li>
                        ))}
                      </ul>
                    </TabPane>
                  ))}
              </Tabs>
            </Section>
          }
        </div>
        <PendingFooter
          onCancel={this.props.onCancel}
          onOk={this.onOk}
          processing={this.loading}
        />
      </RoleEditorWrapper>
    )
  }

  private onOk = () => {
    const { onOk, role } = this.props

    if (!role.name) {
      message.error('请填写角色名称')
      return
    }

    if (this.errorMessage) {
      // 校验不通过
      message.error(this.errorMessage)
      return
    }

    // edit role
    if (this.editMode) {
      // update role
      this.updateLoading(true)
      role
        .update()
        .then(() => {
          message.success('角色更新成功')
          onOk && onOk()
        })
        .catch(e => {
          if (e.fake) {
            if (e.success) {
              // 关闭对话框
              this.props.onCancel()
            }
          }
        })
        .finally(() => {
          this.updateLoading(false)
          checkPerm()
        })
    } else {
      // add role
      this.updateLoading(true)
      RoleList.add({
        name: role.name,
        comment: role.comment,
        perms: this.permIds
        // permNames: this.permNames
      })
        .then(() => {
          message.success('角色新增成功')
          onOk && onOk()
        })
        .catch(e => {
          if (e.fake) {
            if (e.success) {
              // 关闭对话框
              this.props.onCancel()
            }
          }
        })
        .finally(() => this.updateLoading(false))
    }
  }

  private onCheckPerm = (checked: boolean, id: number) => {
    let newPermIds = [...this.props.role.permIds]

    if (checked) {
      newPermIds.push(id)
    } else {
      newPermIds = newPermIds.filter(pid => pid !== id)
    }

    const system = this.props.role.permList.systemPerms
      .filter(p => newPermIds.includes(p.id))
      .map(p => p.name)
    const app = this.props.role.permList.subAppPerms
      .filter(p => newPermIds.includes(p.id))
      .map(p => p.name)

    this.updatePermIds(newPermIds, [...system, ...app])
  }

  private onChangeName = e => this.updateName(e.target.value)
  private onChangeComment = e => this.updateComment(e.target.value)
}
