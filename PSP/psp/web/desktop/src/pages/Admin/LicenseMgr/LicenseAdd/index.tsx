/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect } from 'react'
import { Page } from '@/components/Page'
import styled from 'styled-components'
import { useParams } from 'react-router'
import { Form, Input, Select, Button, message, Drawer } from 'antd'
import { useLocalStore, observer } from 'mobx-react-lite'
import { history, Http } from '@/utils'
import ConfigForm from './ConfigForm'
import { lmList } from '@/domain'
import { BackButton, BashEditor } from '@/components'

const Option = Select.Option
const TextArea = Input.TextArea

export const StyledLayout = styled.div`
  padding: 30px;
  overflow-y: hidden;

  .form {
    width: 600px;
    .back {
      font-size: 18px;
      font-family: PingFangSC-Medium;
      padding: 16px;
    }
    .back:hover {
      cursor: pointer;
      color: #3182ff;
    }
  }

  .footer {
    padding: 20px;
    .btn {
      margin: 10px;
    }
  }
`

const colProps = { labelCol: { span: 6 }, wrapperCol: { span: 14, offset: 1 } }

export default observer(function LicenseAdd() {
  const { id } = useParams<{ id: string }>()
  const [license, setLicense] = useState(null)
  const [show, setShow] = useState(false)
  const [form] = Form.useForm(null)

  const state = useLocalStore(() => ({
    appType: '',
    setAppType(str) {
      this.appType = str
    },
    appTypeList: [],
    setAppTypeList(list) {
      this.appTypeList = list
    },
    appTypeListLoading: false,
    appNameVersionList: [],
    setAppNameVersionList(list) {
      this.appNameVersionList = list
    },
    appNameVersionListLoading: false,
    appId: '',
    setAppId(str) {
      this.appId = str
    },
    provider: null, // 当前 provider
    setProvider(index, value) {
      this.provider = {
        index,
        value
      }
    },
    providerConfigs: []
  }))

  const open = index => {
    const license_info = form.getFieldValue('license_info')
    if (license_info[index]) {
      state.setProvider(index, license_info[index])
      setTimeout(() => {
        setShow(true)
      }, 0)
    } else {
      message.error('设置前，请输入提供者名称')
    }
  }

  const onClose = () => {
    state.setProvider(null, null)
    setShow(false)
  }

  useEffect(() => {
    ;(async function fetch() {
      if (id) {
        // 后端需要提供获取单个 license 的信息，与添加的返回一致
        // 通过联合查询返回 app_id 相关的 app_type, app_version
        const { data } = await Http.get(`/licenseManagers/${id}`)
        setLicense(data)

        form.setFieldsValue({
          appType: data?.app_type,
          os: data?.os,
          desc: data?.desc,
          compute_rule: data?.compute_rule
          // license_info: data?.listLicenseInfo.map(
          //   item => item.provider
          // ) || [undefined]
        })
        state.setAppType(data?.app_type)
        // state.providerConfigs =
        //   data?.item?.listLicenseInfo.map(item => {
        //     return {
        //       ...item,
        //       begin_time: item.begin_time.seconds
        //         ? moment(item.begin_time.seconds * 1000)
        //         : undefined,
        //       end_time: item.end_time.seconds
        //         ? moment(item.end_time.seconds * 1000)
        //         : undefined
        //     }
        //   }) || []
      }
    })()
  }, [])

  // useEffect(() => {
  //   ;(async function fetch() {
  //     // app type list
  //     state.appTypeListLoading = true
  //     const res = await Http.get('licenseManagers/typeList')
  //     state.setAppTypeList(res?.data.license_type_infos || [])
  //     // 非编辑模式
  //     if (!id) {
  //       state.setAppType(res?.data.license_type_infos[0])
  //       form.setFieldsValue({
  //         appType: res?.data.license_type_infos[0]
  //       })
  //     }
  //     state.appTypeListLoading = false
  //   })()
  // }, [])

  // useEffect(() => {
  //   ;(async function fetch() {
  //     // app name version list
  //     state.appNameVersionListLoading = true
  //     const res = await Http.get(
  //       `/licenseMgr/apps?type=${encodeURIComponent(state.appType)}`
  //     )
  //     state.setAppNameVersionList(res?.data || [])
  //     // 非编辑模式
  //     if (!id) {
  //       state.setAppId(res?.data[0]?.id)
  //       form.setFieldsValue({
  //         appId: res?.data[0]?.isUsed ? null : res?.data[0]?.id
  //       })
  //     }
  //     state.appNameVersionListLoading = false
  //   })()
  // }, [state.appType])

  function onFinish(values) {
    // let list = [...new Set(values.license_info.filter(Boolean))]

    // if (!state.providerConfigs[list.length - 1]) {
    //   message.error(`请设置提供者${list.length}的License`)
    //   return
    // }

    // let providerConfigsJSON = toJS(state.providerConfigs)

    let bodyData = {
      id,
      app_type: values.appType,
      os: values.os,
      desc: values.desc,
      compute_rule: values.compute_rule
      // license_info: list.map((item, index) => {
      //   let providerInfos = providerConfigsJSON[index]
      //   delete providerInfos.time
      //   return {
      //     ...providerInfos,
      //     provider: item,
      //     begin_time: providerInfos.begin_time?.valueOf(),
      //     end_time: providerInfos.end_time?.valueOf()
      //   }
      // })
    }

    if (id) {
      lmList
        .edit(id, bodyData)
        .then(res => {
          message.success('License编辑成功')
          back()
        })
        .catch(e => {
          if(e !==22014 ) message.error('License编辑失败')
        })
    } else {
      lmList
        .add(bodyData)
        .then(res => {
          message.success('License添加成功')
          back()
        })
        .catch(e => {
          if(e !==22014 ) message.error('License添加失败')
        })
    }
  }

  function editLine(value, index, type) {
    const keys = form.getFieldValue(type)

    keys.splice(index, 1, value)

    form.setFieldsValue({
      [type]: keys
    })
  }

  function addLine(type) {
    const keys = form.getFieldValue(type)

    if (keys.some(item => item === '' || item === undefined)) {
      form.validateFields(['license_info'])
      return
    }

    const list = keys.filter(Boolean)

    if (list.length > 1 && list.length !== new Set(list).size) {
      form.validateFields(['license_info'])
      // message.error('提供者名称不能重复')
      return
    }

    if (!state.providerConfigs[list.length - 1]) {
      form.validateFields(['license_info'])
      // message.error(`请设置提供者${list.length}的License`)
      return
    }

    form.setFieldsValue({
      [type]: [...new Set([...list, undefined])]
    })
  }

  function back() {
    history.push('/sys/license_mgr')
  }

  const onChangeContent = newcode => {
    if (!newcode) {
      return
    }
    form.setFieldsValue({ compute_rule: newcode })
  }
  const computeRule = form.getFieldValue('compute_rule')
  return (
    <Page header={null}>
      <StyledLayout>
        <Drawer
          width={860}
          title={`许可证${state?.provider?.index + 1}配置`}
          placement='right'
          onClose={onClose}
          visible={show}>
          <ConfigForm
            key={state?.provider?.index}
            onCancel={() => onClose()}
            onSubmit={values => {
              state.providerConfigs[state?.provider?.index] = values
              onClose()
            }}
            licenseConfig={
              state.providerConfigs[state?.provider?.index] || null
            }
          />
        </Drawer>
        <div className='form'>
          <BackButton
            title='返回许可证管理'
            style={{
              fontSize: 20
            }}
            onClick={() => window.history.back()}>
            {id ? '编辑License' : '添加License'}
          </BackButton>
          <Form
            name='license_edit_form'
            form={form}
            {...colProps}
            initialValues={{
              appType: state.appType,
              os: license?.os || 1,
              desc: license?.desc || '',
              compute_rule: license?.compute_rule || ''
            }}
            onFinish={onFinish}>
            <Form.Item
              label='许可证类型'
              name='appType'
              required
              rules={[{ required: true, message: '许可证类型不能为空' }]}>
              <Input style={{ width: 260 }} placeholder='请输入许可证类型' />
            </Form.Item>

            <Form.Item label='操作系统' name='os'>
              <Select>
                <Option value={1}>Linux</Option>
                <Option value={2}>Windows</Option>
              </Select>
            </Form.Item>
            <Form.Item
              label='计算规则'
              name='compute_rule'
              required
              rules={[{ required: true, message: '计算规则不能为空' }]}>
              <BashEditor
                code={computeRule}
                onChange={onChangeContent}
                readOnly={false}
                placeholder='请输入计算规则'
              />
            </Form.Item>
            <Form.Item label='描述' name='desc'>
              <TextArea placeholder='请输入描述' maxLength={300} />
            </Form.Item>
          </Form>
          <div className='footer'>
            <Button className='btn' onClick={() => back()}>
              关闭
            </Button>
            <Button
              className='btn'
              type='primary'
              onClick={() => form.submit()}>
             保存
            </Button>
          </div>
        </div>
      </StyledLayout>
    </Page>
  )
})
