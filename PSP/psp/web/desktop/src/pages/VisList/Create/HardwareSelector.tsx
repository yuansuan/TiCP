import React, { useState, useEffect, useMemo } from 'react'
import { ListActionWrapper, ModalListDataWrapper } from '../style'
import { Hardware } from '@/domain/Vis'
import { Spin, Select, Space, Tooltip, Table } from 'antd'
import { observer } from 'mobx-react-lite'
import { propertyMapReduce } from '../utils'
import { Content, Layout, Sider, SiderTitle } from './styles'

const selectStyles = { width: '200px' }

interface IProps {
  loading?: boolean
  hardware: Array<Hardware>
  onSelect: (hardwareId: string) => void
  defaultId?: String
}

export const HardwareSelector = observer(
  ({ hardware, loading, onSelect, defaultId }: IProps) => {
    const [cpu, setCpu] = useState(0)
    const [mem, setMem] = useState(0)
    const [gpu, setGpu] = useState(0)
    const [filteredHardware, setFilteredHardware] = useState(hardware)
    const [selectRowKeys, setSelectRowKeys] = useState([])

    useEffect(() => {
      setFilteredHardware(
        hardware.filter(item => {
          if (cpu != 0 && item.cpu != cpu) return false
          if (mem != 0 && item.mem != mem) return false
          if (gpu != 0 && item.gpu != gpu) return false
          return true
        })
      )
    }, [cpu, mem, gpu, hardware])

    useEffect(() => {
      onSelect(selectRowKeys.toString())
    }, [selectRowKeys])

    //软件切换后，实例自动选择软件对应的预设清单中的默认配置
    useEffect(() => {
      setSelectRowKeys([defaultId])
      if (defaultId) {
        document.querySelector('.validate_hard_tip').innerHTML = ''
      }
    }, [defaultId])

    const columns = useMemo(() => {
      return [
        {
          title: '机型',
          dataIndex: 'name',
          width: 200,
          render: (_, record) => (
            <Tooltip placement='topLeft' title={record.name}>
              {record.name}
            </Tooltip>
          )
        },
        {
          title: '规格',
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
          title: '处理器型号',
          dataIndex: '',
          width: 100,
          render: (_, record) => {
            return (
              <Space>
                <span> {record.cpu_model || '--'}</span>
              </Space>
            )
          }
        },

        {
          title: 'vCPU',
          dataIndex: 'cpu',
          width: 100,
          render: (_, record) => {
            return (
              <Space>
                <span>{record.cpu} 核</span>
              </Space>
            )
          }
        },
        {
          title: '内存',
          dataIndex: 'mem',
          width: 100,
          render: (_, record) => {
            return (
              <Space>
                <span>{record.mem} GB</span>
              </Space>
            )
          }
        },
        {
          title: 'vGPU',
          dataIndex: 'gpu',
          width: 100,
          render: (_, record) => {
            return (
              <Space>
                <span>{record.gpu + '个' + record.gpu_model} </span>
              </Space>
            )
          }
        }
      ]
    }, [])

    const rowSelectionChange = (selectRowKeys, rowSelection) => {
      setSelectRowKeys(selectRowKeys)
      if (selectRowKeys) {
        document.querySelector('.validate_hard_tip').innerHTML = ''
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
          <SiderTitle>实例规格</SiderTitle>
        </Sider>
        <Content>
          <Spin spinning={loading}>
            <ListActionWrapper style={{ padding: 0 }}>
              <div className='item'>
                <Select
                  className={'status'}
                  value={cpu}
                  onChange={setCpu}
                  style={selectStyles}>
                  <Select.Option value={0}>所有vCPU</Select.Option>
                  {propertyMapReduce(hardware, 'cpu').map(v => (
                    <Select.Option value={v} key={v}>
                      {v} vCPU
                    </Select.Option>
                  ))}
                </Select>
              </div>
              <div className='item'>
                <Select
                  className={'status'}
                  value={mem}
                  onChange={setMem}
                  style={selectStyles}>
                  <Select.Option value={0}>所有内存</Select.Option>
                  {propertyMapReduce(hardware, 'mem').map(v => (
                    <Select.Option value={v} key={v}>
                      {v} GiB
                    </Select.Option>
                  ))}
                </Select>
              </div>
              <div className='item'>
                <Select
                  className={'status'}
                  value={gpu}
                  onChange={setGpu}
                  style={selectStyles}>
                  <Select.Option value={0}>所有GPU</Select.Option>
                  {propertyMapReduce(hardware, 'gpu').map(v => (
                    <Select.Option value={v} key={v}>
                      {v} vGPU
                    </Select.Option>
                  ))}
                </Select>
              </div>
            </ListActionWrapper>
            <ModalListDataWrapper style={{ padding: 0, marginTop: '10px' }}>
              <Table
                size='small'
                rowKey={record => record.id}
                onRow={record => {
                  return {
                    onClick: () => {
                      rowSelectionChange([record.id], [record])
                    }
                  }
                }}
                rowSelection={{ ...rowSelection }}
                dataSource={filteredHardware}
                columns={columns}
                pagination={false}></Table>
            </ModalListDataWrapper>
          </Spin>
          {hardware && (
            <div className='validate_tip validate_hard_tip'>
              {'请选择实例规格'}
            </div>
          )}
        </Content>
      </Layout>
    )
  }
)
