import * as React from 'react'
import { observer, Provider, disposeOnUnmount } from 'mobx-react'
import { observable, reaction } from 'mobx'
import List from './List'
import ListPagination from './ListPagination'
import ListActions from './ListActions'
import { currentUser } from '@/domain'
import { JobList as jobListModal } from '@/domain/JobList'

interface IProps {
  nodeName?: string
  workspaceId?: number
}

let intervalId = null

@observer
export default class JobList extends React.Component<IProps> {
  @observable loading = false

  async componentDidMount() {}

  componentDidUpdate() {}

  componentWillUnmount() {
    clearInterval(intervalId)
    intervalId = null
  }

  get hasViewAllJobPermission() {
    return currentUser.perms.includes('system-view_all_job')
  }

  // 更新job list的信息
  updateJobList = async () => {
    this.loading = true
    // await jobListModal.fetchList(this.options)
    this.loading = false
  }

  @disposeOnUnmount
  disposer = reaction(
    () => this.options,
    () => {
      this.updateJobList()
    },
    {
      delay: 600
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
      orderAsc
    } = jobListModal

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
      workspaceFilter: this.props.workspaceId || -1 // 如果 workspaceId 存在，说明是需要通过 workspaceId 过滤 List
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
      'cpuTime'
    ]

    return (
      <Provider store={jobListModal}>
        <ListActions
          hasUserFilter={this.hasViewAllJobPermission}
          // jobList={jobListModal.list}
          hasActions={{
            buttons: true,
            buttonsOption: {
              reSubmit: true
            },
            search: true,
            moreFilter: true
          }}
          // queueList={jobListModal.queueList}
          updateJobList={this.updateJobList}
        />
        <List
          workspaceId={this.props.workspaceId || -1}
          isHistory={false}
          reSubmit={true}
          // jobList={jobListModal.list}
          columns={columns}
          isPopupJobDetail={false}
          loading={this.loading}
          hasCheckbox={true}
        />
        {/* <ListPagination totalItems={jobListModal.totalItems} /> */}
      </Provider>
    )
  }
}
