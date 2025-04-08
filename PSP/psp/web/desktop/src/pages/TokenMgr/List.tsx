import React from 'react'
import { Table, Button, Modal } from '@/components'
import styled from 'styled-components'
import { env, useResize } from '@/domain'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Http as v2Http } from '@/utils/v2Http'
import { message } from 'antd'
import moment from 'moment'

const StyledLayout = styled.div``

interface IProps {
  model: { list: any[] }
  action: { fetch: () => void }
  loading: boolean
}

export const List = observer(function List({
  model,
  loading,
  action: { fetch }
}: IProps) {
  const [rect, ref] = useResize()
  const { dataSource } = useLocalStore(() => ({
    get dataSource() {
      return model.list.map(item => ({
        ...item,
        create_time: moment(item.createTime).format('YYYY/MM/DD HH:mm:ss')
      }))
    }
  }))

  const onRemove = async ({ token }) => {
    await Modal.showConfirm({
      title: '确认删除',
      content: `确认删除这条密钥？`
    })
    await v2Http.delete(`/kms/token/${env.project?.id}/${token}`)
    fetch()
    message.success('删除成功')
  }

  return (
    <StyledLayout ref={ref}>
      <Table
        columns={[
          {
            header: '密钥',
            dataKey: 'token'
          },
          {
            header: '创建时间',
            dataKey: 'create_time'
          },
          {
            header: '操作',
            props: {
              width: '20%'
            },
            cell: {
              props: {
                dataKey: 'token'
              },
              render: ({ rowData, dataKey }) => (
                <Button
                  type='link'
                  onClick={() =>
                    onRemove({
                      token: rowData[dataKey]
                    })
                  }>
                  删除
                </Button>
              )
            }
          }
        ]}
        props={{
          rowKey: 'token',
          width: rect.width,
          autoHeight: true,
          data: dataSource,
          loading
        }}
      />
    </StyledLayout>
  )
})
