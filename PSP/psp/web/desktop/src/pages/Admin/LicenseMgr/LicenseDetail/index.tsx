/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect, useRef } from 'react'
import { Page } from '@/components/Page'
import styled from 'styled-components'
import { useParams } from 'react-router'
import { Http } from '@/utils'
import { Descriptions, Tooltip, Drawer } from 'antd'
import { AppInfo } from '@/domain/LicenseMgr/AppInfo'
import { Modal, Button, Table } from '@/components'
import { QuantityDetail } from '../QuantityDetail'
import { BackButton } from '@/components'
import { ellipsisModal } from '..'
import Report from '@/pages/UniteReport/Report'

import ConfigForm from '../LicenseAdd/ConfigForm'
import moment from 'moment'
export const StyledLayout = styled.div`
  padding: 20px;
  overflow: auto;
  background: #fff;
  height: calc(100vh - 160px);
  .back {
    font-size: 18px;
    font-family: PingFangSC-Medium;
  }
  .ellipsis {
    display: inline-block;
    width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; 
  }
  .back:hover {
    cursor: pointer;
    color: #3182ff;
  }
`

const ListWrapper = styled.div`
  padding-top: 20px;
  padding-bottom: 30px;
  height: calc(100vh - 200px);
  > button {
    margin: 10px 0;
  }
`

type LicenseTableProps = {
  data: AppInfo
  show: boolean
  isAdd: boolean
  setIsAdd: Function
  setShow: Function
  refresh: Function
}

const StyledWrapper = styled.div`
  width: 100%;
  height: calc(100vh - 250px);
  .body {
    width: 100%;
    overflow: auto;
  }
`

function LicenseTable(props: LicenseTableProps) {
  const { data, show, setShow, refresh } = props
  const [visible, setVisible] = useState(false)
  const [currentLicense, setCurrentLicense] = useState<any>({})
    useState<any>([])
  const [height, setHeight] = useState(800)
  const ref = useRef<HTMLDivElement>()

  useEffect(() => {
    const resizeObserver = new ResizeObserver(entries => {
      for (let entry of entries) {
        setHeight(entry.contentRect.height)
      }
    })

    resizeObserver.observe(ref.current)

    setTimeout(() => {
      ref.current.style.paddingRight = 1 + 'px'
    }, 3000)

    return () => {
      resizeObserver && resizeObserver.disconnect()
    }
  }, [])
  const onClose = () => {
    setShow(false)
    setCurrentLicense({})
  }

  const fetchModuleConfigs = async id => {
    const { data } = await Http.get(`/licenseInfos/${id}/moduleConfigs`)
    return data
  }
  async function goToQuantityDetail(id) {
    const res = await fetchModuleConfigs(id)

    await Modal.show({
      title: '数量详情',
      footer: null,
      width: 1000,
      bodyStyle: {
        height: 600
      },
      content: ({ onCancel, onOk }) => {
        return (
          <QuantityDetail
            licenseId={id}
            moduleInfos={res?.module_config_infos || []}
            usedPercent={res?.used_percent}
            onOk={onOk}
          />
        )
      }
    })
  }
  async function showAppModuleReport(record) {
    Modal.show({
      title: '许可证模块使用情况',
      footer: null,
      width: 1000,
      bodyStyle: {
        height: 600
      },
      content: ({ onCancel, onOk }) => {
        return (
          <Report
            chartType={'LICENSE_APP_MODULE_USED_UT_AVG'}
            showSelect={false}
            licenseId={record.id}
            licenseType={record.license_type}
          />
        )
      }
    })
  }
  const deleteLicense = ({ id, license_name }) => {
    Modal.confirm({
      title: '删除许可证',
      content: `确认删除「${license_name}」！`,
      okText: '确认',
      visible,
      cancelText: '取消',
      onOk: async () => {
        await Http.delete(`/licenseInfos/${id}`)
          .then(res => {
            if (res.success) {
            }
          })
          .finally(() => {
            setVisible(false)
          })
        refresh()
      }
    })
  }

  const columns = [
    {
      header: '许可证名称',
      dataKey: 'license_name',
      props: {
        width: 180,
        fixed: true,
        resizable: true
      }
    },
    // {
    //   title: '是否授权',
    //   dataIndex: 'auth',
    //   key: 'auth',
    //   width: 100,
    //   render: (text, record) => <>{text === true ? '是' : '否'}</>
    // },
    {
      header: '许可证服务器',
      dataKey: 'license_url',
      props: {
        width: 200
      }
    },
    {
      header: '端口',
      dataKey: 'port'
    },
    {
      header: '工具路径',
      dataKey: 'tool_path',
      props: {
        width: 200,
      },
      cell: {
        render: ({ rowData }) => (
          <Tooltip title={rowData.tool_path} placement='topLeft'>
            {rowData.tool_path}
          </Tooltip>
        )
      }
    },
    {
      header: 'Mac地址',
      dataKey: 'mac_addr',
      props: {
        width: 200,
      }    
    },
    {
      header: '模块总数',
      dataKey: 'module_config_infos',
      cell: {
        render: ({ rowData }) => (
          <>{rowData?.module_config_infos?.length || 0}</>
        )
      }
    },
    {
      header: '调度优先级',
      dataKey: 'weight',
      props: {
        width: 150,
      }
    },
    {
      header: '使用有效期',
      dataKey: 'time',
      props: {
        width: 350
      },
      cell: {
        render: ({ rowData }) => (
          <>
            {rowData?.begin_time} - {rowData?.end_time}
          </>
        )
      }
    },
    {
      header: '操作',
      dataKey: 'operation',
      props: {
        width: 220,
        fixed: 'right' as 'right'
      },
      cell: {
        render: ({ rowData }) => {
          return (
            <>
              <Button
                type='link'
                onClick={() => {
                  goToQuantityDetail(rowData?.id)
                }}>
                数据
              </Button>
              <Button type='link' onClick={() => showAppModuleReport(rowData)}>
                图表
              </Button>
              <Button
                type='link'
                onClick={() => {
                  setShow(true)
                  setCurrentLicense(rowData)
                }}>
                编辑
              </Button>
              <Button type='link' onClick={() => deleteLicense(rowData)}>
                删除
              </Button>
            </>
          )
        }
      }
    }
  ]
  return (
    <StyledWrapper ref={ref}>
      <Drawer
        width={860}
        title={`${currentLicense?.license_name ? '编辑' : '添加'}许可证`}
        placement='right'
        onClose={onClose}
        destroyOnClose={true}
        visible={show}>
        <ConfigForm
          onCancel={() => onClose()}
          onSubmit={values => {
            if (values === 'ok') {
              refresh()
              onClose()
            }
          }}
          licenseConfig={{
            id: data?.id,
            ...currentLicense,
            begin_time: currentLicense?.begin_time
              ? moment(currentLicense?.begin_time)
              : '',
            end_time: currentLicense?.end_time
              ? moment(currentLicense?.end_time)
              : ''
          }}
        />
      </Drawer>
      <div className='body'>
      <Table
        columns={columns}
        props={{
          data: data?.license_infos || [],
          height: height-200,
          rowKey: 'id'
        }}
      />
      </div>
      
    </StyledWrapper>
  )
}

