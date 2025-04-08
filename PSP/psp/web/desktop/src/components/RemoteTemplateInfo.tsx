import * as React from 'react'
import styled from 'styled-components'
import { observable, action } from 'mobx'
import { observer } from 'mobx-react'
import { Input, Radio, message, Upload } from 'antd'

import { ValidInput } from '@/components'
import { createMobxStream, Http } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { filter } from 'rxjs/operators'

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;

  .table {
    flex-grow: 1;
  }

  .module {
    display: flex;
    align-items: center;
    margin: 8px 0;

    .required {
      color: red;
    }

    .name {
      min-width: 140px;
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
        width: 160px;
      }
      .app-name {
        font-weight: bold;
      }

      .desc {
        word-break: break-all;
      }
    }
    .visual-apps {
      height: 160px;
      overflow-y: auto;
      display: flex;
      align-items: center;
      .container {
        height: 100%;
      }
      .ant-radio-wrapper {
        height: 30px;
      }
      .header {
        padding-left: 24px;
        font-weight: 500;
        line-height: 40px;
      }
      .title {
        width: 160px;
      }
      .os {
        width: 100px;
      }
      .title,
      .os,
      .gpu {
        display: inline-block;
        font-size: 14px;
      }
    }

    .ant-upload.ant-upload-select-picture-card {
      position: relative;
      width: 120px;
      height: 64px;
      border: 1px dashed #d9d9d9;
      margin: 5px 0;

      img {
        width: 92px;
        height: 46px;
      }

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
  model: {
    name: string
    updateName?: (value: string) => void

    appId?: string
    updateAppId?: (value: string) => void

    osType?: string
    updateOsType?: (value: string) => void

    gpuEnabled?: boolean
    updateGpuEnabled?: (value: boolean) => void

    appPath: string
    updateAppPath?: (value: string) => void

    appParamsWithFile: string
    updateAppParamsWithFile?: (appParam: string) => void

    appParamsWithoutFile: string
    updateAppParamsWithoutFile?: (appParam: string) => void

    appFileTypes: string
    updateAppFileTypes?: (value: string) => void

    icon: string
    updateIcon?: (value: string) => void

    multiFiles: boolean
    updateMultiFiles?: (multiFile: boolean) => void

    desc: string
    updateDesc?: (desc: string) => void
  }
}

@observer
export default class RemoteTemplateInfo extends React.Component<IProps> {
  @observable uploading = false
  @observable iconUrl = ''
  @observable appSummary = ''
  @observable virtualApps = []
  @action
  updateVirtualApps = apps => (this.virtualApps = apps)
  @action
  updateUploading = uploading => (this.uploading = uploading)
  @action
  updateIconUrl = iconUrl => (this.iconUrl = iconUrl)

  columns: any[] = [
    {
      header: '',
      props: {
        width: 60,
      },
      cell: {
        props: {
          dataKey: 'id',
        },
        render: ({ dataKey, rowData }) => (
          <Radio
            checked={this.props.model.appId == rowData[dataKey]}
            onClick={() => this.props.model.updateAppId(rowData[dataKey])}
          />
        ),
      },
    },
    {
      props: {
        resizable: true,
        flexGrow: 1,
      },
      header: '名称',
      dataKey: 'name',
    },
    {
      props: {
        resizable: true,
      },
      header: '平台',
      dataKey: 'osType',
    },
    {
      props: {
        resizable: true,
        align: 'center',
      },
      header: 'GPU启用',
      cell: {
        props: {
          dataKey: 'gpuEnabled',
        },
        render: ({ dataKey, rowData }) => (
          <div style={{ textAlign: 'center' }}>
            {rowData[dataKey] ? '是' : '否'}
          </div>
        ),
      },
    },
  ]

  componentDidMount() {
    const { icon, appId, osType, gpuEnabled } = this.props.model
    Http.get('/visual/app/list', { baseURL: '' }).then(res => {
      this.updateVirtualApps(
        res.data.map(item => ({
          id: item.id,
          name: item.name,
          osType: item.os_type,
          gpuEnabled: item.gpu_support,
          appParam: item.app_param,
        }))
      )
      this.virtualApps.map(va => {
        if (appId && va.id == appId) {
          this.appSummary = `${va.name}, ${osType}, ${
            gpuEnabled ? '启用GPU' : '未启用GPU'
          }`
        }
      })
    })
    this.updateIconUrl(icon)
    createMobxStream(() => this.props.model.icon)
      .pipe(untilDestroyed(this))
      .subscribe(this.updateIconUrl)

    createMobxStream(() => this.props.model.appId)
      .pipe(
        untilDestroyed(this),
        filter(appId => !!appId)
      )
      .subscribe(appId => {
        const app = this.virtualApps.find(item => item.id == appId)
        if (app) {
          const { model } = this.props
          model.updateGpuEnabled(app.gpuEnabled)
          model.updateOsType(app.osType)
          model.updateAppParamsWithFile(app.appParam)
        }
      })
  }

