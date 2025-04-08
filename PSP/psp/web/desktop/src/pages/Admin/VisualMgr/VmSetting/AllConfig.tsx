import * as React from 'react'
import { observer } from 'mobx-react'
import { Collapse, Spin } from 'antd'
import { machineSetting } from '@/domain/Visual'
import UserTaskConfig from './UserTaskConfig'
import { observable } from 'mobx'

const { Panel } = Collapse

@observer
export default class AllConfig extends React.Component<any> {
  @observable loading = true

  async componentDidMount() {
    await Promise.all([machineSetting.fetch()])
    this.loading = false
  }

  render() {
    return (
      <>
        {this.loading ? (
          <div className='loading'>
            <Spin />
          </div>
        ) : (
          <Collapse bordered={false} defaultActiveKey={['0']}>
            <Panel header='任务数量设置' key='0'>
              <UserTaskConfig />
            </Panel>
          </Collapse>
        )}
      </>
    )
  }
}
