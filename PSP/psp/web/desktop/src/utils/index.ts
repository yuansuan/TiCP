/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import { observable, computed, runInAction } from 'mobx'
import { Observable } from 'rxjs'
export { downloadTestFile } from './downloadTestFile'
export { default as history } from './history'
export { RawHttp } from './Http'
export { NewHttp, NewRawHttp } from './NewHttp'
export { Http as v2Http } from './v2Http'
export * from './formatter'
export * from './Visual'
export * from './theme'
export * from './graphqlClient'
export { buryPoint } from './buryPoint'
export { pageStateStore } from './pageStateStore'
export { apolloClient } from './apolloClient'
import { default as history } from './history'
import moment from 'moment'
export * from './apps'
export * from './HttpClient'
export {interceptError, interceptResponse, createHttp} from './HttpClient2'
export type {AxiosInstance} from './HttpClient2/IAxiosInstance'
export * from './single'
export * from './nextTick'
export { default as eventEmitter } from './EventEmitter'
import * as Validator from './Validator'
export { Validator }
import styled from 'styled-components'
import { sysConfig } from '@/domain'
export { calculateCellDimensions } from './calculateCellDimensions'
import { newBoxServer } from '@/server'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { showFailure } from '@/components/NewFileMGT'
import { message } from 'antd'
const server = serverFactory(newBoxServer)

export interface IEventData {
  message: any
}

export const getComputeType = computeType => {
  return (
    sysConfig.globalConfig?.compute_types?.find(
      item => item['compute_type'] === computeType
    )?.show_name || ''
  )
}
export const isContentOverflow = (element, content) => {
  // -创建临时 element
  let cloneElement = element.cloneNode(true)
  cloneElement.style.display = 'inline-block'
  cloneElement.style.opacity = 0
  cloneElement.style.width = element.offsetWidth + 'px'
  document.body.appendChild(cloneElement)

  // 创建临时 span
  let span = document.createElement('span')
  const computedStyle = window.getComputedStyle(element)
  span.innerHTML = content
  span.style.opacity = '0'
  span.style.display = 'inline-block'
  span.style.fontSize = computedStyle.getPropertyValue('font-size')
  span.style.fontFamily = computedStyle.getPropertyValue('font-family')
  document.body.appendChild(span)

  // 得到文本是否超出的 flag
  const isContentOverflow = span.offsetWidth > cloneElement.offsetWidth

  // 移除临时元素
  document.body.removeChild(span)
  document.body.removeChild(cloneElement)

  return isContentOverflow
}
export const formatRegExpStr = str => str.replace(/[|\\{}()[\]^$+*?.]/g, '\\$&')

export { DataDashPlugin } from './tablePlugins/DataDashPlugin'

export const getFilenameByPath = (p: string) => {
  const dashIndex = p.lastIndexOf('/')
  const name = p.substring(dashIndex + 1)
  const dotIndex = name.lastIndexOf('.')
  return name.substring(0, dotIndex > 0 ? dotIndex : name.length)
}
export const getDisplayRunTime = (nums: number) => {
  const hour = Math.floor(nums / 3600)
    .toString()
    .padStart(2, '0')

  const minute = Math.floor((nums % 3600) / 60)
    .toString()
    .padStart(2, '0')
  const second = (nums % 60).toString().padStart(2, '0')

  return `${hour}:${minute}:${second}`
}
export const formatPath = (path: string): string => {
  if (!path?.length) {
    return './'
  } else {
    if (/^\/.*$/.test(path)) {
      return `.${path}`
    } else if (/^[^\/].*$/.test(path)) {
      return `./${path}`
    } else {
      return './'
    }
  }
}

export const ceilNumber = (value: number) => {
  if (!value) return 0
  let bite = 0
  if (value < 10) {
    return 10
  }
  while (value >= 10) {
    value /= 10
    bite += 1
  }
  return Math.ceil(value) * Math.pow(10, bite)
}

export const formatDate = dateISOString => {
  if (dateISOString) {
    return formatUnixTime(dateISOString)
  }
}

export const copyText = (
  text: string,
  success = () => {},
  failure = () => {}
) => {
  const input = document.createElement('input')
  document.body.appendChild(input)
  input.setAttribute('value', text)
  input.select()

  document.execCommand('copy') ? success() : failure()

  document.body.removeChild(input)
}

export const getSearchParamByKey = (searchParmasStr, key) => {
  return new URLSearchParams(searchParmasStr).get(key)
}

export const getMonths = (start, end) => {
  let months = []

  while (start.format('YYYY-MM') !== end.format('YYYY-MM')) {
    months.push(start.format('YYYY-MM'))
    start.add(1, 'months')
  }

  months.push(end.format('YYYY-MM'))

  return months
}

export const copy2clipboard = text => {
  if (navigator.clipboard) {
    // clipboard api 复制
    navigator.clipboard.writeText(text)
  } else {
    const textarea = document.createElement('textarea')
    document.body.appendChild(textarea)
    // 隐藏此输入框
    textarea.style.position = 'fixed'
    textarea.style.clip = 'rect(0 0 0 0)'
    textarea.style.top = '10px'
    // 赋值
    textarea.value = text
    // 选中
    textarea.select()
    // 复制
    document.execCommand('copy', true)
    // 移除输入框
    document.body.removeChild(textarea)
  }
}

export function extractPathAndParamsFromURL(url): any {
  if (!url) return ''
  let startPath = url?.split('/')[1]
  let extractedPath = startPath?.split('?')[0]
  const endIndex = url?.indexOf('?')
  const params = {}

  if (endIndex !== -1) {
    const queryString = url?.substring(endIndex + 1)
    const paramPairs = queryString?.split('&')

    paramPairs?.forEach(pair => {
      const [key, value] = pair.split('=') // 拆分参数名和参数值
      params[key] = value // 将参数名和参数值添加到 params 对象中
    })
  }

  return {
    path: extractedPath,
    ...params
  }
}

