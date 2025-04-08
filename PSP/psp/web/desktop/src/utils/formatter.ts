/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import moment from 'moment'

export const formatAmount = (amount: number) => (amount / 100000).toFixed(2)
/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */
export const formatTime = time => {
  const milliSec = parseInt(time) * 1000
  return moment(milliSec).format('YYYY-MM-DD HH:mm:ss')
}

export const formatDateFromMilliSec = milliSec => {
  if (!milliSec) return null
  return moment(milliSec).format('YYYY-MM-DD HH:mm:ss')
}

export const formatDateFromMilliSecWithTimeZone = milliSec => {
  if (!milliSec) return null
  let timeZoneSec = new Date().getTimezoneOffset() * 60
  return moment(milliSec - timeZoneSec).format('YYYY-MM-DD HH:mm:ss')
}

const fileUnits = {
  Byte: 'B',
  Kilobyte: 'KB',
  Megabyte: 'MB',
  Gigabyte: 'GB',
  Terabyte: 'TB',
  Petabyte: 'PB'
}
const NodeUnits = {
  Megabyte: 'MB',
  Gigabyte: 'GB'
}

export const round = (num: number) => {
  return Math.round(num * 100) / 100
}

export const formatNodeAttrNumber = (num: number) => {
  let unit = 'MB'
  let size = num
  let units = Object.keys(NodeUnits)
  units.shift()

  units.reduce((result, key) => {
    result = result / 1024
    if (result >= 1) {
      unit = NodeUnits[key]
      size = result
      return result
    } else {
      units = []
      return undefined
    }
  }, num)

  return `${round(size)}${unit}`
}
export const formatlsfAttrNumber = attr => {
  const t = typeof attr
  if (t === 'number') {
    return round(attr)
  } else {
    return '--'
  }
}

/**
 * @name formatFileSize : format Byte to Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte
 * @param number unit is Byte
 *
 */
export const formatFileSize = (num: number) => {
  let unit = 'B'
  let size = num
  let units = Object.keys(fileUnits)
  units.shift()

  units.reduce((result, key) => {
    result = result / 1024
    if (result >= 1) {
      unit = fileUnits[key]
      size = result
      return result
    } else {
      units = []
      return undefined
    }
  }, num)

  return `${round(size)}${unit}`
}

export const roundNumber = (value, precision) => {
  const step = Math.pow(10, precision)
  return Math.round(value * step) / step
}

export const toTimeDuration = s => {
  if (s === -1 || s === null || s === undefined) return '--'
  let hour = Math.floor(s / (60 * 60))
  let minu = Math.floor((s / 60) % 60)
  let sec = Math.floor(s % 60)

  return `${hour}:${minu}:${sec}`
}

export const formatCreditQuota = number =>
  new Intl.NumberFormat('zh-CN', {
    currency: 'CNY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(number)

export const formatMoney = number =>
  new Intl.NumberFormat('zh-CN', {
    currency: 'CNY',
    minimumFractionDigits: 2,
    maximumFractionDigits: 4
  }).format(number)

export const formatMoneyWithUnit = (number, digits = 2) =>
  new Intl.NumberFormat('zh-CN', {
    currency: 'CNY',
    style: 'currency',
    minimumFractionDigits: 2,
    maximumFractionDigits: digits
  }).format(number)

export const formatCulearTime = number =>
  new Intl.NumberFormat('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 5
  }).format(number)

export const replaceRoleTag = str =>
  str.replace(/<Role\s+id=\{(\d+)\}>/g, '').replace(/<\/Role>/g, '')

export const replaceGroupTag = str =>
  str.replace(/<Group\s+id=\{(\d+)\}>/g, '').replace(/<\/Group>/g, '')

export const replaceUserTag = str =>
  str.replace(/<User\s+id=\{(\d+)\}>/g, '').replace(/<\/User>/g, '')

export const replaceCustomTags = str =>
  replaceUserTag(replaceGroupTag(replaceRoleTag(str)))

export const isISO = (input: any) =>
  moment(input, moment.ISO_8601, true).isValid()

/**
 * dateStr: a string like `2021-01-08T14:42:34.678Z` (format: ISO 8601).
 * format: target format string like "YYYY-MM-DD HH:mm:dd"
*/
export const formatISODateStr = (dateStr, format: string) => {
  if (isISO(dateStr)) {
    return moment(dateStr).format(format)
  } else {
    return '--'
  }
}

export const formatTimestamp = (timestamp: string) => {
  const date = new Date(timestamp)
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return hours + ':' + minutes
}