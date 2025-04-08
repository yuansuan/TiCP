import * as React from 'react'
import { observer } from 'mobx-react'

import { Wrapper } from './style'
import { List } from './List'
import { userLogList } from '@/domain/UserLog'

@observer
export default class UserLogPage extends React.Component<any> {
  render() {
    return (
      <Wrapper>
        <div className='body'>
          <List logList={userLogList} />
        </div>
      </Wrapper>
    )
  }
}
