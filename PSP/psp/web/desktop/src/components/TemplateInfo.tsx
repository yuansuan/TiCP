import * as React from 'react'
import styled from 'styled-components'
import { observable, action, toJS } from 'mobx'
import { observer } from 'mobx-react'
import {
  Select,
  message,
  Upload,
  Input as AntdInput,
  Descriptions,
  Switch
} from 'antd'
import TodoList from '@/components/TodoList'
import TodoList2 from '@/components/TodoList2'
import { ValidInput } from '@/components'
import { createMobxStream } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { AppList, RemoteAppList } from '@/domain/Applications'
import { sysConfig } from '@/domain'
import { DeployMode } from '@/constant'

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;

  .module {
    display: flex;
    align-items: center;
    margin: 10px 0;

    .required {
      color: red;
    }

    .name {
      width: 120px;
      font-family: 'PingFangSC-Regular';
      font-size: 12px;
      color: rgba(38, 38, 38, 0.65);
      text-align: right;
      margin-right: 10px;
    }

    .widget {
      max-width: 60%;

      input,
      .ant-select {
        width: 200px;
        font-family: 'PingFangSC-Regular';
        font-size: 12px;
        color: rgba(38, 38, 38, 0.65);
      }

      .desc {
        word-break: break-all;
      }
      .binPathWrap {
        display: flex;
        flex-direction: column;
        max-width: 650px;
        font-size: 12px;
        .binPathItem {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 5px 0;
          word-break: break-all;
        }
      }
    }

    .iconImage {
      width: 92px;
      height: 46px;
    }

    .ant-upload.ant-upload-select-picture-card {
      position: relative;
      width: 120px;
      height: 64px;
      border: 1px dashed #d9d9d9;
      margin: 5px 0;

      .loading {
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right: 0;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
`

function getBase64(img, callback) {
  const reader = new FileReader()
  reader.addEventListener('loadend', () => callback(reader.result))
  reader.readAsDataURL(img)
}

interface IProps {
  appList?: AppList
  remoteAppList?: RemoteAppList
  isRemote?: boolean
  model: {
    name: string
    updateName?: (name: string) => void
    icon: string
    updateIcon?: (icon: string) => void
    description: string
    updateDesc?: (desc: string) => void
    newVersion: string
    updateVersion?: (version: string) => void
    newType?: string
    updateType?: (newType: string) => void
    updateBaseName?: (baseName: string) => void
    updateBaseVersion?: (baseVersion: string) => void
    updateBaseState?: (baseState: string) => void
    binPath?: Array<{ key: string; value: string }>
    updateBinPath?: (binPath: any) => void
    schedulerParam?: Array<{ key: string; value: string }>
    updateSchedulerParam?: (params: any) => void
    queues?: Array<{ queue_name: string; cpu_number: number; select: boolean }>
    updateQueue?: (queue: any) => void
    licenses?: Array<{
      id: string
      name: string
      select: boolean
      licence_valid: boolean
    }>
    updateLicense?: (license: any) => void
    image?: string
    updateImage?: (name: string) => void
    cloudOutAppId?: string
    updateCloudOutAppId?: (id: string) => void
    enableResidual?: boolean
    updateEnableResidual?: (bool: boolean) => void
    enableSnapshot?: boolean
    updateEnableSnapshot?: (bool: boolean) => void
    residualLogParser?: string
    updateResidualLogParser?: (name: string) => void
  }
}

@observer
export default class TemplateInfo extends React.Component<IProps> {
  @observable uploading = false
  @observable iconUrl = ''
  @action
  updateUploading = uploading => (this.uploading = uploading)
  @action
  updateIconUrl = iconUrl => (this.iconUrl = iconUrl)

  async componentDidMount() {
    createMobxStream(() => this.props.model.icon)
      .pipe(untilDestroyed(this))
      .subscribe(this.updateIconUrl)
  }

  private uploadIcon = ({ file }: any) => {
    if (file.type.match('image.*')) {
      getBase64(file, this.props.model.updateIcon)
    }
  }

  private beforeUpload = file => {
    const isImage = ['image/jpg', 'image/jpeg', 'image/png'].includes(file.type)
    if (!isImage) {
      message.error('只能上传图片文件')
    }
    const sizeLimit = file.size / 1024 / 1024 < 1
    if (!sizeLimit) {
      message.error('只能上传小于 1MB 的图片')
    }
    return isImage && sizeLimit
  }

  // select base template
  private onSelectBase = async name => {
    const { updateBaseName, updateBaseState, updateBaseVersion } =
      this.props.model

    if (name) {
      const app = this.props.appList.list.get(name)
      if (app) {
        updateBaseName(name)
        updateBaseState(app.state)
        updateBaseVersion(app.version)
      }
    }
  }

  private validate = (value, field) => {
    if (value.length > 64 && field === 'image') {
      message.error('镜像名最长限制为 64 个字符')
      return false
    }

    return true
  }
  private onBlurHandle = e => {
    const newName = e.target.value
    if (this.props.appList?.list.get(newName)) {
      message.error(`模版${newName}已经出现，请重新输入`)

      return Promise.reject(`模版${newName}已经出现，请重新输入`)
    } else {
      return true
    }
  }

  render() {
    const {
      description,
      updateDesc,
      newVersion = '',
      updateVersion,
      newType = '',
      updateType,
      updateBaseName,
      updateIcon,
      image,
      updateImage,
      binPath = [],
      updateBinPath,
      schedulerParam = [],
      updateSchedulerParam,
      queues = [],
      updateQueue,
      licenses = [],
      updateLicense,
      cloudOutAppId,
      updateCloudOutAppId,
      enableResidual,
      updateEnableResidual,
      enableSnapshot,
      updateEnableSnapshot,
      residualLogParser,
      updateResidualLogParser
    } = this.props.model

    const defaultCloudOutApp =
      this.props.remoteAppList?.list?.size > 0
        ? [...this.props.remoteAppList]?.find(
            item => item.outAppId === cloudOutAppId
          )
        : {}
    const currentQueues =
      queues.length > 0 ? queues : this.props.appList?.queueList
    const newQueue = currentQueues
      ? JSON.parse(JSON.stringify(currentQueues))
      : []
    const onSelectQueues = async values => {
      const { updateQueue } = this.props.model
      newQueue.forEach(item => {
        if (values?.includes(item.queue_name)) {
          item.select = true
        } else {
          item.select = false
        }
      })

      if (newQueue) {
        updateQueue(newQueue)
      }
    }
    const currentLicense =
      licenses.length > 0 ? licenses : this.props.appList?.licenseList
    const newLicense = currentLicense ? toJS(currentLicense) : []
    const onSelectLicense = async values => {
      const { updateLicense } = this.props.model
      newLicense.forEach(item => {
        if (values === item.name) {
          if (!item.licence_valid) {
            message.warn(`当前许可证服务器${item.name}已过期，请确认使用！`)
          }
          item.select = true
        } else {
          item.select = false
        }
      })

      if (newLicense) {
        updateLicense(newLicense)
      }
    }
    const RenderBinPath = ({ data = [] }) => {
      return (
        <div className='binPathWrap'>
          {data.map((item, index) => {
            return (
              <Descriptions className='binPathItem' key={index} size='small'>
                <Descriptions.Item className='binPathKey'>
                  {item?.key}
                </Descriptions.Item>
                <Descriptions.Item className='binPathValue'>
                  {item?.value}
                </Descriptions.Item>
              </Descriptions>
            )
          })}
        </div>
      )
    }
    return (
      <Wrapper>
        <div className='module'>
          <span className='name'>
            <span className='required'>*</span>类型：
          </span>
          <div className='widget'>
            {updateType ? (
              <ValidInput
                value={newType}
                onChange={e => updateType(e.target.value)}
              />
            ) : (
              newType
            )}
          </div>
        </div>
        <div className='module'>
          <span className='name'>
            <span className='required'>*</span>版本：
          </span>
          <div className='widget'>
            {updateVersion ? (
              <ValidInput
                value={newVersion}
                // autoFocus
                // onFocus={e => e.target.select()}

                onChange={e => updateVersion(e.target.value)}
              />
            ) : (
              newVersion
            )}
          </div>
        </div>

        {!this.props.isRemote && (
          <>
            <div className='module'>
              <span className='name'>
                <span className='required'>
                  {binPath?.length > 0 ? '' : '*'}
                </span>
                镜像名：
              </span>
              <div className='widget'>
                {updateImage ? (
                  <ValidInput
                    value={image}
                    onChange={e => updateImage(e.target.value)}
                    validator={value => this.validate(value, 'image')}
                  />
                ) : (
                  image
                )}
              </div>
            </div>
            <div className='module'>
              <span className='name'>
                <span className='required'>{image ? '' : '*'}</span>
                可执行文件路径：
              </span>
              <div className='widget'>
                {updateBinPath ? (
                  <TodoList
                    onChange={values => {
                      updateBinPath(values)
                    }}
                    defaultValues={binPath || []}
                    zoneList={this.props.appList?.zoneList}
                  />
                ) : (
                  binPath?.length > 0 && <RenderBinPath data={binPath} />
                )}
              </div>
            </div>
            <div className='module'>
              <span className='name'>
                <span className='required'></span>
                调度器参数：
              </span>
              <div className='widget'>
                {updateSchedulerParam ? (
                  <TodoList2
                    onChange={values => {
                      updateSchedulerParam(values)
                    }}
                    defaultValues={schedulerParam || []}
                  />
                ) : (
                  schedulerParam?.length > 0 && (
                    <RenderBinPath data={schedulerParam} />
                  )
                )}
              </div>
            </div>
            <div className='module'>
              <span className='name'>
                <span className='required'></span>许可证类型
              </span>
              <div className='widget'>
                {updateLicense ? (
                  <Select
                    showArrow
                    onChange={onSelectLicense}
                    allowClear={true}
                    defaultValue={licenses
                      .filter(item => item.select)
                      .map(item => item.name)}
                    placeholder='全部'>
                    {newLicense?.map(item => {
                      return (
                        <Select.Option
                          title={item.name}
                          key={item.name}
                          value={item.name}>
                          {item.name}
                        </Select.Option>
                      )
                    })}
                  </Select>
                ) : (
                  licenses
                    .filter(item => item.select)
                    .map(item => item.name)
                    .join(';')
                )}
              </div>
            </div>

            <div className='module'>
              <span className='name'>队列：</span>
              <div className='widget'>
                {updateQueue ? (
                  <Select
                    showArrow
                    mode='multiple'
                    onChange={onSelectQueues}
                    allowClear={true}
                    defaultValue={queues
                      .filter(item => item.select)
                      .map(item => item.queue_name)}
                    placeholder='全部'>
                    {newQueue?.map(item => {
                      return (
                        <Select.Option
                          title={item.queue_name}
                          key={item.queue_name}
                          value={item.queue_name}>
                          {item.queue_name}
                        </Select.Option>
                      )
                    })}
                  </Select>
                ) : (
                  queues
                    .filter(item => item.select)
                    .map(item => item.queue_name)
                    .join(';')
                )}
              </div>
            </div>
          </>
        )}
        <div className='module'>
          <span className='name'>残差图：</span>
          <div className='widget'>
            {updateEnableResidual && !this.props.isRemote ? (
              <Switch
                size='small'
                checked={enableResidual}
                checkedChildren='开启'
                unCheckedChildren='关闭'
                onClick={checked => {
                  updateEnableResidual(checked)
                }}
              />
            ) : enableResidual ? (
              '开启'
            ) : (
              '关闭'
            )}
          </div>
        </div>
        <div className='module'>
          <span className='name'>云图：</span>
          <div className='widget'>
            {updateEnableSnapshot && !this.props.isRemote ? (
              <Switch
                size='small'
                checked={enableSnapshot}
                checkedChildren='开启'
                unCheckedChildren='关闭'
                onClick={checked => {
                  updateEnableSnapshot(checked)
                }}
              />
            ) : enableSnapshot ? (
              '开启'
            ) : (
              '关闭'
            )}
          </div>
        </div>
        <div className='module'>
          <span className='name'>残差图日志解析器：</span>
          <div className='widget'>
            {updateResidualLogParser && !this.props.isRemote ? (
              <Select
                value={residualLogParser}
                onChange={name => {
                  updateResidualLogParser(name)
                }}
                allowClear={true}
                placeholder='全部'>
                {['starccm', 'fluent'].map(name => {
                  return (
                    <Select.Option title={name} key={name} value={name}>
                      {name}
                    </Select.Option>
                  )
                })}
              </Select>
            ) : (
              residualLogParser
            )}
          </div>
        </div>
        {updateBaseName && (
          <div className='module'>
            <span className='name'>基于模版：</span>
            <div className='widget'>
              <Select
                onChange={this.onSelectBase}
                allowClear={true}
                placeholder='全部'>
                {[...this.props.appList].map(item => {
                  const name = item.version
                    ? item.name + `(${item.version})`
                    : item.name
                  return (
                    <Select.Option
                      title={item.name + `(${item.version})`}
                      key={item.name}
                      value={item.name}>
                      {name}
                    </Select.Option>
                  )
                })}
              </Select>
            </div>
          </div>
        )}
        <div className='module'>
          <span className='name'>模版图标：</span>
          <div className='widget'>
            {updateIcon ? (
              <Upload
                listType='picture-card'
                showUploadList={false}
                beforeUpload={this.beforeUpload}
                customRequest={this.uploadIcon}>
                <div>
                  {this.iconUrl ? (
                    <img className='iconImage' src={this.iconUrl} />
                  ) : null}
                </div>
              </Upload>
            ) : (
              <img
                className='iconImage'
                src={
                  this.iconUrl
                    ? this.iconUrl
                    : require('@/assets/images/defaultApp.svg')
                }
              />
            )}
          </div>
        </div>

        <div className='module'>
          <span className='name' style={{ alignSelf: 'start' }}>
            模版描述：
          </span>
          <div className='widget'>
            {updateDesc ? (
              <AntdInput.TextArea
                value={description}
                style={{ resize: 'none', width: 420, height: 180 }}
                onChange={e => updateDesc(e.target.value)}
              />
            ) : (
              <span className='desc'>{description}</span>
            )}
          </div>
        </div>
      </Wrapper>
    )
  }
}
