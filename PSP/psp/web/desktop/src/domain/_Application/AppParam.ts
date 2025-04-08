/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable } from 'mobx'

export enum FieldType {
  radio = 'radio',
  checkbox = 'checkbox',
  text = 'text',
  list = 'list',
  multiple = 'multiple',
  label = 'label',
  date = 'date',
  lsfile = 'lsfile',
  texarea = 'textarea',
  node_selector = 'node_selector',
  cascade_selector = 'cascade_selector'
}

export interface IAppParam {
  id: string
  action: string
  default_value: string
  default_values: string[]
  hidden: boolean
  label: string
  help: string
  options: string[]
  options_from: string
  options_script: string
  post_text: string
  required: boolean
  type: FieldType
  value: string
  values: string[]
  custom_json_value_string: string
}

export default class AppParam {
  @observable id: string
  @observable action: string
  @observable defaultValue: string
  @observable defaultValues: string[]
  @observable hidden: boolean
  @observable label: string
  @observable help: string
  @observable options: string[]
  @observable optionsFrom: string
  @observable optionsScript: string
  @observable postText: string
  @observable required: boolean
  @observable type: FieldType
  @observable value: string
  @observable values: string[]
  @observable customJSONValueString: string

  constructor(data: IAppParam) {
    Object.assign(this, {
      id: data.id,
      action: data.action || '',
      defaultValue: data.default_value || '',
      defaultValues: data.default_values || [],
      hidden: data.hidden || false,
      label: data.label || '',
      help: data.help || '',
      options: data.options || [],
      optionsFrom: data.options_from || 'custom',
      optionsScript: data.options_script || '',
      postText: data.post_text || '',
      required: data.required || false,
      type: data.type,
      value: data.value || '',
      values: data.values || [],
      customJSONValueString: data.custom_json_value_string || '{}',
    })
  }
}
