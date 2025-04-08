/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import styled from 'styled-components'
import { Button } from '@/components'
import { Input, DatePicker } from 'antd'
import moment from 'moment'

const { RangePicker } = DatePicker

const Wrapper = styled.div`
  display: flex;
  justify-content: space-between;
  margin: 0 28px 16px 0;

  .filter {
    display: flex;
    justify-content:space-between;
    min-width: 250px;
  }
`

interface IProps {
  onAdd: () => void
  onSearch: (values) => void
}


@observer
class Action extends React.Component<IProps> {
  @observable license_type = ''
  
  add = () => {
    this.props.onAdd()
  }

  search = (e) => {
    this.license_type=e.target.value
    this.props.onSearch({
      license_type: e.target.value
    })
  }

  render() {

    return (
      <Wrapper>
        <Button icon='add' type='primary' onClick={this.add}>
          添加
        </Button>
        <div className='filter'>
          <label>
              许可证类型: <Input style={{width: 160}} value={this.license_type} onChange={e =>this.search(e)} placeholder='输入许可证类型' />
          </label>
          {/* <label>
              有效期: <RangePicker showTime value={this.time} onChange={(dates) => this.time = (dates || []) as [moment.Moment, moment.Moment]} />
          </label> */}
          {/* <Button onClick={this.search} type='primary'>查询</Button> */}
        </div>
      </Wrapper>
    )
  }
}

export default Action
