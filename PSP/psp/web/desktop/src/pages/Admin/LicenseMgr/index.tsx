import React, { useState, useEffect, useRef } from 'react'
import { observer } from 'mobx-react-lite'
import { reaction } from 'mobx'
import { lmList } from '@/domain'
import { Button, message, Tooltip, Pagination } from 'antd'
import { Modal, Table,Icon} from '@/components'
import Action from './Action'
import { history,Http } from '@/utils'
import styled from 'styled-components'
import { QuantityDetail } from './QuantityDetail'
import {BashEditor} from '@/components'

const StyledWrapper = styled.div`
  padding: 20px;
  width: 100%;
  background: #fff;
  height: calc(100vh - 150px);
  .body {
    width: 100%;
    overflow: auto;
  }
  .ellipsis {
    display: inline-block;
    width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; 
  }
  
  .pagination {
    margin: 20px;
    text-align: center;
  }
`



export const ellipsisModal = (text) => {
  const showScript = (
    <div>
      {text?.length >= 60 ? (
        <div style={{display: 'flex',alignItems: 'center' }}>
          <span className="ellipsis">{text?.substring(0, 60)}...</span>
          <a
            style={{ padding: 0 }}
            onClick={() => {
              return Modal.show({
                title: '更多内容',
                bodyStyle: { maxHeight: '86vh', overflow: 'auto',padding:0 },
                width: 800,
                footer: null,
                content: () => (
                  <BashEditor
                    width={'800px'}
                    height={'600px'}
                    code={text}
                    readOnly={true}
                  />
                )
              })
            }}>
            展示更多
          </a>
        </div>
      ) : (
        <span title={text}>{text}</span>
      )}
    </div>
  )

  return showScript
}

const expandedRowRender = (record, index, indent, expanded) => {
  async function goToQuantityDetail(licenseInfo) {
    await Modal.show({
      title: '数量详情',
      footer: null,
      width: 800,
      bodyStyle: {
        height: 600
      },
      content: ({ onCancel, onOk }) => {
        return <QuantityDetail licenseInfo={licenseInfo} onOk={onOk} />
      }
    })
  }
  const columns = [
    {
      header: '编号',
      dataKey: 'id',
      render: (text, record) => (
        <a
          onClick={() => {
            goToQuantityDetail(record)
          }}>
          {text}
        </a>
      )
    },

    {
      title: '是否被授权',
      dataIndex: 'auth',
      key: 'auth',
      render: (text, record) => <>{text === true ? '是' : '否'}</>
    },
    {
      title: '使用有效期',
      dataIndex: 'time',
      key: 'time',
      render: (text, record) => (
        <>
          {record?.begin_time.dateString} - {record?.end_time.dateString}
        </>
      )
    },
    { title: '许可证服务器', dataIndex: 'license_url', key: 'license_url' },
    { title: '端口', dataIndex: 'port', key: 'port' },
    {
      title: '工具路径',
      dataIndex: 'tool_path',
      key: 'tool_path',
      width: 200,
      textWrap: 'word-break',
      ellipsis: true,
      render: (text, record) => (
        <Tooltip title={text} placement='topLeft'>
          {text}
        </Tooltip>
      )
    },
    // { title: '超算ID', dataIndex: 'sc_id', key: 'sc_id' },
    { title: 'Mac地址', dataIndex: 'mac_addr', key: 'mac_addr' },
    {
      title: '模块总数',
      dataIndex: 'module_conf',
      key: 'module_conf',
      render: text => <>{text?.length || 0}</>
    },
    { title: '调度执行顺序', dataIndex: 'weight', key: 'weight' }
  ]
  return (
    <Table
      rowKey={'provider'}
      columns={columns}
      dataSource={record['listLicenseInfo']}
      pagination={false}
    />
  )
}

// new license mgr
const LicenseMgr = observer(function LicenseMgr() {
  const goToLicenseDetail = id => history.push(`/sys/license_mgr/${id}`)
  const goToLicenseEdit = id => history.push(`/sys/license_mgr/license/${id}`)
  const [height, setHeight] = useState(800)
  const [visible, setVisible] = useState(false)
  const [loading, setLoading] = useState(false)
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

  useEffect(() => {
    try{
      setLoading(true)
      lmList.fetch()
    }finally{
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    let disposer = reaction(
      () => ({
        index: lmList.index,
        size: lmList.size
      }),
      () => {
        lmList.fetch()
      }
    )

    return () => {
      disposer()
    }
  }, [])

  const deleteLicense = ({id,license_type}) => {
    Modal.confirm({
      title: '删除许可证管理',
      content: `删除前需要移除「${license_type}」下所有的许可证服务，确认删除？`,
      okText: '确认',
      visible,
      cancelText: '取消',
      onOk: async () => {
        await Http.delete(`/licenseManagers/${id}`)
          .then(res => {
            if (res.success) {
            }
          })
          .finally(() => {
            setVisible(false)
          })
          lmList.fetch()
      }
    })
  }
  

  const columns = [
    {
      props: {
        resizable: true,
        width: 120,
        fixed: true
      },
      header: '编号',
      dataKey: 'id',
      cell: {
        props: {
          dataKey: 'id'
        },
        render: ({ rowData, dataKey }) => (
          <a
            onClick={() => {
              goToLicenseDetail(rowData['id'])
            }}>
            {rowData['id']}
          </a>
        )
      }
    },
    {
      header: '许可证类型',
      dataKey: 'license_type',
      props: {
        resizable: true,
        width: 150,
      }
    },
    {
      header: '操作系统',
      dataKey: 'os_name',
      props: {
        resizable: true,
        width: 100,
      }
    },
    {
      header: '描述',
      dataKey: 'desc',
      props: {
        resizable: true,
        width: 300,
      }
    },
    {
      header: '计算规则',
      dataKey: 'compute_rule',
      props: {
        resizable: true,
        width: 350
      },
      cell: {
        render: ({ rowData, dataKey }) => ellipsisModal(rowData.compute_rule)
      }
    },
    {
      header: '创建时间',
      dataKey: 'create_time',
      props: {
        resizable: true,
        width: 250
      }
    },
    {
      props: {
        resizable: true,
        width: 160,
        fixed: 'right' as 'right'
      },
      header: '操作',
      dataKey: 'operation',
      cell: {
        render: ({ rowData, dataKey }) => {
          return (
            <>
             <Button
                type='link'
                onClick={() => goToLicenseEdit(rowData['id'])}>
                编辑
              </Button>
              <Button type='link' onClick={() => deleteLicense({...rowData})}>
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
      <Action
        onSearch={values => {
          const { license_type } = values
          lmList.setFilterParams(license_type)
          lmList.fetch()
        }}
        onAdd={() => {
          history.push('/sys/license_mgr-add')
        }}
      />
      <div className='body'>
      <Table
        columns={columns}
        props={{
          data: lmList.list,
          height: height-200,
          loading: loading,
          shouldUpdateScroll: false,
          locale: {
            emptyMessage: '没有许可证数据',
            loading: '数据加载中...'
          },
          rowKey: 'id'
        }}
      />
      </div>
     
      {/* 暂不支持分页 */}
      {/* <div className='pagination'>
        <Pagination
          current={lmList.index}
          pageSize={lmList.size}
          total={lmList.total}
          showSizeChanger
          onChange={lmList.onPageChange.bind(lmList)}
          onShowSizeChange={lmList.onSizeChange.bind(lmList)}
        />
      </div> */}
    </StyledWrapper>
  )
})

export default LicenseMgr
