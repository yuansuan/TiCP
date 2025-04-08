/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { useLocalStore } from 'mobx-react-lite'
import { createStore } from '@/utils/reducer'
import { runInAction } from 'mobx'
import { CompanyUsers, UserQueryOrderBy } from '@/domain/CompanyUsers'
import { companyServer, departmentServer } from '@/server'
import { env } from '@/domain'
import { DepartmentList } from '@/domain/DepartmentList'

export const useModel = () =>
  useLocalStore(() => ({
    loading: true,
    setLoading(bool) {
      this.loading = bool
    },

    departmentList: new DepartmentList(),
    setDepartmentList(list) {
      this.departmentList = list
    },

    members: new CompanyUsers(),
    setMember(v) {
      this.members = new CompanyUsers(v)
    },

    selectedKeys: [],
    setSelectedKeys(arr) {
      this.selectedKeys = arr
    },

    currentUserRole: '',
    setCurrentUserRole(role) {
      this.currentUserRole = role
    },

    page_index: 1,
    page_size: 10,
    total: 0,
    setPage(index, size?, total?) {
      this.page_index = index
      this.setSelectedKeys([])
      if (size !== undefined) this.page_size = size
      if (total !== undefined) this.total = total
    },

    query: {
      key: '',
      sortKey: '',
      sortType: '',
      company_id: '',
      order_by: UserQueryOrderBy.ORDERBY_NULL,
      status: 1
    },
    setQuery(query) {
      this.query = { ...this.query, ...query } as any
      this.setPage(1)
    },

    async fetch() {
      try {
        this.setLoading(true)

        let orderBy = UserQueryOrderBy.ORDERBY_NULL
        if (this.query.sortKey === 'create_time') {
          if (this.query.sortType === 'asc') {
            orderBy = UserQueryOrderBy.ORDERBY_JOINTIMEASC
          } else {
            orderBy = UserQueryOrderBy.ORDERBY_JOINTIMEDESC
          }
        } else if (this.query.sortKey === 'last_login_time') {
          if (this.query.sortType === 'asc') {
            orderBy = UserQueryOrderBy.ORDERBY_LASTLOGINTIMEDASC
          } else {
            orderBy = UserQueryOrderBy.ORDERBY_LASTLOGINTIMEDESC
          }
        }

        const depData = await departmentServer.getList({
          company_id: env.company.id,
          status: 1,
          page_index: 1,
          page_size: 1000
        })
        const depList = depData?.data?.list

        const { data } = await companyServer.getUserRole()

        const {
          data: { list, page_ctx }
        } = await companyServer.queryUsers({
          page_index: this.page_index,
          page_size: this.page_size,
          ...this.query,
          order_by: orderBy,
          company_id: env.company.id
        })
        runInAction(() => {
          this.setPage(page_ctx.index, page_ctx.size, page_ctx.total)
          this.setMember({ list })
          this.setCurrentUserRole(data?.role_list[0].name)
          this.setDepartmentList(depList)
          this.setLoading(false)
        })
      } catch (e) {
        runInAction(() => {
          this.setLoading(false)
        })
      }
    }
  }))

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
