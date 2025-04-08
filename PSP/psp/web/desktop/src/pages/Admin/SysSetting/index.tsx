import * as React from 'react'
import { observer } from 'mobx-react'
import AllConfig from './AllConfig'

import { Wrapper } from './style'
import { observable } from 'mobx'

@observer
export default class SysSetting extends React.Component<any> {
  @observable search = ''

  render() {
    return (
      <Wrapper>
        <div className='body'>
          <AllConfig filterKey={this.search} />
        </div>
      </Wrapper>
    )
  }
}
