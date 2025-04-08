/* Copyright (C) 2016-present, Yuansuan.cn */
import { observable, runInAction } from 'mobx'
import { Http } from '@/utils'
import { IncomingPageAware, OutcomingPageAware } from '../supports/page'
import { HardwareList } from './hardware'
import qs from 'qs'
import { message } from 'antd'
import { Timestamp } from '@/domain/common'
import currentUser from '@/domain/User'
import moment from 'moment'
// <rootId, isReady>
const CACHEMAP: Map<string, boolean> = new Map()

export const HARDWARE_PLATFORM_ALL = 0
export const HARDWARE_PLATFORM_LINUX = 1
export const HARDWARE_PLATFORM_WINDOWS = 2

export const HARDWARE_PLATFORM_MAP = {
  0: '-',
  1: 'Linux',
  2: 'Windows'
}

export const HARDWARE_DISPLAY_UNKNOWN = 0
export const HARDWARE_DISPLAY_DESKTOP = 1
export const HARDWARE_DISPLAY_APPLICATION = 2

export const HARDWARE_DISPLAY_MAP = {
  0: '-',
  1: '桌面模式',
  2: '应用模式'
}

export const INSTANCE_STATUS_PENDING = 0
export const INSTANCE_STATUS_CREATED = 1
export const INSTANCE_STATUS_RUNNING = 2
export const INSTANCE_STATUSTERMINATED = 3

export const INSTANCE_STATUS_MAP = {
  PENDING: '创建中',
  CREATED: '已创建',
  RUNNING: '运行中',
  TERMINATED: '已结束'
}

export const CHARGE_TYPE_MAP = {
  1: '按量付费',
  2: '竞价付费'
}
export const SESSION_STATUS_ALL = 0
export const SESSION_STATUS_PENDING = 1
export const SESSION_STATUS_STARTING = 2
export const SESSION_STATUS_STARTED = 3
export const SESSION_STATUS_CLOSING = 4
export const SESSION_STATUS_CLOSED = 5
export const SESSION_STATUS_POWERING_OFF = 6
export const SESSION_STATUS_POWER_OFF = 7
export const SESSION_STATUS_POWERING_ON = 8
export const SESSION_STATUS_REBOOTING = 9

export const SESSION_STATUS_MAP = {
  PENDING: '等待资源',
  STARTING: '启动中',
  STARTED: '已启动',
  CLOSING: '删除中',
  CLOSED: '已删除',
  UNAVAILABLE: '不可用',
  REBOOTING: '重启中',
  'POWERING OFF': '关机中',
  'POWER OFF': '已关机',
  'POWERING ON': '开机中'
}

export const statusMapping = {
  等待资源: 'primary',
  启动中: 'primary',
  已启动: 'success',
  出错: 'error',
  不可用: 'error',
  删除中: 'warn',
  已删除: 'canceled',
  重启中: 'primary',
  关机中: 'primary',
  已关机: 'canceled',
  开机中: 'primary'
}

export const SESSION_STATUS_BUTTON_LOADING = {
  READYING: '就绪中',
  READIED: '已就绪'
}

export class Project {
  @observable id: string
  @observable name: string

  constructor(props: Partial<Project>) {
    Object.assign(this, props)
  }
}

export class Hardware {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable network: number
  @observable cpu: number
  @observable mem: number
  @observable gpu: number
  @observable cpu_model: number
  @observable gpu_model: number
  @observable instance_family: string
  @observable instance_type: string
  @observable default_preset: boolean

  constructor(props: Partial<Hardware>) {
    Object.assign(this, props)
  }
}

export class Software {
  @observable id: string
  @observable name: string
  @observable desc: string
  @observable platform: number
  @observable display: number
  @observable icon: string
  @observable image_id: string
  @observable presets: []
  @observable init_script: string
  @observable update_time: Timestamp
  @observable create_time: Timestamp

  @observable remote_apps: []
  @observable gpu_desired: boolean

  constructor(props: Partial<Software>) {
    Object.assign(this, props)
  }
}

