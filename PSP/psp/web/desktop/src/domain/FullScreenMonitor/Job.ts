import { action, observable } from "mobx"

export class JobStatus {
  @observable job_count: number
  @observable status: string
  @observable timestamp: number

  constructor(props?: Partial<JobStatus>) {
    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<JobStatus>) {
    Object.assign(this, props)
  }
}
