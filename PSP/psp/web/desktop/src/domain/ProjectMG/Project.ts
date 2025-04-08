import { observable } from 'mobx'
import moment from 'moment'
import { currentUser } from '..'
import { Http } from '@/utils'

interface IProps {
  id: string
  project_name: string
  state: string
  create_time: string
  start_time: string
  end_time: string

  is_project_owner: boolean
  project_owner_id: string
  project_owner_name: string

  comment: string

  members: any[]

}

export default class Project {
  @observable id: string
  @observable project_name: string
  @observable state: string
  @observable create_time: string

  @observable start_time: string
  @observable end_time: string

  @observable is_project_owner: boolean
  @observable project_owner_id: string
  @observable project_owner_name: string

  @observable comment: string

  @observable members: any[] = []

  @observable start_time_momnet: any
  @observable end_time_momnet: any

  get isOwner() {
    return this.project_owner_id === currentUser.id
  }

  async getDetail() {
    const res = await Http.get(`/project/detail`, {
      params: {
        project_id: this.id
      }
    })

    this.members = res.data.members

    return res
  }

  constructor({
    create_time,
    start_time,
    end_time,
    ...props
  }: Partial<IProps>) {
    Object.assign(this, props)

    if (start_time) {
      this.start_time_momnet = moment(start_time)
      this.start_time = moment(start_time).format('YYYY-MM-DD HH:mm:ss')
    }

    if (end_time) {
      this.end_time_momnet = moment(end_time)
      this.end_time = moment(end_time).format('YYYY-MM-DD HH:mm:ss')
    }

    if (create_time) {
      this.create_time = moment(create_time).format('YYYY-MM-DD HH:mm:ss')
    }
  }
}
