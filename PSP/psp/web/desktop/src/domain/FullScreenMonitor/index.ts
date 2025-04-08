import { action, observable } from 'mobx'
import { JobStatus } from './Job'

export class JobStatusList {
  @observable list: JobStatus[]

  @action
  update({ list }) {
    if (list) {
      this.list = list.map(item => new JobStatus(item))
    }
  }
}

export class ClusterInfo {
  @observable clusterName: string
  @observable cores: number
  @observable freeCores: number
  @observable usedCores: number
  @observable totalNodeNum: number
  @observable availableNodeNum: number

  constructor(props?: Partial<ClusterInfo>) {
    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<ClusterInfo>) {
    Object.assign(this, props)
  }
}