export class Session {
  @observable id: string
  @observable out_app_id: string //实际实例id
  @observable status: string
  @observable stream_url: string
  @observable start_time: string = '--'
  @observable update_time: string = '--'
  @observable end_time: string = '--'
  @observable exit_reason: string
  @observable create_time: string = '--'
  @observable hardware: Hardware
  @observable software: Software
  @observable project_name: string = '--'

  constructor(props: Partial<Session>) {
    Object.assign(this, props)
    if (props.start_time) {
      this.start_time = moment(props.start_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (props.update_time) {
      this.update_time = moment(props.update_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (props.create_time) {
      this.create_time = moment(props.create_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (props.end_time) {
      this.end_time = moment(props.end_time).format('YYYY-MM-DD HH:mm:ss')
    }
  }
}

export class SessionUser {
  @observable real_name: string
  @observable phone: string
  @observable username: string
  @observable display_name: string

  constructor(props: Partial<SessionUser>) {
    Object.assign(this, props)
  }
}

export class ListHardwareRequest extends OutcomingPageAware {
  @observable number_of_cpu?: number = 0
  @observable number_of_mem?: number = 0
  @observable number_of_gpu?: number = 0

  constructor(props?: Partial<ListHardwareRequest>) {
    super(props)
    Object.assign(this, props)
  }
}

class ListHardwareResponse extends IncomingPageAware {
  @observable hardwares: Array<Hardware> = []

  constructor(props: any) {
    super(props)

    props?.hardwares?.map((item: any) => {
      this.hardwares.push(new Hardware(item))
    })
  }
}

class ListSoftwareRequest extends OutcomingPageAware {
  @observable id: string
  @observable remote_apps: []
  @observable desc: string
  @observable image_id: string
  @observable init_script: string
  @observable icon: string
  @observable gpu_desired: boolean
  @observable name: string
  @observable state: string
  @observable platform: number
  @observable display: number
  @observable create_time: Timestamp
  @observable update_time: Timestamp

  constructor(props?: Partial<ListSoftwareRequest>) {
    super(props)
    Object.assign(this, props)
  }
}

export class ListSoftwareResponse extends IncomingPageAware {
  @observable softwares: Array<Software> = []

  constructor(props: any) {
    super(props)

    props?.softwares?.map((item: any) => {
      this.softwares.push(new Software(item))
    })
  }
}

export class StartSessionResponse {
  @observable session_id: string
  @observable webrtc_url: string
  @observable room_id: string
  @observable stream_url: string

  constructor(props: any) {
    Object.assign(this, props)
  }
}

export class ListSessionRequest extends OutcomingPageAware {
  @observable statuses: string[] = []
  @observable hardware_ids: string[] = []
  @observable software_ids: string[] = []
  @observable project_ids: string[] = []
  @observable user_name: string = ''
  @observable page_index: number = 1
  @observable page_size: number = 10

  constructor(props?: Partial<ListSessionRequest>) {
    super(props)
    Object.assign(this, props)
  }
}
export class pollRequestItem {
  @observable loading: boolean = false
  @observable timer: any = null
  @observable status: string
  @observable realStatus: string
  @observable token: string
  @observable refreshSession: boolean = false
  @observable timerList: Array<any> = []
}
export class ListSessionItem extends pollRequestItem {
  @observable session: Session
  constructor(props: any) {
    super()
    this.session = new Session(props)
    this.refresh(props)
  }
  clearSessionTimers() {
    this.timerList.forEach(item => {
      if (item.timer !== null) {
        clearInterval(item.timer)
      }
    })
    this.timerList = []
  }
  refresh(session) {
    this.status = session.status
    this.roomId = session?.id
    const isReady = CACHEMAP.get(this.roomId)

    // 从之前缓存里面拿，room的状态，如果已经 ready，没必要再去请求
    if (isReady) {
      this.loading = false
      if (session.status === 'STARTED') {
        this.realStatus = 'READIED'
      } else if (session.status === 'CLOSING' || session.status === 'CLOSED') {
        this.realStatus = ''
        CACHEMAP.delete(this.roomId)
      }
      return
    }

    // 修复打开会话 button loading 导致的闪烁问题
    this.loading =
      session.status === 'STARTING' ||
      session.status === 'STARTED' ||
      session.status === 'POWERING OFF' ||
      session.status === 'POWERING ON' ||
      session.status === 'REBOOTING'

    if (session.status === 'STARTED') {
      this.realStatus = 'READYING'
      if (this.realStatus === 'READIED') return
      runInAction(async () => {
        try {
          await new Vis().pollFetchRequest(session.id).then(res => {
            this.loading = !res.data.ready
            if (res.data.ready) {
              this['realStatus'] = 'READIED'
            }
            CACHEMAP.set(this.roomId, res.data.ready)
          })
        } catch (error) {
          // this['realStatus'] = 'READIED'
        }
      })
    }
  }
}

export class ListSessionResponse extends IncomingPageAware {
  @observable sessions: Array<ListSessionItem> = []

  constructor(props?: any) {
    super(props)

    props?.sessions.map((item: any) => {
      this.sessions.push(new ListSessionItem(item))
    })
  }
}
export class RelatedHardWareResponse extends IncomingPageAware {
  @observable hardwares: Array<ListSessionItem> = []

  constructor(props?: any) {
    super(props)

    props?.data.map((item: any) => {
      this.hardwares.push(item)
    })
  }
}

const acquireAllItems = { page_size: 1000, page_index: 1 }

class VisBase {}

export default class Vis extends VisBase {
  @observable filterParams: {
    statuses: []
    hardware_ids: []
    software_ids: []
    project_ids: []
    user_name: ''
    page_index: number
    page_size: number
  }

  constructor(props?: Partial<VisBase>) {
    super()
  }

  async getProjectName() {
    const res = await Http.get('/vis/session/projectNames', {
      params: {
        has_used: !currentUser.roles.some(r => [1, 2].includes(r.id)) // 1,2 系统管理员、超级管理员
      }
    })
    return res.data?.names
  }

  async getProjects(is_admin) {
    return Http.get('/project/listForParam', {
      params: {
        is_admin: is_admin ?? currentUser.hasSysMgrPerm
      }
    })
  }

  async getCurrentUserProjects(): Promise<Array<Project>> {
    const { data } = await Http.get('/project/list/current', {
      params: {
        state: 'Running'
      }
    })
    return data?.projects || []
  }

  async listHardwareFilter(
    request: ListHardwareRequest
  ): Promise<ListHardwareResponse> {
    return new ListHardwareResponse(
      (
        await Http.get('/vis/hardware', {
          params: {
            ...request,
            has_used: true,
            is_admin: false
          }
        })
      ).data
    )
  }
  async listHardware(
    request: ListHardwareRequest
  ): Promise<ListHardwareResponse> {
    return new ListHardwareResponse(
      (
        await Http.get('/vis/hardware', {
          params: {
            ...request,
            is_admin: false
          }
        })
      ).data
    )
  }

  async listAllHardware(): Promise<Array<Hardware>> {
    return (
      await this.listHardwareFilter(
        new ListHardwareRequest({ ...acquireAllItems })
      )
    ).hardwares
  }
  async listPermHardware(): Promise<Array<Hardware>> {
    return (
      await this.listHardware(new ListHardwareRequest({ ...acquireAllItems }))
    ).hardwares
  }
  get filterQuery() {
    return (
      this.filterParams || {
        statuses: [],
        hardware_ids: [],
        software_ids: [],
        project_ids: [],
        user_name: '',
        page_index: 1,
        page_size: 20
      }
    )
  }

  setFilterParams(filterParams) {
    this.filterParams = filterParams
  }

  async listSoftwareFilter(
    request: Partial<ListSoftwareRequest>
  ): Promise<ListSoftwareResponse> {
    return new ListSoftwareResponse(
      (
        await Http.get('/vis/software', {
          params: {
            name: '',
            has_used: true, // has_permission 和has_used 不能同时用
            ...request,
            is_admin: false
          }
        })
      ).data
    )
  }

  async listSoftware(
    request: Partial<ListSoftwareRequest>
  ): Promise<ListSoftwareResponse> {
    return new ListSoftwareResponse(
      (
        await Http.get('/vis/software', {
          params: {
            name: '',
            has_permission: true, //true代表需要查询权限
            ...request,
            state: 'published',
            is_admin: false
          }
        })
      ).data
    )
  }

  async getUsingSoftware() {
    return await Http.get('/vis/software/usingStatuses')
  }
  async listAllSoftware(): Promise<Array<Software>> {
    return (
      await this.listSoftwareFilter(
        new ListSoftwareRequest({ ...acquireAllItems })
      )
    ).softwares
  }

  async listPermSoftware(): Promise<Array<Software>> {
    return (
      await this.listSoftware(new ListSoftwareRequest({ ...acquireAllItems }))
    ).softwares
  }
  async startSession(props: {
    hardware_id: string
    software_id: string
    project_id: string
    project_name: string
  }): Promise<StartSessionResponse> {
    return new StartSessionResponse(
      (
        await Http.post(
          '/vis/session',
          {
            ...props,
            user_name: currentUser.name,
            user_id: currentUser.id
          },
          {}
        )
      ).data
    )
  }

  async closeSession(session_id: string, exit_reason?: string): Promise<any> {
    CACHEMAP.delete(session_id)
    return (
      await Http.post(
        '/vis/session/close',
        {
          session_id,
          exit_reason
        },
        {}
      )
    ).data
  }

  async restartSession(
    session_id: string,
    admin = false,
    reason = ''
  ): Promise<any> {
    CACHEMAP.delete(session_id)
    return Http.post('/vis/session/reboot', {
      session_id,
      reason,
      admin
    })
  }

  async powerOnSession(session_id: string): Promise<any> {
    CACHEMAP.delete(session_id)
    return Http.post('/vis/session/powerOn', {
      session_id
    })
  }

  async powerOffSession(session_id: string): Promise<any> {
    CACHEMAP.delete(session_id)
    return Http.post('/vis/session/powerOff', {
      session_id
    })
  }
  async deleteSession(session_id: string): Promise<any> {
    CACHEMAP.delete(session_id)
    return (
      await Http.delete('/vis/session', {
        data: { session_id }
      })
    ).data
  }
  async updateSession(
    session_id: string,
    autoClose: boolean,
    time?: string
  ): Promise<any> {
    CACHEMAP.delete(session_id)
    return await Http.put(`/vis/session/${session_id}`, {
      is_auto_close: autoClose,
      auto_close_time: time
    })
  }

  // 获取所有会话列表
  async listSession(
    request: Partial<ListSessionRequest>
  ): Promise<ListSessionResponse> {
    const newReq = Object.assign(this.filterQuery, request)
    return new ListSessionResponse(
      (
        await Http.get('/vis/session', {
          params: {
            ...newReq,
            user_name: currentUser.name,
            is_admin: false
          }
        })
      ).data
    )
  }
  async getRemoteAppUrl(session_id, remote_app_name) {
    return (
      await Http.get('/vis/session/remoteAppUrl', {
        params: {
          session_id,
          remote_app_name
        }
      })
    ).data?.url
  }
  async pollFetchRequest(session_id: string) {
    return await Http.get('/vis/session/ready', {
      params: {
        session_id
      },
      baseURL: 'api/v1'
    })
  }

  async autoChoseHardware(id: string) {
    return (
      await Http.get('/vis/software/preset', {
        params: {
          software_id: id
        }
      })
    ).data
  }

  // 获取当前公司下，所有活的 3D 会话，产品经理定义：非关闭都算活的
  async getActiveSessions() {
    return (await Http.get('/vis/activeSessions')).data
  }

  async getComboList(zoneId: string, productId: string) {
    return (
      await Http.get('/vis/combo', {
        params: {
          zone: zoneId,
          product_id: productId
        }
      })
    ).data
  }

  async getTerminalLimit() {
    return (await Http.get('/vis/maxTerminalLimit')).data
  }
}

export const hardwareList = new HardwareList()
