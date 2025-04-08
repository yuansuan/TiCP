/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import * as React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import VirtualMachine from '@/domain/Visual/VirtualMachine'
import { Descriptions } from 'antd'
import { vmList } from '@/domain/Visual'
import moment from 'moment'
const Wrapper = styled.div`
  .update_time {
    margin-top: 40px;
  }
`

@observer
export class VMDetail extends React.Component<{ vm: VirtualMachine }> {
  render() {
    let vm = this.props.vm
    vm = vmList.get(vm.id, !!vm.children)
    const updateTime = new Date(vm.update_time)
    return (
      <Wrapper>
        <Descriptions title={vm.name}>
          <Descriptions.Item label='CPU使用率'>
            {Math.round(vm.cpu_usage) + '%'}
          </Descriptions.Item>
          <Descriptions.Item label='内存使用率'>
            {Math.round(vm.mem_usage) + '%'}
          </Descriptions.Item>
          <Descriptions.Item label='硬盘使用率'>
            {Math.round(vm.disk_usage) + '%'}
          </Descriptions.Item>
          <Descriptions.Item label='GPU使用率'>
            {vm.gpu_usage}
          </Descriptions.Item>
          <Descriptions.Item label='GPU内存使用率'>
            {vm.gpu_memory}
          </Descriptions.Item>
        </Descriptions>
        <div className='update_time'>
          更新时间:{moment(updateTime).format('YYYY-MM-DD HH:mm:ss')}
        </div>
      </Wrapper>
    )
  }
}
