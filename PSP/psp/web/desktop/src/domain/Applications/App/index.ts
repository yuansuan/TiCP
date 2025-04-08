/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, computed, observable, runInAction, toJS } from 'mobx'
import nanoid from 'nanoid'

import { arrayDifference, arrayDifferenceByKey, difference } from './utils'
import { Http } from '@/utils'
import AppIcons from '../AppIcons'
import HelpDoc, { IRequest as IHelpDocRequest } from './HelpDoc'
import SubForm, { IRequest as ISubFormRequest } from './SubForm'

export interface IRequest {
  sub_form: ISubFormRequest
  description?: string
  help_doc: IHelpDocRequest
  state: string
  last_modified_time: string
  last_modifier: string
  name: string
  id: string
  image?: string
  license_manager_id?: string
  bin_path?: Array<{ key: string; value: string }>
  scheduler_param?: Array<{ key: string; value: string }>
  out_app_id?: string
  cloud_out_app_id?: string
  cloud_out_app_name?: string
  residual_log_parser: string
  enable_residual: boolean
  enable_snapshot: boolean
  script: string
  icon?: string
  type?: string
  compute_type?: string
  isLiked?: boolean
  version?: string // app version
  queues?: Array<{ queue_name: string; cpu_number: number; select: boolean }>
  licenses?: Array<{ id: string; name: string; select: boolean;licence_valid: boolean }>
}

type ReqFields = Array<{
  id: string
  type: string
  value: string
  values: string[]
  master_slave: string
  required: boolean
}>

interface ITestRequest {
  name: string
  req_fields: ReqFields
  scheduler: string
  upload_sub_token: string
  script: string
  state: string
}

interface ISubmitRequest {
  name: string
  req_fields: ReqFields
  upload_sub_token: string
  resubmit?: boolean
  job_id?: any
}

interface IApp {
  appId: string
  _key: string
  subForm: SubForm
  description: string
  helpDoc: HelpDoc
  state: string
  type: string
  computeType: string
  version: string
  lastModifiedTime: string
  lastModifier: string
  name: string
  script: string
  iconData: string
  scriptData: string
  image?: string
  license_manager_id?: string
  cloud_out_app_id?: string
  cloud_out_app_name?: string
  residual_log_parser: string
  enable_residual:boolean
  enable_snapshot:boolean
  bin_path?: Array<{ key: string; value: string }>
  scheduler_param?: Array<{ key: string; value: string }>
  queues?: Array<{ queue_name: string; cpu_number: number; select: boolean }>
  licenses?: Array<{ id: string; name: string; select: boolean;licence_valid: boolean  }>
}

export default class App implements IApp {
  @computed
  get isPublished() {
    return this.state === 'published'
  }

  @computed
  get fieldIds() {
    return this.subForm.sections.reduce((ids, section) => {
      ids = [
        ...ids,
        ...section.fields
          .filter(item => item.id !== undefined)
          .map(item => item.id)
      ]

      return ids
    }, [])
  }

  public static fetch = ({ name, state, version, compute_type }) =>
    Http.get(`/app/template`, {
      params: { name, version, state, compute_type }
    }).then(res => new App(res.data.app))

  public static fetchYSCloud = ({ app_id }) =>
    Http.get(`/application/yscloud`, {
      params: { app_id }
    }).then(res => new App(res.data.app))

  public static resubmit = jobId =>
    Http.get('/application/resubmit', {
      params: { job_id: jobId }
    }).then(res => new App(res.data.app))

  public initialValue
  @observable public _key = nanoid()
  @observable public subForm
  @observable public description = ''
  @observable public helpDoc = null
  @observable public state = ''
  @observable public lastModifiedTime = ''
  @observable public lastModifier = ''
  @observable public name = ''
  @observable public icon = ''
  @observable public type = ''
  @observable public computeType = ''
  @observable public appId = ''
  @observable public script = ''
  @observable public iconData = ''
  @observable public scriptData = ''
  @observable public isCloud = false
  @observable public isLiked = false
  @observable public version = ''
  @observable public image = ''
  @observable public bin_path = []
  @observable public scheduler_param = []
  @observable public queues = []
  @observable public licenses = []
  @observable public license_manager_id = ''
  @observable public cloud_out_app_id = ''
  @observable public cloud_out_app_name = ''
  @observable public enable_residual = false
  @observable public enable_snapshot = false
  @observable public residual_log_parser = ''
  

  constructor(request?: IRequest) {
    this.init(request)
  }

  @action
  public init = (request?: IRequest) => {
    // update _key trigger react render
    this._key = nanoid()
    request &&
      Object.assign(this, {
        subForm: new SubForm(request.sub_form),
        description: request.description,
        helpDoc: new HelpDoc(request.help_doc),
        state: request.state,
        lastModifiedTime: request.last_modified_time,
        lastModifier: request.last_modifier,
        name: request.name,
        appId: request.id,
        script: request.script,
        iconData: request.icon,
        isLiked: request.isLiked || false,
        version: request.version,
        type: request.type,
        computeType: request.compute_type,
        bin_path: request.bin_path,
        scheduler_param: request.scheduler_param,
        image: request.image,
        queues: request.queues,
        licenses: request.licenses,
        license_manager_id: request.license_manager_id,
        cloud_out_app_id: request.cloud_out_app_id,
        cloud_out_app_name: request.cloud_out_app_name,
        enable_residual: request.enable_residual,
        enable_snapshot: request.enable_snapshot,
        residual_log_parser: request.residual_log_parser,
      })
    this.initialValue = this.snapshot()
  }

  public snapshot = () => {
    const copy: any = this.toRequest()
    // hack: get icon data URIs
    copy.icon_data = this.iconData

    return toJS(copy)
  }

