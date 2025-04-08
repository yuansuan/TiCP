import { runInAction } from 'mobx'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { Http } from '@/utils'
import { SoftwareList, HardwareList } from '@/domain/VIsIBV'
import { SessionList } from '@/domain/SessionList'
import { currentUser } from '@/domain'
import { pageStateStore } from '@/utils'
import { vis } from '@/domain'

type Params = {
  statuses: string[]
  user_name: string
  hardware_ids: string[]
  software_ids: string[]
  project_ids: string[]
  page_index: number
  page_size: number
  is_admin: boolean
}
export function useModel() {
  // 暂时不支持 pageStateStore， 之后如果考虑的话，恢复一下就好
  const sessionQuery = {} as any // pageStateStore.getByPath('sessionList') as any

  return useLocalStore(() => ({
    model: new SessionList(),
    fetching: false,
    setFetching(flag) {
      this.fetching = flag
    },
    statuses: sessionQuery?.statuses || [],
    user_name: sessionQuery?.user_name,
    project_ids: sessionQuery?.project_ids || [],
    hardware_ids: sessionQuery?.hardware_ids || [],
    software_ids: sessionQuery?.software_ids || [],
    page_index: sessionQuery?.page_index || 1,
    setSessionPageIndex(index) {
      this.page_index = index
    },
    page_size: sessionQuery?.page_size || 10,
    setSessionPageSize(size) {
      this.page_size = size
    },
    changeSearchItems(changedValues, allValues) {
      Object.assign(this, allValues)
    },

    get params(): Params {
      return {
        statuses: this.statuses,
        user_name: '',
        project_ids: this.project_ids,
        hardware_ids: this.hardware_ids,
        software_ids: this.software_ids,
        page_index: this.page_index,
        page_size: this.page_size,
        is_admin: currentUser.hasSysMgrPerm,
      }
    },

    async fetchSessionList() {
      try {
        this.setFetching(true)
        pageStateStore.setByPath('sessionList', this.params)
        await this.model.fetch(this.params)
      } finally {
        this.setFetching(false)
      }
    },
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    tabType: '1',
    setTabType(value) {
      this.tabType = value
    },
    pageIndex: 1,
    setPageIndex(index) {
      this.pageIndex = index
    },
    pageSize: 10,
    setPageSize(size) {
      this.pageSize = size
    },
    name: '',
    setName(name) {
      this.name = name
    },
    software: new SoftwareList(),
    hardware: new HardwareList(),
    async refreshSoftware(size?: number) {
      this.setLoading(true)
      try {
        const { data } = await Http.get('/vis/software', {
          params: {
            name: this.name,
            platform: '',
            is_admin: currentUser.hasSysMgrPerm,
            page_index: this.pageIndex,
            page_size: size === 0 ? 0 : this.pageSize
          }
        })

        runInAction(() => {
          this.software.update({
            list: data.softwares || [],
            page_ctx: {
              index: 1,
              size: 10,
              total: data.total
            }
          })
        })
      } finally {
        this.setLoading(false)
      }
    },

    async refreshHardware(size?: number) {
      this.setLoading(true)
      try {
        const { data } = await Http.get('/vis/hardware', {
          params: {
            name: this.name,
            is_admin: currentUser.hasSysMgrPerm,
            page_index: this.pageIndex,
            page_size: size === 0 ? 0 : this.pageSize
          }
        })

        runInAction(() => {
          this.hardware.update({
            list: data.hardwares || [],
            page_ctx: {
              index: 1,
              size: 10,
              total: data.total
            }
          })
        })
      } finally {
        this.setLoading(false)
      }
    },

    projects: [],
    setProjects (vals) {
      this.projects = vals
    },
    async fetchProjects() {
      const res = await vis.getProjects(currentUser.hasSysMgrPerm)
      this.setProjects(res?.data?.projects?.map(item => ({
        key: item?.id, 
        name: item?.name
      })) || [])
    }

  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
