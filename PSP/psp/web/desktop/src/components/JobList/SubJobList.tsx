import * as React from 'react'
import { observer, Provider, disposeOnUnmount } from 'mobx-react'
import { observable, reaction } from 'mobx'
import List from './List'
import ListPagination from './ListPagination'
import ListActions from './ListActions'
import { arrayJobList, arrayHistoryJobList } from '@/domain/JobList'

interface IProps {
  jobId?: number
  isHistory: boolean
  workspaceId?: number
}

@observer
export default class SubJobList extends React.Component<IProps> {
  JobListModal = this.props.isHistory ? arrayHistoryJobList : arrayJobList

  @observable loading = false

  async componentDidMount() {
    // this.JobListModal.clearStatus()
    await Promise.all([
      this.updateJobList(),
      this.JobListModal.fetchQueueList(),
    ])
  }

  // 更新job list的信息
  updateJobList = async () => {
    this.loading = true
    await this.JobListModal.fetchList(this.options)
    this.loading = false
  }

  @disposeOnUnmount
  disposer = reaction(
    () => this.options,
    () => {
      this.updateJobList()
    }
  )

  get options() {
    const {
      currentIndex,
      job_id,
      pageSize,
      stateFilter,
      userFilter,
      queueFilter,
      appFilter,
      submitPastTime,
      endPastTime,
      fuzzy,
      orderBy,
      orderAsc,
    } = this.JobListModal

    return {
      currentIndex,
      pageSize,
      stateFilter,
      userFilter,
      queueFilter,
      appFilter,
      submitPastTime,
      endPastTime,
      fuzzy,
      job_id: this.props.jobId || job_id,
      orderBy,
      orderAsc,
      workspaceFilter: this.props.workspaceId || -1, // 如果 workspaceId 存在，说明是需要通过 workspaceId 过滤 List
    }
  }

  render() {
    let columns = [
      'id',
      'jobName',
      'jobStatus',
      'app',
      'submitTime',
      'startTime',
      'endTime',
      'userName',
      'queue',
      'exHostStr',
      'numExecProcs',
      'cpuTime',
    ]

    return (
      <Provider store={this.JobListModal}>
        <ListActions
          jobList={this.JobListModal.list}
          hasActions={{
            buttons: true,
            buttonsOption: {
              reSubmit: false,
            },
            search: false,
            moreFilter: false,
          }}
          isDropDown={true}
          queueList={this.JobListModal.queueList}
          updateJobList={this.updateJobList}
        />
        <List
          workspaceId={this.props.workspaceId || -1}
          isSubJob={true}
          reSubmit={false}
          isHistory={this.props.isHistory}
          jobList={this.JobListModal.list}
          columns={columns}
          isPopupJobDetail={true}
          loading={this.loading}
          hasCheckbox={true}
        />
        <ListPagination totalItems={this.JobListModal.totalItems} />
      </Provider>
    )
  }
}