export function getUrlParams(url?: string): any {
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH') || ''
  return extractPathAndParamsFromURL(url || currentPath)
}
export function parseUrlParam(url) {
  if (!url.includes('?')) return {}
  const params = {}
  // 截取问号之后的的参数
  const urlParams = url.slice(1)
  // 根据&截取
  const dataArr = urlParams.split('&')
  dataArr.forEach(item => {
    const [key, value] = item.split('=')
    params[key] = value
  })
  return params
}

export const createMobxStream = (selector, immediately = true) =>
  new Observable(function (observer) {
    const computedValue = computed(selector)
    return computedValue.observe(
      ({ newValue }) => observer.next(newValue),
      immediately
    )
  })

export const fromStream = <T>(
  stream$: Observable<T>,
  mobxValue?: { current: T; dispose?: () => void }
) => {
  const observer =
    mobxValue || observable({ current: null, dispose: undefined })

  observer.dispose = stream$.subscribe((val: any) => {
    runInAction(() => {
      observer.current = val
    })
  })

  return observer
}

export function checkImgExists(imgurl) {
  return new Promise(function (resolve, reject) {
    let ImgObj = new Image()
    ImgObj.src = imgurl
    ImgObj.onload = function (res) {
      resolve({ success: true, res })
    }
    ImgObj.onerror = function (err) {
      resolve({ success: false, err })
    }
  })
}
export const getBytes = name =>
  encodeURI(name).split(/%(?:u[0-9A-F]{2})?[0-9A-F]{2}|./).length - 1
interface IGoogleTimestamp {
  seconds: number
  nanos: number
}

export class Timestamp implements IGoogleTimestamp {
  @observable seconds = null
  @observable nanos = null

  constructor(timestampObj?: IGoogleTimestamp) {
    if (!timestampObj) return
    this.nanos = timestampObj.nanos
    this.seconds = timestampObj.seconds
  }

  set moment(moment) {
    if (moment) {
      this.seconds = Math.floor(moment.valueOf() / 1000)
      this.nanos = moment.valueOf() % 1000
    } else {
      this.seconds = null
      this.nanos = null
    }
  }

  @computed
  get dateString() {
    return this.formatDate()
  }

  formatDate(format = 'YYYY-MM-DD HH:mm:ss') {
    if (this.seconds === null && this.nanos === null) return ''
    const milliSec = parseInt(this.seconds) * 1000 + this.nanos
    return moment(milliSec).format(format)
  }

  @computed
  get moment() {
    return this.seconds
      ? moment(parseInt(this.seconds) * 1000 + this.nanos)
      : null
  }
}

export const needLogin = localStorage.getItem('needLogin') === 'true'

type HistoryType = 'browser' | 'hash'

export const getSelectedKey = (routerType?: HistoryType) => {
  if (typeof window === 'undefined') {
    return undefined
  }

  const type = typeof history === 'string' ? history : routerType

  if (type === 'browser') {
    return window.location.pathname.replace(/\/$/, '') || '/'
  } else {
    return (window.location.hash.replace(/^#/, '') || '/').split('?')[0]
  }
}

export const Tips = styled.span`
  font-family: PingFangSC-Regular;
  font-size: 12px;
  color: #999999;
  line-height: 22px;
`

export function getByteLength(str) {
  const encoder = new TextEncoder()
  const byteArray = encoder.encode(str)
  return byteArray.length
}

export function download(href: string, name: string) {
  const a = document.createElement('a')
  a.style.display = 'none'
  document.body.appendChild(a)
  a.href = href
  a.download = name
  a.click()
  a.remove()
}

export async function moveTo(path, selectedNodes) {
  function getCurrentDir(file_path) {
    if (!file_path) return '/'
    let last_slash_index = file_path.lastIndexOf('/')
    let current_dir = file_path.substring(0, last_slash_index).slice(0)
    return current_dir
  }

  const nodes = selectedNodes

  // check duplicate
  const targetDir = await server.fetch(path)
  const rejectedNodes = []
  const resolvedNodes = []

  nodes.forEach(item => {
    if (targetDir.getDuplicate({ id: undefined, name: item.name })) {
      rejectedNodes.push(item)
    } else {
      resolvedNodes.push(item)
    }
  })

  const destMoveNodes = [...resolvedNodes]
  if (rejectedNodes.length > 0) {
    const coverNodes = await showFailure({
      actionName: '移动',
      items: rejectedNodes
    })
    if (coverNodes.length > 0) {
      // coverNodes 中的要删除
      // await server.delete(coverNodes.map(item => `${path}/${item.name}`))
      destMoveNodes.push(...coverNodes)
    }
  }

  if (destMoveNodes.length > 0) {
    const srcPaths = destMoveNodes[0]?.selectedname
      ? destMoveNodes[0]?.selectedname
          .split(',')
          .map(name => `${getCurrentDir(destMoveNodes[0].path)}/${name}`)
      : [`${getCurrentDir(destMoveNodes[0].path)}/${destMoveNodes[0].name}`]

    const destPath = path ? path : '/'
    await server.move({ srcPaths, destPath, overwrite: true })

    await server.fetch(path)

    message.success('文件移动成功')
  }
}

export function escapeRegExp(text) {
  return text.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, '\\$&')
}

export const formatUnixTime = (timestamp: number) => {
  return require('dayjs')(timestamp * 1000).format('YYYY/MM/DD HH:mm:ss')
}
