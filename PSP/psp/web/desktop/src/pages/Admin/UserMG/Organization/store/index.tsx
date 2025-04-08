import { useCallback } from 'react'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import organization from '@/domain/UserMG/UserOfOrgList'

export function useModel() {
  return useLocalStore(() => ({
    userList: null,
    setUserList(list) {
      this.userList = list
    },
    totalItems: 0,
    setTotalItems(total) {
      this.totalItems = total
    },
    org: null,
    setOrg(org) {
      this.org = org
    },
    nodeId: undefined,
    setNodeId(id) {
      this.nodeId = id
    },

    searchKey: '',
    setSearchKey(key) {
      this.setPage(1, 10)
      this.searchKey = key
    },
    orderAsc: true,
    orderBy: 'name',
    setOrder(sortType, sortKey) {
      this.orderAsc = sortType
      this.orderBy = sortKey
    },

    page_index: 1,
    page_size: 10,
    total: 0,
    setPage(index, size?, total?) {
      this.page_index = index
      if (size != undefined) this.page_size = size
      if (total != undefined) this.total = total
    },
    selectedKeys: [],
    setSelectedKeys(keys) {
      this.selectedKeys = [...keys]
    },

    listLoading: false,
    setListLoading(flag) {
      this.listLoading = flag
    },
    orgLoading: false,
    setOrgLoading(flag) {
      this.orgLoading = flag
    },

    getUserList: function getUserList(): [() => Promise<void>, boolean] {
      const store = this
      const fetch = useCallback(
        async function fetch() {
          const {
            nodeId,
            page_index,
            page_size,
            searchKey,
            orderAsc,
            orderBy,
          } = store

          if (!nodeId) {
            return
          }
          try {
            store.setListLoading(true)

            const {
              data: { list, page_ctx },
            } = await organization.getUserList({
              id: nodeId,
              index: page_index,
              size: page_size,
              orderAsc: orderAsc,
              orderBy: orderBy,
              search_value: searchKey,
            })
            store.setUserList(organization.userList)
            store.setSelectedKeys([])
            store.setPage(page_ctx.index, page_ctx.size, page_ctx.total)
          } finally {
            store.setListLoading(false)
          }
        },
        [
          this.nodeId,
          this.page_index,
          this.page_size,
          this.searchKey,
          this.orderAsc,
          this.orderBy,
        ]
      )
      return [fetch, this.listLoading]
    },
    getOrganization: function getOrganization(): [
      () => Promise<void>,
      boolean
    ] {
      const store = this
      const fetch = useCallback(async function fetch() {
        try {
          store.setOrgLoading(true)
          const org = await organization.getOrganization()
          store.setOrg(org)
          store.setNodeId([org?.data.org.id])
        } finally {
          store.setOrgLoading(false)
        }
      }, [])
      return [fetch, this.orgLoading]
    },
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
