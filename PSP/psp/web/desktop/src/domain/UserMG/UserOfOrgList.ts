import { Http } from '@/utils'
import { action, observable } from 'mobx'
import UserOfOrg from './UserOfOrg'

export class Organization {
  @observable userList: UserOfOrg[]

  syncOrganizationStructure = () => {
    return Http.put('/organize/sync')
  }
  getOrganization = async () => {
    const res = await Http.get('/organize')
    return res
  }

  @action
  getUserList(params) {
    return Http.get('/organize/userList', {
      params: {
        id: params.id,
        index: params.index,
        size: params.size,
        orderAsc: params.orderAsc,
        orderBy: params.orderBy,
        name: params.name,
        created_at: params.created_at,
        search_value: params.search_value,
      },
    }).then(res => {
      this.userList = res?.data.list.map(user => new UserOfOrg(user))
      return res
    })
  }

  @action
  update(id, body) {
    return Http.put(`/organize/userList/${id}`, body)
  }

  @action
  public active = (id, name, body) => {
    return Http.post(`/organize/active/${id}?username=${name}`, {
      ...body,
    })
  }

  @action
  public inactive = (id, name, body) => {
    return Http.post(`/organize/inactive/${id}?username=${name}`, {
      ...body,
    })
  }
}
export default new Organization()
