/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState, useEffect, useMemo } from 'react'
import { ListActionWrapper, ModalListDataWrapper } from '../style'
import {
  HARDWARE_DISPLAY_MAP,
  HARDWARE_PLATFORM_MAP,
  Software
} from '@/domain/Vis'
import { Spin, Select, Tooltip, Table } from 'antd'
import { observer } from 'mobx-react-lite'
import { propertyMapReduce } from '../utils'
import { Content, Layout, Sider, SiderTitle } from './styles'
const selectStyles = { width: '200px' }
interface IProps {
  loading: boolean
  software: Array<Software>
  onSelect?: (softwareId: string) => void
}

export const SoftwareSelector = observer(
  ({ software, loading, onSelect }: IProps) => {
    const [platform, setPlatform] = useState(0)
    const [display, setDispay] = useState(0)
    const [softwareId, setSoftwareId] = useState('')
    const [filteredSoftware, setFilteredSoftware] = useState(software)
    const [selectRowKeys, setSelectRowKeys] = useState([])

    useEffect(() => {
      setFilteredSoftware(
        software.filter(item => {
          if (display !== 0 && item.display != display) return false
          if (platform !== 0 && item.platform != platform) return false
          if (softwareId !== '' && item.id !== softwareId) return false
          return true
        })
      )
    }, [platform, softwareId, display])

    useEffect(() => {
      onSelect(selectRowKeys.toString())
    }, [selectRowKeys])

    const columns = useMemo(() => {
      return [
        {
          title: '软件名称',
          dataIndex: 'name',
          width: 200,
          render: (_, record) => {
            return (
              <Tooltip placement='topLeft' title={record.name}>
                {record.name}
              </Tooltip>
            )
          }
        },

        {
          title: '软件描述',
          dataIndex: 'desc',
          width: 200,
          onCell: () => {
            return {
              style: {
                maxWidth: 200,
                overflow: 'hidden',
                whiteSpace: 'nowrap',
                textOverflow: 'ellipsis',
                cursor: 'pointer'
              }
            }
          },
          render: (_, record) => (
            <Tooltip placement='topLeft' title={record.desc}>
              {record.desc}
            </Tooltip>
          )
        },
        {
          title: '软件平台',
          dataIndex: 'platform',
          width: 200,
          render: (_, record) => {
            return (
              <Tooltip placement='topLeft' title={record.platform}>
                {record.platform}
              </Tooltip>
            )
          }
        }
        // {
        //   title: '软件形式',
        //   dataIndex: 'display',
        //   width: 100,
        //   render: (_, record) => {
        //     return (
        //       <Tooltip
        //         placement='topLeft'
        //         title={HARDWARE_DISPLAY_MAP[record.display]}>
        //         {HARDWARE_DISPLAY_MAP[record.display]}
        //       </Tooltip>
        //     )
        //   }
        // }
      ]
    }, [])
    const rowSelectionChange = (selectRowKeys, rowSelection) => {
      setSelectRowKeys(selectRowKeys)
      if (selectRowKeys) {
        document.querySelector('.validate_soft_tip').innerHTML = ''
      }
    }
    const rowSelection: any = {
      type: 'radio',
      onChange: rowSelectionChange,
      selectedRowKeys: selectRowKeys
    }
    return (
      <Layout>
        <Sider>
          <SiderTitle>选择软件</SiderTitle>
        </Sider>
        <Content>
          <Spin spinning={loading}>
            <ListActionWrapper style={{ padding: 0 }}>
              <div className='item'>
                <Select
                  className={'status'}
                  value={platform}
                  onChange={setPlatform}
                  style={selectStyles}>
                  <Select.Option value={0}>所有平台</Select.Option>
                  {propertyMapReduce(software, 'platform').map(v => (
                    <Select.Option value={v} key={v}>
                      <Tooltip placement='topLeft' title={v}>
                        {v}
                      </Tooltip>
                    </Select.Option>
                  ))}
                </Select>
              </div>
              {/* <div className='item'>
                <Select
                  className={'status'}
                  value={display}
                  onChange={setDispay}
                  style={selectStyles}>
                  <Select.Option value={0}>所有形式</Select.Option>
                  {propertyMapReduce(software, 'display').map(v => (
                    <Select.Option value={v} key={v}>
                      <Tooltip
                        placement='topLeft'
                        title={v}>
                        {v}
                      </Tooltip>
                    </Select.Option>
                  ))}
                </Select>
              </div> */}
              <div className='item'>
                <Select
                  className={'status'}
                  value={softwareId}
                  onChange={setSoftwareId}
                  style={selectStyles}>
                  <Select.Option value={''}>所有软件</Select.Option>
                  {software.map(v => (
                    <Select.Option value={v.id} key={v.id}>
                      <Tooltip placement='topLeft' title={v.name}>
                        {v.name}
                      </Tooltip>
                    </Select.Option>
                  ))}
                </Select>
              </div>
            </ListActionWrapper>
            <ModalListDataWrapper style={{ padding: 0, marginTop: '10px' }}>
              <Table
                rowKey={record => record.id}
                size='small'
                onRow={record => {
                  return {
                    onClick: () => {
                      rowSelectionChange([record.id], [record])
                    }
                  }
                }}
                rowSelection={{ ...rowSelection }}
                dataSource={filteredSoftware}
                columns={columns}
                pagination={false}></Table>
            </ModalListDataWrapper>
          </Spin>
          {software && (
            <div className='validate_tip validate_soft_tip '>
              {'请选择软件规格'}
            </div>
          )}
        </Content>
      </Layout>
    )
  }
)
