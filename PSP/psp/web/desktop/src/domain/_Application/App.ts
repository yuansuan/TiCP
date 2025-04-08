/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable } from 'mobx'
import { Http } from '@/utils'
import { ISection } from './Section'
import AppParam from './AppParam'

export class AppProps {
  @observable id: string
  @observable out_app_id: string
  @observable name: string
  @observable icon: string
  @observable compute_type: string
  @observable type: string
  @observable script: string
  @observable version: string
  @observable state: string
  @observable sub_form: AppParam[]
  @observable description: string
  @observable help_doc: string
  @observable image_id: string
  @observable image?: string
  @observable bin_path?: []
  @observable scheduler_param?: []
  @observable queues?: []
  @observable licenses?: []
  @observable cloud_out_app_id?: string
  @observable cloud_out_app_name?: string
}

export default class App extends AppProps {
  constructor(data: AppProps) {
    super()
    Object.assign(this, data)
  }

  async getParams() {
    let params = []

    const transformedData = this.sub_form?.section?.map(section => {
      const newSection = section?.field?.filter(item => !item?.hidden).map(item => new AppParam(item))

      return {
        field: newSection,
        name: section.name
      }
    })
    return transformedData || []
  }
}
