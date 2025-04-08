import * as React from 'react'
import { observer, Provider, disposeOnUnmount } from 'mobx-react'
import { observable, reaction } from 'mobx'
import List from './List'
import ListPagination from './ListPagination'
import ListActions from './ListActions'
import { hostJobList } from '@/domain/JobList'

interface IProps {
  nodeName?: string
}

@observer
export default class JobListByHost extends React.Component<IProps> {
  JobListModal = hostJobList

  @observable loading = false

  async componentDidMount() {
    // this.JobListModal.clearStatus()
    await Promise.all([this.updateJobList()])
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
      nodeName: this.props.nodeName,
      job_id,
      orderBy,
      orderAsc,
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
              reSubmit: true,
            },
            search: false,
            moreFilter: false,
          }}
          isDropDown={true}
          queueList={this.JobListModal.queueList}
          updateJobList={this.updateJobList}
        />
        <List
          isSubJob={false}
          reSubmit={true}
          isHistory={false}
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
