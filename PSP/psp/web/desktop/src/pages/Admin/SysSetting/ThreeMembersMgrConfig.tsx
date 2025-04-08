import * as React from 'react'
import { observer } from 'mobx-react'
import { Select } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable } from 'mobx'

const { Option } = Select

@observer
export default class GlobalConfig extends React.Component<any> {
  @observable approve_user_id = null
  @observable users = [] // 指所有的安全管理员

  async componentDidMount() {
    await sysConfig?.fetchThreeMemberMgrConfig()
    this.approve_user_id = sysConfig?.threeMemberMgrConfig?.defaultApprover?.id || null

    fetch('/api/v1/user/optionList?filterPerm=8')
      .then(response => response.json())  
      .then(res => {
        const opts = res?.data?.map(d => ({id: d.key, name: d.title})) || []
        this.users = opts
      })
  }

  changeApproveUser = value => {
    this.approve_user_id = value
    if ( this.approve_user_id) {
      sysConfig.updateThreeMemberMgrConfig({
        id: value,
        name: this.users.find(u => u.id === value)?.name,
      })
    } else {
      // for clear
      sysConfig.updateThreeMemberMgrConfig({
        id: '',
        name: '',
      })
    }
    
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label align={'left'}>默认审批人(安全管理员)</Label>
          <Select
            value={this.approve_user_id}
            placeholder={'请选择默认审批人'}
            onChange={this.changeApproveUser}
            allowClear
            size={'small'}
            style={{ width: 240, marginLeft: 30 }}>
            {
              this.users.map(u => <Option key={u.id} value={u.id}>{u.name}</Option>)
            }
          </Select>
        </div>
      </ConfigWrapper>
    )
  }
}
