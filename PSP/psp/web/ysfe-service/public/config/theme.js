/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

function format(str) {
  return str.replace(/([A-Z1-9])/g, '-$1').toLowerCase()
}

const fontSize = ['12px', '14px', '16px', '20px', '24px', '30px', '38px']

const _antTheme = {
  fontFamily: 'PingFangSC-Regular',
  btnBorderRadiusBase: '2px',

  primaryColor: '#005dfc',
  linkColor: '#3182ff',

  errorColor: '#f5222d',
  warningColor: '#ffa726',
  successColor: '#52c41a',
  infoColor: '#3182ff',

  borderColorBase: 'rgba(0,0,0,0.10)',
  borderColorSplit: 'rgba(0,0,0,0.10)',

  // disabled bg color
  backgroundColorBase: '#f5f5f5',
  disabledColor: 'rgba(0,0,0,0.25)',

  fontSizeSm: fontSize[0],
  fontSizeBase: fontSize[1],
  fontSizeLg: fontSize[2],
  heading1Size: fontSize[6],
  heading2Size: fontSize[5],
  heading3Size: fontSize[4],
  heading4Size: fontSize[3],
  heading5Size: fontSize[2],
}

const antTheme = Object.entries(_antTheme).reduce(
  (res, [key, value]) => {
    res[format(key)] = value
    return res
  },
  {}
)

const theme = {
  ..._antTheme,
  backgroundColorHover: '#F6F8FA',
  secondaryColor: '#3182FF',
  cancelColor: '#BFBFBF',
  cancelHighlightColor: '#8C8C8C',

  fontSize,
}