  // get the diff info after the latest common
  public diff = () => {
    const newValue = this.snapshot()
    const appDiff = difference(newValue, this.initialValue)

    if (appDiff && appDiff.application) {
      const appDiffSubForm = appDiff.application.sub_form
      if (appDiffSubForm && appDiffSubForm.section) {
        const sectionDiff = arrayDifference(
          newValue.application.sub_form.section,
          this.initialValue.application.sub_form.section,
          'name'
        )

        // elaborate the section update
        sectionDiff.update = sectionDiff.update
          .map(item => {
            const diff = difference(item.new, item.old)

            if (!diff) {
              return null
            }

            // elaborate the field update
            if (diff && diff.field) {
              diff.field = arrayDifference(item.new.field, item.old.field, 'id')
              diff.field.update = diff.field.update
                .map(fieldItem => {
                  const diff = difference(fieldItem.new, fieldItem.old)
                  if (!diff) {
                    return null
                  }

                  diff.options = arrayDifference(
                    fieldItem.new.options,
                    fieldItem.old.options
                  )
                  diff.default_values = arrayDifference(
                    fieldItem.new.default_values,
                    fieldItem.old.default_values
                  )

                  return {
                    props: diff,
                    key: fieldItem.old.id
                  }
                })
                .filter(fieldItem => !!fieldItem)
            }

            return { ...diff, key: item.old.name }
          })
          .filter(item => !!item)

        appDiffSubForm.section = sectionDiff
      }

      if (appDiff.application?.bin_path) {
        const diffBinPath = arrayDifferenceByKey(
          newValue.application.bin_path,
          this.initialValue.application.bin_path,
          'key'
        )
        appDiff.application.bin_path = diffBinPath
      }

      if (appDiff.application?.scheduler_param) {
        const diffSchedulerParam = arrayDifferenceByKey(
          newValue.application.scheduler_param,
          this.initialValue.application.scheduler_param,
          'key'
        )
        appDiff.application.scheduler_param = diffSchedulerParam
      }

      if (appDiff.application?.queues) {
        const diffQueues = arrayDifferenceByKey(
          newValue.application.queues
            .filter(item => item.select)
            ?.map(item => item.queue_name),
          this.initialValue.application.queues
            .filter(item => item.select)
            ?.map(item => item.queue_name),
          'name'
        )
        appDiff.application.queues = diffQueues
      }
      if (appDiff.application?.licenses) {
        const diffLicense = arrayDifferenceByKey(
          newValue.application.licenses
            .filter(item => item.select)
            ?.map(item => item.name),
          this.initialValue.application.licenses
            .filter(item => item.select)
            ?.map(item => item.name),
          'name'
        )
        appDiff.application.licenses = diffLicense
      }
    }

    return appDiff
  }

  @action
  public setScriptData = scriptData => (this.scriptData = scriptData)

  public reset = () => this.init(this.initialValue.application)

  public fetch = () =>
    Http.get(`/app/template`, {
      params: { name: this.name, compute_type: this.computeType }
    }).then(res => this.init(res.data.app))

  public fetchScript = (computeType?: string) =>
    Http.get('/app/template', {
      params: {
        name: this.name,
        compute_type: this.computeType || computeType
      }
    }).then(res =>
      runInAction(() => {
        const scriptData = res.data.app.script
        this.setScriptData(scriptData)
        this.initialValue.script_data = scriptData
      })
    )

  public publish = () =>
    Http.put('/app/template/publish', {
      names: [this.name],
      state: 'published',
      compute_type: this.computeType
    }).then(res => {
      this.state = 'published'
      this.fetch()

      return res
    })

  public unpublish = () =>
    Http.put('/app/template/publish', {
      names: [this.name],
      state: 'unpublished',
      compute_type: this.computeType
    }).then(res => {
      this.state = 'unpublished'
      this.fetch()

      return res
    })

  public save = () => {
    const params = this.toRequest()

    return Http.put('/app/template', {
      app: params.application,
      state: this.state
    }).then(res => {
      // update AppIcons cache
      AppIcons.list.set(this.name, this.iconData)

      return res
    })
  }

  public saveAs = newVersion => {
    const request = this.toRequest()
    const copyParams = request.application
    copyParams.version = newVersion
    delete copyParams.name

    return Http.put(
      '/app/template',
      {
        app: copyParams,
        base_name: this.name,
        base_state: this.state,
        compute_type: this.computeType
      },
    )
  }

  public test = (params: ITestRequest) => Http.post('/application/test', params)

  public static submit = (params: ISubmitRequest) =>
    Http.post('/application/submit', params)

  public static submitYSCloud = (id, params: ISubmitRequest) =>
    Http.post(`/application/submit-yscloud/${id}`, params)

  public toRequest = (): {
    application: IRequest
    icon_data: string
    script_data: string
  } => ({
    application: {
      sub_form: this.subForm && this.subForm.toRequest(),
      description: this.description,
      help_doc: this.helpDoc && this.helpDoc.toRequest(),
      state: this.state,
      last_modified_time: this.lastModifiedTime,
      last_modifier: this.lastModifier,
      name: this.name,
      type: this.type,
      compute_type: this.computeType,
      id: this.appId,
      version: this.version,
      image: this.image,
      bin_path: this.bin_path,
      scheduler_param: this.scheduler_param,
      queues: this.queues,
      licenses: this.licenses,
      cloud_out_app_id: this.cloud_out_app_id,
      cloud_out_app_name: this.cloud_out_app_name,
      enable_residual: this.enable_residual,
      enable_snapshot: this.enable_snapshot,
      residual_log_parser: this.residual_log_parser,
      script: this.scriptData || this.script,
      icon: this.iconData || this.icon
    },
    icon_data: this.iconData,
    script_data: this.scriptData
  })
}
