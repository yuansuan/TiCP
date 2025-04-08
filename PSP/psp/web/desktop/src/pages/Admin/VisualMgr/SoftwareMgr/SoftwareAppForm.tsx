import React, { useState, useRef, useEffect } from 'react'
import { observer } from 'mobx-react'
import { Descriptions, message, Input, Button, Switch, Upload } from 'antd'
import { Label } from '@/components'
import { Validator } from '@/utils'
import styled from 'styled-components'

const TextArea = Input.TextArea

export const FormWrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    flex-direction: column;
  }

  .formItem {
    width: 400px;
  }

  .ant-descriptions-item {
  }

  .iconImage {
    width: 100px;
    height: 100px;
  }

  .footer {
    position: absolute;
    display: flex;
    bottom: 0px;
    right: 0;
    width: 100%;
    line-height: 64px;
    height: 64px;
    background: white;

    .footerMain {
      margin-left: auto;
      margin-right: 8px;

      button {
        margin: 0 8px;
      }
    }
  }
`

export const Tips = styled.span`
  font-family: PingFangSC-Regular;
  font-size: 12px;
  color: #999999;
  line-height: 22px;
`
function getBase64(img, callback) {
  const reader = new FileReader()
  reader.addEventListener('loadend', () => callback(reader.result))
  reader.readAsDataURL(img)
}

export const SoftwareAppForm = observer(props => {
  const inputRef = useRef(null)
  const [adding, setAdding] = useState(false)

  const [appAttr, setAppAttr] = useState(Object.assign({
    name: '',
    version: '',
    path: '',
    icon_data: '',
    os_type: '',
    gpu_support: false,
    app_param: '',
    app_param_paths: '',
    only_show_desktop: false,
  }, props.data || {}))

  const [validateErr, setValidateErr] = useState({
    id: false,
    name: false,
    version: false,
    path: false,
    icon_data: false,
    os_type: false,
    gpu_support: false,
    app_param: false,
    app_param_paths: false,

    // TODO
    vm_list: false,
    only_show_desktop: false,
    process_name: false,
  })

  const validateItems = ['name', 'version', 'os_type', 'path']

  const getErrMsg = (attr, value) => {
    if (attr === 'name') {
      if (value === '') return '软件名称不能为空'
      if (!Validator.isValidLicenseName(value)) {
        return '软件名称只能包含字母，数字，下划线，中括号或中划线'
      }
    }

    if (attr === 'version') {
      if (value === '') return '软件版本不能为空'
    }

    if (attr === 'os_type') {
      if (value === '') return '操作系统不能为空'
    }

    if (!appAttr.only_show_desktop) {
      if (attr === 'path') {
        if (value === '') return '软件路径不能为空'
      }
    }

    return ''
  }

  const validate = (attr, value) => {
    if (getErrMsg(attr, value)) {
      setValidateErr(prevState => ({
        ...prevState,
        [attr]: true,
      }))
    } else {
      setValidateErr(prevState => ({
        ...prevState,
        [attr]: false,
      }))
    }
  }

  const submit = async () => {

    // 全面触发校验，不用等字段自己 onblur
    validateItems.forEach(key => validate(key, appAttr[key]))

    // 检查校验状态
    if (validateItems.every(key => getErrMsg(key, appAttr[key]) === '')) {
      try {
        setAdding(true)
        await props.onOk(appAttr)
      } finally {
        setAdding(false)
      }
    }
  }

  useEffect(() => {
    inputRef.current.focus()
  }, [])

  const uploadIcon = ({ file }: any) => {
    if (file.type.match('image.*')) {
      getBase64(file, (base64Data) => {
        setAppAttr({
          ...appAttr,
          icon_data: base64Data
        })
      })
    }
  }

  const beforeUpload = file => {
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

  return (
    <FormWrapper>
      <Descriptions title='' column={1} style={{ margin: '0 0 50px 0' }}>
        <Descriptions.Item label={<Label required>软件名称</Label>}>
          <div className='item'>
            <Input
              ref={inputRef}
              style={{
                borderColor: validateErr.name ? '#f5222d' : 'inherit',
              }}
              className='formItem'
              maxLength={64}
              placeholder='请输入软件名称'
              value={appAttr.name}
              onBlur={e => {
                validate('name', e.target.value.trim())
              }}
              onChange={e => {
                setAppAttr({
                  ...appAttr,
                  name: e.target.value.trim(),
                })
              }}
            />
            <Tips>只能包含字母，数字，下划线，中括号或中划线</Tips>
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required>软件版本</Label>}>
          <div className='item'>
            <Input
              style={{
                borderColor: validateErr.version ? '#f5222d' : 'inherit',
              }}
              className='formItem'
              minLength={1}
              maxLength={64}
              placeholder='请输入软件版本'
              value={appAttr.version}
              onBlur={e => {
                validate('version', e.target.value.trim())
              }}
              onChange={e => {
                setAppAttr({
                  ...appAttr,
                  version: e.target.value.trim(),
                })
              }}
            />
            <Tips>输入软件版本号，例如 1.0 3.0 等</Tips>
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required>是否显示桌面</Label>}>
          <div className='item'>
            <Switch checkedChildren="是" unCheckedChildren="否" 
              checked={appAttr.only_show_desktop}
              onChange={(checked) => setAppAttr({
                ...appAttr,
                only_show_desktop: checked,
              })} 
              />
          </div>
        </Descriptions.Item>
        {
          !appAttr.only_show_desktop && (
            <Descriptions.Item label={<Label required>软件路径</Label>}>
              <div className='item'>
                <Input
                    style={{
                      borderColor: validateErr.path ? '#f5222d' : 'inherit',
                    }}
                    className='formItem'
                    placeholder='请输入软件路径'
                    value={appAttr.path}
                    onBlur={e => {
                      validate('path', e.target.value.trim())
                    }}
                    onChange={e => {
                      setAppAttr({
                        ...appAttr,
                        path: e.target.value.trim(),
                      })
                    }}
                  />
                  <Tips>输入软件路径，例如：/home/pbsadmin/app_script/hypermesh.sh</Tips>
              </div>
            </Descriptions.Item>
          )
        }
        <Descriptions.Item label={<Label>软件图标</Label>}>
          <div className='item'>
            <Upload
              listType='picture-card'
              showUploadList={false}
              beforeUpload={beforeUpload}
              customRequest={uploadIcon}>
              <div>
                {appAttr.icon_data ? (
                  <img className='iconImage' src={appAttr.icon_data} />
                ) : null}
              </div>
            </Upload>
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required>操作系统</Label>}>
          <div className='item'>
            <Input
                style={{
                  borderColor: validateErr.os_type ? '#f5222d' : 'inherit',
                }}
                className='formItem'
                placeholder='请输入操作系统'
                value={appAttr.os_type}
                onBlur={e => {
                  validate('os_type', e.target.value.trim())
                }}
                onChange={e => {
                  setAppAttr({
                    ...appAttr,
                    os_type: e.target.value.trim(),
                  })
                }}
              />
              <Tips>例如：centos7 或 win7</Tips>
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required>是否支持GPU</Label>}>
          <div className='item'>
            <Switch checkedChildren="是" unCheckedChildren="否" 
              checked={appAttr.gpu_support}
              onChange={(checked) => setAppAttr({
                ...appAttr,
                gpu_support: checked,
              })} 
              />
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label>参数</Label>}>
          <div className='item'>
            <TextArea
                style={{
                  borderColor: validateErr.app_param ? '#f5222d' : 'inherit',
                }}
                className='formItem'
                placeholder='请输入参数'
                value={appAttr.app_param}
                onBlur={e => {
                  validate('app_param', e.target.value.trim())
                }}
                onChange={e => {
                  setAppAttr({
                    ...appAttr,
                    app_param: e.target.value.trim(),
                  })
                }}
              />
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label>参数路径</Label>}>
          <div className='item'>
            <TextArea
                style={{
                  borderColor: validateErr.app_param_paths ? '#f5222d' : 'inherit',
                }}
                className='formItem'
                placeholder='请输入参数路径'
                value={appAttr.app_param_paths}
                onBlur={e => {
                  validate('app_param_paths', e.target.value.trim())
                }}
                onChange={e => {
                  setAppAttr({
                    ...appAttr,
                    app_param_paths: e.target.value.trim(),
                  })
                }}
              />
          </div>
        </Descriptions.Item>
      </Descriptions>
      <div className='footer'>
        <div className='footerMain'>
          <Button type='primary' loading={adding} onClick={submit}>
            确认
          </Button>
          <Button onClick={props.onCancel}>取消</Button>
        </div>
      </div>
    </FormWrapper>
  )
})
