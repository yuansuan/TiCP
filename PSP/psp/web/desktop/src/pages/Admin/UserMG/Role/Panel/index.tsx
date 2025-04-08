import React, { useEffect } from 'react'
import { observer } from 'mobx-react'
import { Switch, message } from 'antd'
import { Table, Modal } from '@/components'
import { currentUser } from '@/domain'
import { Role, RoleList } from '@/domain/UserMG'
import { ListQuery } from '@/pages/Admin/UserMG/utils'
import Operators from './Operators'
import RolePreview from '../Preview'
import { PanelWrapper } from './style'
import { DataDashPlugin } from '@/utils'

interface IProps {
  width: number
  loading: boolean
  height: number

  listQuery: ListQuery
  updateTotal: (number) => void
  store: typeof RoleList
}

export default observer(function RolePanel(props: IProps) {
  const { width, loading, listQuery, store, updateTotal, height } = props

  useEffect(() => {
    updateTotal(store.totalRoles)
  }, [store.roleList])

  useEffect(() => {
    const { page, pageSize, query } = listQuery
    store.filter = {
      name_filter: query,
      page: {
        index: page,
        size: pageSize
      }
    }
    store.fetch()
  }, [listQuery.query, listQuery.page, listQuery.pageSize])

  const onSwitchChange = async (check, id) => {
    if (check) {
      await RoleList.setDefaultRole(id).then(() =>
        message.success('默认角色设置成功')
      )
    }
  }
  const columns = () => {
    let all = [
      {
        props: {
          flexGrow: 1,
          minWidth: 200
        },
        header: '角色名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <div
              title={rowData[dataKey]}
              className='roleName'
              onClick={() =>
                Modal.show({
                  title: '预览角色',
                  bodyStyle: {
                    height: 710,
                    background: '#F0F5FD',
                    overflow: 'auto'
                  },
                  width: 1130,
                  footer: null,
                  content: () => (
                    <RolePreview role={new Role(rowData.toRequest())} />
                  )
                })
              }>
              <span className='name'>{rowData[dataKey]}</span>
            </div>
          )
        }
      },
      {
        props: {
          width: width * 0.45,
          resizable: true
        },
        header: '描述',
        cell: {
          props: {
            dataKey: 'comment'
          },
          render: ({ rowData, dataKey }) => (
            <div className='comment'>
              <span className='name' title={rowData[dataKey]}>
                {rowData[dataKey]}
              </span>
            </div>
          )
        }
      },
      ...(currentUser.isLdapEnabled
        ? []
        : [
            {
              props: {
                width: width * 0.15
              },
              header: '默认角色',
              cell: {
                render: ({ rowData, dataKey }) => (
                  <Switch
                    size='small'
                    disabled={rowData.isDefault}
                    checked={rowData.isDefault}
                    onChange={check => onSwitchChange(check, rowData.id)}
                  />
                )
              }
            }
          ]),
      {
        props: {
          width: width * 0.15
        },
        header: '操作',
        cell: {
          render: ({ rowData }) => <Operators rowData={rowData} />
        }
      }
    ]
    return all
  }
  return (
    <PanelWrapper>
      <Table
        props={{
          height: height,
          data: store.roleList || [],
          rowKey: 'id',
          loading
        }}
        plugins={[new DataDashPlugin()]}
        columns={columns() as any}
      />
    </PanelWrapper>
  )
})
