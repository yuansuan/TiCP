/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { BillUserList } from './Model/index'
import { runInAction } from 'mobx'
import { billUserServer } from '@/server'
import { env, currentUser } from '@/domain'

type Store = {
  model: BillUserList
  loading: boolean
  queryKey: Partial<QueryKey>
  pageIndex: number
  pageSize: number
}

type QueryKey = {
  types: number[]
  merchandise_id: string
  billing_month: string
}

export function useModel() {
  const store = useLocalStore(() => ({
    model: new BillUserList(),
    loading: false,
    queryKey: {
      types: [],
      merchandise_id: '',
      billing_month: ''
    },
    pageIndex: 1,
    pageSize: 10,
    update(data: Partial<Store>) {
      Object.assign(store, data)
    },
    expandedRowKeys: [],
    setExpandedRowKey(val) {
      this.expandedRowKeys = val
    },
    async fetch() {
      try {
        store.update({
          loading: true
        })
        const { list, page_ctx, total_amount, total_refund_amount } =
          await billUserServer.getBillUserList({
            ...store.queryKey,
            project_id: env.project.id,
            user_id: currentUser?.id,
            company_id: !env.isPersonal ? env?.company?.id : '1',
            pageIndex: this.pageIndex,
            pageSize: this.pageSize
          })
        runInAction(() => {
          this.model.update({
            list: list.map(({ job, bill, merchandise_name }) => ({
              ...bill,
              bill_id: bill.id,
              bill_job_id: bill.job_id,
              job_id: job.job.id,
              merchandise_name
            })),
            page_ctx,
            total_amount,
            total_refund_amount
          })
        })
      } finally {
        store.update({
          loading: false
        })
      }
    }
  }))

  return store
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
