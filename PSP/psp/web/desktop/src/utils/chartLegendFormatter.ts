export const legendFormatterByLength = (string, maxLength) => {
  if (string.replace(/[\u4e00-\u9fa5]/g, '**').length <= maxLength) {
    return string
  } else {
    let len = 0
    let tmpStr = ''
    for (let i = 0; i < string.length; i++) {
      if (/[\u4e00-\u9fa5]/.test(string[i])) {
        len += 2
      } else {
        len += 1
      }
      if (len > maxLength - 2) {
        break
      } else {
        tmpStr += string[i]
      }
    }
    return tmpStr + '...'
  }
}

export const legendFormatter = string => {
  return legendFormatterByLength(string, 10)
}