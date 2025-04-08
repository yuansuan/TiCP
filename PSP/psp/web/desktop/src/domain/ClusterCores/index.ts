import { observable, runInAction } from 'mobx'
import { Http } from '@/utils'

class ClusterCores {
  @observable total_cores = 0
  @observable available_cores = 0

  async getClusterCoreInfo() {
    const res = await Http.get('/node/coreNum')

    runInAction(() => {
      this.total_cores = res.data?.total_num || 0
      this.available_cores = res.data?.free_num || 0
    })

    return res
  }
}

export const clusterCores = new ClusterCores()
