import * as React from 'react'
import { observer, Provider, disposeOnUnmount } from 'mobx-react'
import { observable, reaction } from 'mobx'
import List from './List'
import ListPagination from './ListPagination'
import ListActions from './ListActions'
import { historyJobList as JobListModal } from '@/domain/JobList'
import { currentUser } from '@/domain'
interface IProps {
  nodeName?: string
  workspaceId?: number
}

@observer
export default class HistoryJobList extends React.Component<IProps> {
  @observable loading = false

  get hasViewAllJobPermission() {
    return currentUser.perms.includes('system-view_all_job')
  }

  async componentDidMount() {
    JobListModal.clearStatus()
    await Promise.all([this.updateJobList(), JobListModal.fetchQueueList()])
  }

  // 更新job list的信息
  updateJobList = async () => {
    this.loading = true
    await JobListModal.fetchList(this.options)
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
      pageSize,
      stateFilter,
      userFilter,
      queueFilter,
      appFilter,
      submitPastTime,
      endPastTime,
      fuzzy,
      nodeName,
      orderBy,
      orderAsc,
    } = JobListModal

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
      nodeName: this.props.nodeName || nodeName, // 如果 nodeName 存在，说明是需要通过 nodeName 过滤 List
      orderBy,
      orderAsc,
      workspaceFilter: this.props.workspaceId || -1,
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
      <Provider store={JobListModal}>
        <ListActions
          hasUserFilter={true}
          isHistory={true}
          jobList={JobListModal.list}
          hasActions={{
            buttons: true,
            buttonsOption: {
              reSubmit: true,
            },
            search: true,
            moreFilter: true,
          }}
          queueList={JobListModal.queueList}
          updateJobList={this.updateJobList}
        />
        <List
          workspaceId={this.props.workspaceId || -1}
          isHistory={true}
          reSubmit={true}
          jobList={JobListModal.list}
          columns={columns}
          isPopupJobDetail={false}
          loading={this.loading}
          hasCheckbox={true}
        />
        <ListPagination totalItems={JobListModal.totalItems} />
      </Provider>
    )
  }
}
