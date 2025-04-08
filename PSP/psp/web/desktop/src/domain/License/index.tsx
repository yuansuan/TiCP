import { action, observable } from 'mobx'
import { Http } from '@/utils'
import Info from './Info'

export class License {
  @observable machineId: string
  @observable info: Info

  @action
  async getMachineId() {
    const res = await Http.get('/auth/license/machineID')
    this.machineId = res.data?.id
    return res
  }

  @action
  async getLicenseInfo() {
    const res = await Http.get('/auth/license')
    this.info = new Info(res.data)
    return res
  }

  updateLicense = (licenseAttribute: {}) => {
    return Http.post('/auth/license', {
      ...licenseAttribute,
      machine_id: this.machineId
    })
  }
}
export default new License()
