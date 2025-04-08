import { getBytes } from '.'
import { Http } from '@/utils'
import qs from 'query-string'

export const isFileExist = path => {
  return Http.post('/file/exist', { paths: [path] }).then(
    res => res.data.isExist
  )
}

export const validateName = name => {
  if (!name.trim()) {
    return '姓名不能为空'
  }

  if (!validateName.reg.test(name)) {
    return '姓名只能包含中文、英文、空格和 , . \' -'
  }

  return true
}
validateName.reg = /^[\u4e00-\u9fa5a-z ,.'-]+$/i

export const validateFilename = name => {
  if (!name) {
    return '文件名不能为空'
  }

  if (!validateFilename.reg.test(name)) {
    return '文件名中包含非法字符：斜杆、反斜杆、单引号、双引号、反引号、逗号或分号'
  }
  if (getBytes(name) > 255) {
    return '文件名长度不能超过 255 字节'
  }

  return true
}
validateFilename.reg = /^[^\\/,;'`"]+$/

export const formatByte = (byte: number) => {
  if (!byte || byte <= 1) return '0.00 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const num = Math.floor(Math.log2(byte) / 10)
  return `${(byte / 2 ** (10 * num)).toFixed(2)} ${units[num]}`
}

export const getUrlParams = () =>
  qs.parse(window.location.search || window.location.hash.split('?')[1])

export function copyToClipboard(s: string) {
  const el = document.createElement('textarea')
  el.value = s
  el.setAttribute('readonly', '')
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  document.body.appendChild(el)
  el.select()
  document.execCommand('copy')
  document.body.removeChild(el)
}

// TODO: wait grpc isDir and Dir exist
export const isDirExist = path => {}

export const filename = name => {
  if (/[\\/,;'`"]/.test(name)) {
    return {
      error: new Error(
        '文件名中包含非法字符：斜杆、反斜杆、单引号、双引号、反引号、逗号或分号'
      )
    }
  }
  if (getBytes(name) > 255) {
    return { error: new Error('文件名长度不能超过 255 字节') }
  }

  return {}
}

const validDomainNameReg =
  /^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$/
const validIpReg =
  /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
const validNameReg = /^[a-zA-Z0-9_]*$/
const validHostNameReg = /^[a-zA-Z0-9_\-\.]+$/
const validInputNameReg = /^[a-zA-Z0-9_\u4e00-\u9fa5]*$/
const validTextAreaReg = /^[a-zA-Z0-9_\u4e00-\u9fa5\n\s]*$/
const validEmailReg =
  /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/
const validPhoneNumber =
  /^(?=\d{11}$)^1(?:3\d|4[57]|5[^4\D]|66|7[^249\D]|8\d|9[89])\d{8}$/
const validPathReg = /^\//
const validLicenseNameReg = /^[\[\]a-zA-Z0-9_\-]+$/

export const isValidInputName = str => validInputNameReg.test(str)

export const isValidTextArea = str => validTextAreaReg.test(str)

export const isValidEmail = str => validEmailReg.test(str)

export const isValidPhoneNumber = str => validPhoneNumber.test(str)

export const isValidName = str => validNameReg.test(str)
export const isValidDomainName = str => validDomainNameReg.test(str)
export const isValidHostName = str => validHostNameReg.test(str)
export const isValidIp = str => validIpReg.test(str)
export const isValidPath = str => validPathReg.test(str)
export const isValidLicenseName = str => validLicenseNameReg.test(str)

export const validateInput = (_, value, label, isRequired) => {
  if (!value && isRequired) {
    return Promise.reject(`${label}不能为空`)
  }

  if (
    value !== '' &&
    _?.field === 'image_id' &&
    !/^[a-zA-Z0-9][a-zA-Z0-9-]{0,64}$/.test(value)
  ) {
    return Promise.reject(
      `${label}不支持中文符号，必须在64个字符以内，由字母、数字和中划线组成`
    )
  }
  if (value?.length > 64) {
    return Promise.reject(`${label}必须小于或等于 64 个字符`)
  }
  return Promise.resolve()
}
export const validateScriptText = (_, value, label, isRequired) => {
  if (!value && isRequired) {
    return Promise.reject(`${label}不能为空`)
  }
  if (value?.length > 65535) {
    return Promise.reject(`${label}必须小于或等于 65535 个字符`)
  }
  return Promise.resolve()
}

export const validateDesc = (_, value, label, isRequired) => {
  if (!value && isRequired) {
    return Promise.reject(`${label}不能为空`)
  }
  if (value?.length > 255) {
    return Promise.reject(`${label}必须小于或等于 255 个字符`)
  }
  return Promise.resolve()
}
