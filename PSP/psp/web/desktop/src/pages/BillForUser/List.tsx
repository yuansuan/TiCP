/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useLayoutRect } from '@/utils/hooks'
import { Icon, Table } from '@/components'
import { formatUnixTime } from '@/utils'
import { useStore } from './store'
import { RESOURCE_TYPE } from '@/constant'

const StyledLayout = styled.div`
  padding: 20px 0;
`

export const NewJobLayout = styled.div`
  display: flex;
  align-items: center;

  .icon {
    margin-left: 6px;
  }
`

const RESOURCE_TYPE_MAP = {
  1: '计算作业',
  2: '可视化应用',
  5: '3D云应用-软件',
  6: '3D云应用-硬件',
  7: (
    <NewJobLayout>
      计算作业
      <Icon type='nys-new' className='icon' />
    </NewJobLayout>
  ),
  103: '3D云应用套餐'
}

export const List = observer(function List() {
  const store = useStore()
  const [rect, ref, resize] = useLayoutRect()
  const { dataSource } = useLocalStore(() => ({
    get dataSource() {
      return store.model.list.map(item => {
        const isAppResource =
          item?.out_resource_type === RESOURCE_TYPE.COMPUTE_APP ||
          item?.out_resource_type === RESOURCE_TYPE.STANDARD_COMPUTE_APP
        const isVisResource =
          item?.out_resource_type === RESOURCE_TYPE.VISUAL_APP

        const isVisIBVTypeResource =
          item?.out_resource_type === RESOURCE_TYPE.IBV_SOFTWARE ||
          item?.out_resource_type === RESOURCE_TYPE.IBV_HARDWARE ||
          item?.out_resource_type === RESOURCE_TYPE.CLOUD_APP_COMBO_USAGE

        const unit = isAppResource
          ? ' 核时'
          : isVisResource || isVisIBVTypeResource
          ? ' 时'
          : ''
        return {
          ...item,
          update_time: formatUnixTime(item?.update_time),
          merchandise_quantity: item?.merchandise_quantity?.toFixed(2) + unit,
          real_amount: (item?.real_amount / 100000).toFixed(2),
          refund_amount: (item?.refund_amount / 100000).toFixed(2),
          out_resource_type: RESOURCE_TYPE_MAP[item?.out_resource_type],
          job_id: isAppResource
            ? item?.bill_job_id
            : isVisResource
            ? item?.out_biz_id
            : '--'
        }
      })
    }
  }))

  useEffect(() => {
    resize()
  }, [])

  return (
    <StyledLayout ref={ref}>
      <Table
        props={{
          rowKey: 'id',
          height: 600,
          data: dataSource,
          loading: store.loading
        }}
        columns={[
          {
            header: '账单编号',
            props: {
              width: 150
            },
            dataKey: 'bill_id'
          },
          {
            header: '作业编号',
            props: {
              width: 200
            },
            dataKey: 'job_id'
          },
          {
            header: '账期',
            props: {
              width: 150
            },
            dataKey: 'billing_month'
          },
          {
            header: '姓名',
            props: {
              width: 150
            },
            dataKey: 'user_name'
          },
          {
            header: '成员ID',
            props: {
              width: 150
            },
            dataKey: 'user_id'
          },
          {
            header: '消费时间',
            props: {
              width: 250
            },
            dataKey: 'update_time'
          },
          {
            header: '商品',
            props: {
              width: 250
            },
            dataKey: 'merchandise_name'
          },
          {
            header: '账单类型',
            props: {
              width: 150
            },
            dataKey: 'out_resource_type'
          },
          {
            header: '用量',
            props: {
              width: 200
            },
            dataKey: 'merchandise_quantity'
          },
          {
            header: '退款金额（元）',
            props: {
              width: 200
            },
            dataKey: 'refund_amount'
          },
          {
            header: '账单金额（元）',
            props: {
              width: 200
            },
            dataKey: 'real_amount'
          }
        ]}
      />
    </StyledLayout>
  )
})