export default function LicenseDetail() {
  const { id } = useParams<{ id: string }>()
  const [license, setLicense] = useState(null)
  const [show, setShow] = useState(false)
  const [isAdd, setIsAdd] = useState(true)
  async function fetch() {
    const { data } = await Http.get(`/licenseManagers/${id}`)
    setLicense(new AppInfo(data))
  }
  useEffect(() => {
    fetch()
  }, [])
  const renderTitle = () => {
    return (
      <BackButton
        title='返回许可证管理'
        onClick={() => window.history.back()}
        style={{
          fontSize: 20
        }}>
        应用信息
      </BackButton>
    )
  }

  return (
      <StyledLayout>
        <Descriptions title={renderTitle()}>
          <Descriptions.Item label='编号'>{license?.id}</Descriptions.Item>
          <Descriptions.Item label='许可证类型'>
            {license?.app_type}
          </Descriptions.Item>
          <Descriptions.Item label='操作系统'>
            {license?.os_name}
          </Descriptions.Item>
          <Descriptions.Item label='创建时间'>
            {license?.create_time}
          </Descriptions.Item>
          <Descriptions.Item label='描述'>{license?.desc}</Descriptions.Item>
          <Descriptions.Item label='计算规则'>
            {ellipsisModal(license?.compute_rule)}
          </Descriptions.Item>
        </Descriptions>
        <ListWrapper>
          <h3>许可证列表</h3>
          <Button
            icon='add'
            type='primary'
            onClick={() => {
              setIsAdd(true)
              setShow(true)
            }}>
            添加
          </Button>
          <LicenseTable
            data={license}
            show={show}
            setShow={setShow}
            refresh={fetch}
            isAdd={isAdd}
            setIsAdd={setIsAdd}
          />
        </ListWrapper>
      </StyledLayout>
  )
}
