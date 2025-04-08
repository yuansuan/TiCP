/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import nanoid from 'nanoid'

export enum FieldType {
  radio = 'radio',
  checkbox = 'checkbox',
  text = 'text',
  list = 'list',
  multiple = 'multiple',
  label = 'label',
  date = 'date',
  // lsfile = 'lsfile',
  texarea = 'textarea',
  // lsfile_yscloud = 'lsfile_yscloud',
  node_selector = 'node_selector',
  cascade_selector = 'cascade_selector',
}

export interface IRequest {
  id: string
  action: string
  default_value: string
  default_values: string[]
  hidden: boolean
  label: string
  help: string
  options_from: 'custom' | 'script' | string
  options_script: string
  options: string[]
  post_text: string
  required: boolean
  type: FieldType
  value: string
  values: string[]
  file_from_type: string
  is_master_slave: boolean
  master_slave: string
  master_include_extensions: string
  master_include_keywords: string
  // 定制化组件，需要额外存储保存的数据
  custom_json_value_string: string

  is_support_master: boolean
  master_file: string

  is_support_workdir: boolean
  workdir: string
}

interface IField {
  _key: string
  id: string
  action: string
  defaultValue: string
  defaultValues: string[]
  hidden: boolean
  label: string
  help: string
  options: string[]
  optionsFrom: 'custom' | 'script'
  optionsScript: string
  postText: string
  required: boolean
  type: FieldType
  value: string
  values: string[]
  fileFromType: string
  isMasterSlave: boolean
  masterIncludeExtensions: string
  masterIncludeKeywords: string
  masterSlave: string
  editing: boolean
  customJSONValueString: string

  // 支持主文件
  isSupportMaster: boolean
  masterFile: string

  // 指定工作目录
  isSupportWorkdir: boolean
  workdir: string
}

export default class Field implements IField {
  @observable _key = nanoid()
  @observable id
  @observable action = ''
  @observable defaultValue = ''
  @observable defaultValues = []
  @observable hidden = false
  @observable label = ''
  @observable help = ''
  @observable options = []
  @observable optionsFrom: 'custom' | 'script' = 'custom'
  @observable optionsScript = ''
  @observable postText = ''
  @observable required = false
  @observable type
  @observable value = ''
  @observable values = []
  @observable fileFromType = ''
  @observable masterSlave = ''
  @observable isMasterSlave = false
  @observable masterIncludeExtensions = ''
  @observable masterIncludeKeywords = ''
  @observable customJSONValueString = '{}'

  @observable isSupportMaster = true // 默认支持主文件
  @observable masterFile = ''

  @observable isSupportWorkdir = false
  @observable workdir = ''

  @observable editing = false

  constructor(request?: Partial<IRequest>) {
    request && this.init(request)

    // hack: new field enter editing mode by default
    if (!request || !request.id) {
      this.updateEditing(true)
    }
  }

  @action
  init = (request: Partial<IRequest>) => {
    Object.assign(this, {
      id: request.id,
      action: request.action || '',
      defaultValue: request.default_value || '',
      defaultValues: request.default_values || [],
      hidden: request.hidden || false,
      label: request.label || '',
      help: request.help || '',
      options: request.options || [],
      optionsFrom: request.options_from || 'custom',
      optionsScript: request.options_script || '',
      postText: request.post_text || '',
      required: request.required || false,
      type: request.type,
      value: request.value || '',
      values: request.values || [],
      fileFromType: request.file_from_type || '',
      isMasterSlave: request.is_master_slave || false,
      masterSlave: request.master_slave || '',
      masterIncludeExtensions: request.master_include_extensions || '',
      masterIncludeKeywords: request.master_include_keywords || '',
      customJSONValueString: request.custom_json_value_string || '{}',
      isSupportMaster: true, // 只支持一种方式: 新主文件模式
      masterFile: request.master_file,
      isSupportWorkdir: request.is_support_workdir,
      workdir: request.workdir
    })
  }

  @action
  updateEditing = editing => (this.editing = editing)

  toRequest = (): IRequest => ({
    id: this.id,
    action: this.action,
    default_value: this.defaultValue,
    default_values: [...this.defaultValues],
    hidden: this.hidden,
    label: this.label,
    help: this.help,
    options: this.type === 'cascade_selector' ? [] : [...this.options],
    options_from: this.optionsFrom,
    options_script: this.optionsScript,
    post_text: this.postText,
    required: this.required,
    type: this.type,
    value: this.value,
    values: [...this.values],
    file_from_type: this.fileFromType,
    is_master_slave: this.isMasterSlave,
    master_slave: this.masterSlave,
    master_include_extensions: this.masterIncludeExtensions,
    master_include_keywords: this.masterIncludeKeywords,
    custom_json_value_string: this.customJSONValueString,
    is_support_master: this.isSupportMaster,
    master_file: this.masterFile,
    is_support_workdir: this.isSupportWorkdir,
    workdir: this.workdir
  })
}
