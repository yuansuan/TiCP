import { monitorServer } from '@/server/monitorServer'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { JobStatusList, ClusterInfo } from '@/domain/FullScreenMonitor'
import { runInAction } from 'mobx'

export const useModel = () => {
  return useLocalStore(() => {
    return {
      clusterInfo: new ClusterInfo(),
      setClusterInfo(clusterInfo) {
        this.clusterInfo.update(clusterInfo)
      },

      jobStatusList: new JobStatusList(),
      setJobStatusList(list) {
        this.jobStatusList.update(list)
      },

      cpuAvgList: [],
      setCpuAvgList(list) {
        this.cpuAvgList = list
      },

      memAvgList: [],
      setMemAvgList(list) {
        this.memAvgList = list
      },

      ioAvgList: [],
      setIoAvgList(list) {
        this.ioAvgList = list
      },

      diskList: { data: [], fields: [] },
      setDiskList(list) {
        this.diskList = list
      },

      nodeList: [],
      setNodeList(list) {
        this.nodeList = list
      },

      userJobList: [],
      setUserJobList(list) {
        this.userJobList = list
      },
      appJobList: { app_jobs: [], app_total: 0 },
      setAppJobList(list) {
        this.appJobList = list
      },
      projectJobList: { projects: [], users: [], jobs: [], cputimes: [] },
      setProjectJobList(list) {
        this.projectJobList = list
      },

      async refresh() {
        let endTime = Date.now()
        let startTime = endTime - 24 * 60 * 60 * 1000

        const { data: jobStateData } = await monitorServer.jobStateData(
          String(startTime),
          String(endTime)
        )
        const { data: userJobData } = await monitorServer.userJobData(
          String(startTime),
          String(endTime)
        )
        const { data: appJobData } = await monitorServer.appJobData(
          String(startTime),
          String(endTime)
        )
        const { data: projectJobData } = await monitorServer.projectJobData(
          String(startTime),
          String(endTime)
        )
        const { data: resInfoData } = await monitorServer.resourceInfo(
          String(startTime),
          String(endTime)
        )
        const { data: clusterInfoData } = await monitorServer.clusterInfo()

        runInAction(() => {
          this.setClusterInfo(clusterInfoData?.clusterInfo)
          this.setJobStatusList({ list: jobStateData?.jobResLatest })
          this.setCpuAvgList(resInfoData.metric_cpu_ut_avg?.[0]?.d)
          this.setMemAvgList(resInfoData.metric_mem_ut_avg?.[0]?.d)
          this.setIoAvgList(resInfoData.metric_io_ut_avg)
          this.setDiskList(clusterInfoData?.disks)
          this.setNodeList(clusterInfoData?.nodeList)
          this.setUserJobList(userJobData)
          this.setAppJobList(appJobData)
          this.setProjectJobList(projectJobData)
        })
      }
    }
  })
}

const store = createStore(useModel)

export const Context = store.Context
export const Provider = store.Provider
export const useStore = store.useStore
