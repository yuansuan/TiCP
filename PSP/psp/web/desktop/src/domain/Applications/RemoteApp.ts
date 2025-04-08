import { action, computed, observable, toJS } from 'mobx'
import nanoid from 'nanoid'
import { Http } from '@/utils'
import RemoteAppIcons from './RemoteAppIcons'
import HelpDoc, { IRequest as IHelpDocRequest } from './App/HelpDoc'
import SubForm, { IRequest as ISubFormRequest } from './App/SubForm'

export interface IRequest {
  sub_form: ISubFormRequest
  description?: string
  help_doc: IHelpDocRequest
  state: string
  last_modified_time: string
  last_modifier: string
  name: string
  id: string
  out_app_id: string
  script: string
  icon?: string
  type?: string
  compute_type?: string
  isLiked?: boolean
  version?: string
  residual_log_parser?: string
  enable_residual?: boolean
  enable_snapshot?: boolean
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
  outAppId: string
  script: string
  iconData: string
  residualLogParser: string
  enableResidual: boolean
  enableSnapshot: boolean
  scriptData: string
}

export default class RemoteApp implements IApp {
  @computed
  get isPublished() {
    return this.state === 'published'
  }

  public static fetch = ({ name, state }) =>
    Http.get(`/remote_app`, {
      params: { name, state, type: 'remote' }
    }).then(res => new RemoteApp(res.data.application))

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
  @observable public outAppId = ''
  @observable public script = ''
  @observable public iconData = ''
  @observable public scriptData = ''
  @observable public isCloud = false
  @observable public isLiked = false
  @observable public version = ''
  @observable public residualLogParser = ''
  @observable public enableResidual = false
  @observable public enableSnapshot = false

  constructor(request?: IRequest) {
    this.init(request)
  }

  @action
  public init = (request?: IRequest) => {
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
        outAppId: request.out_app_id,
        script: request.script,
        iconData: request.icon,
        isLiked: request.isLiked || false,
        version: request.version,
        type: request.type,
        computeType: request.compute_type,
        residualLogParser: request.residual_log_parser,
        enableResidual: request.enable_residual,
        enableSnapshot: request.enable_snapshot,
      })
  }
  public updateIcon = icon => {
    this.icon = icon
    RemoteAppIcons.list.set(this.name, icon)
  }
  public snapshot = () => {
    const copy: any = this.toRequest()
    // hack: get icon data URIs
    copy.icon_data = this.iconData

    return toJS(copy)
  }

  public fetch = () =>
    Http.get(`/app/template`, {
      params: { name: this.name, compute_type: 'cloud' }
    }).then(res => this.init(res.data.app))

  public publish = () =>
    Http.put('/app/template/publish', {
      names: [this.name],
      state: 'published',
      compute_type: 'cloud'
    }).then(res => {
      this.state = 'published'
      this.fetch()

      return res
    })

  public unpublish = () =>
    Http.put('/app/template/publish', {
      names: [this.name],
      state: 'unpublished',
      compute_type: 'cloud'
    }).then(res => {
      this.state = 'unpublished'
      this.fetch()

      return res
    })

  public save = () => {
    // update AppIcons cache
    if (this.icon) {
      RemoteAppIcons.list.set(this.icon, this.iconData)
    }

    return Http.put('/remote_app/save', this.toRequest())
  }

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
      out_app_id: this.outAppId,
      version: this.version,
      script: this.scriptData || this.script,
      icon: this.iconData || this.icon,
      residual_log_parser: this.residualLogParser,
      enable_residual: this.enableResidual,
      enable_snapshot: this.enableSnapshot
    },
    icon_data: this.iconData,
    script_data: this.scriptData
  })
}