  private validate = name => {
    if (/[^\u4e00-\u9fa5\w-\.]/.test(name)) {
      message.error('模版名称只能包含中文、数字、字母、下划线、中划线和点字符')
      return false
    } else if (name.length > 32) {
      message.error('模版名称最长限制为 32 个字符')
      return false
    }

    return true
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

  render() {
    const {
      name,
      updateName,
      icon,
      updateAppId,
      appParamsWithFile,
      updateAppParamsWithFile,
      appParamsWithoutFile,
      updateAppParamsWithoutFile,
      appFileTypes,
      updateAppFileTypes,
      desc,
      updateDesc,
    } = this.props.model

    return (
      <Wrapper>
        <div className='module'>
          <span className='name'>
            <span className='required'>*</span>模版名：
          </span>
          <div className='widget'>
            {updateName && (
              <ValidInput
                value={name}
                autoFocus
                onFocus={e => e.target.select()}
                validator={this.validate}
                onChange={e => updateName(e.target.value)}
              />
            )}
            {!updateName && <div className='app-name'>{name}</div>}
          </div>
        </div>
        <div className='module'>
          <span className='name'>
            <span className='required'>*</span>可视化应用：
          </span>
          {updateAppId && (
            <div className='visual-apps widget'>
              <div className='container'>
                <div className='header'>
                  <span className='title'>名称</span>
                  <span className='os'>平台</span>
                  <span className='gpu'>GPU启用</span>
                </div>
                <Radio.Group
                  value={this.props.model.appId}
                  onChange={e => this.props.model.updateAppId(e.target.value)}>
                  {this.virtualApps.map(app => {
                    return (
                      <Radio key={app.id} value={app.id.toString()}>
                        <span className='title'>{app.name}</span>
                        <span className='os'>{app.osType}</span>
                        <span className='gpu'>
                          {app.gpuEnabled ? '是' : '否'}
                        </span>
                      </Radio>
                    )
                  })}
                </Radio.Group>
              </div>
            </div>
          )}
          {!updateAppId && (
            <div className='widget'>
              <span>{this.appSummary}</span>
            </div>
          )}
        </div>

        <div className='module'>
          <span className='name'>模版图标：</span>
          <div className='widget'>
            {updateAppId && (
              <Upload
                listType='picture-card'
                showUploadList={false}
                beforeUpload={this.beforeUpload}
                customRequest={this.uploadIcon}>
                <div>{this.iconUrl ? <img src={this.iconUrl} /> : null}</div>
              </Upload>
            )}
            {!updateAppId && <img src={icon} width={46} height={46} />}
          </div>
        </div>
        {false && (
          <div className='module'>
            <span className='name'>参数模版：</span>
            <div className='widget'>
              {updateAppId && (
                <Input
                  value={appParamsWithoutFile}
                  onChange={e => updateAppParamsWithoutFile(e.target.value)}
                />
              )}
              {!updateAppId && <span>{appParamsWithoutFile}</span>}
            </div>
          </div>
        )}
        {false && (
          <div className='module'>
            <span className='name'>打开文件参数模版：</span>
            <div className='widget'>
              {updateAppId && (
                <Input
                  value={appParamsWithFile}
                  onChange={e => updateAppParamsWithFile(e.target.value)}
                />
              )}
              {!updateAppId && <span>{appParamsWithFile}</span>}
            </div>
          </div>
        )}
        <div className='module'>
          <span className='name'>关联文件类型：</span>
          <div className='widget'>
            {updateAppId && (
              <div>
                <Input
                  value={appFileTypes}
                  onChange={e => updateAppFileTypes(e.target.value)}
                />
                <span>(例如：*.hm, *3d*, *p?lot)</span>
              </div>
            )}
            {!updateAppId && <span>{appFileTypes}</span>}
          </div>
        </div>
        <div className='module'>
          <span className='name'>模版描述：</span>
          <div className='widget'>
            {updateAppId && (
              <Input.TextArea
                value={desc}
                style={{ resize: 'none', width: 420, height: 180 }}
                onChange={e => updateDesc(e.target.value)}
              />
            )}
            {!updateAppId && <span className='desc'>{desc}</span>}
          </div>
        </div>
      </Wrapper>
    )
  }
}
