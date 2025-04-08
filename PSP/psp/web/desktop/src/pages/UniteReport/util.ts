import html2canvas from 'html2canvas'
import XLSX from 'xlsx-color'
import { Http } from '@/utils'
import { currentUser } from '@/domain'

interface excelData {
  excelName: string
  sheets: {
    sheetName: string
    data: Array<Array<string | number>>
  }[]
}

const sheetHeaderStyle = {
  font: {
    name: '微软雅黑',
    sz: 14,
    color: { rgb: 'ffffff' },
    bold: true,
    italic: false,
    underline: false,
  },
  fill: {
    patternType: 'solid',
    fgColor: { rgb: '4f80bd' },
  },
}

function formatWSHeader(worksheet, len, style) {
  'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    .substring(0, len)
    .split('')
    .forEach(colNo => {
      // ws[`A1`].v = 'Test'
      worksheet[`${colNo}1`].s = style
    })
}

function formatWSContent(worksheet, oddStyle, evenStyle) {
  const range = XLSX.utils.decode_range(worksheet['!ref'])
  // note: range.s.r + 1 skips the header row
  for (let row = range.s.r + 1; row <= range.e.r; ++row) {
    for (let col = range.s.c; col <= range.e.c; ++col) {
      const ref = XLSX.utils.encode_cell({ r: row, c: col })
      if (worksheet[ref]) {
        if (row % 2 === 1) {
          worksheet[ref].s = oddStyle
        } else {
          worksheet[ref].s = evenStyle
        }
      }
    }
  }
}

export function exportExeclByUrl(type, start, end, execlName) {
  download(`/api/v1/report/export?type=${type}&start=${start}&end=${end}`, execlName)
}

export function exportExcel(data: excelData, auditLogBody?: any) {
  const wb = XLSX.utils.book_new()

  data.sheets.forEach(d => {
    const sheet = XLSX.utils.aoa_to_sheet(d.data)
    XLSX.utils.book_append_sheet(wb, sheet, d.sheetName)

    sheet['!cols'] = [
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
      { wpx: 200 },
    ].slice(0, d.data[0].length)

    formatWSHeader(sheet, d.data[0].length, sheetHeaderStyle)

    let fontStyle = {
      name: '微软雅黑',
      sz: 12,
      color: { rgb: '000000' },
      bold: false,
      italic: false,
      underline: false,
    }

    formatWSContent(
      sheet,
      {
        font: fontStyle,
        fill: {
          patternType: 'solid',
          fgColor: { rgb: 'dbe6f1' },
        },
      },
      {
        font: fontStyle,
        fill: {
          patternType: 'solid',
          fgColor: { rgb: 'b5cce3' },
        },
      }
    )
  })

  const s = XLSX.write(wb, {
    bookType: 'xlsx',
    bookSST: false,
    type: 'binary',
  })
  let buf = new ArrayBuffer(s.length)
  let view = new Uint8Array(buf)
  for (let i = 0; i !== s.length; ++i) {
    view[i] = s.charCodeAt(i) & 0xff
  }

  let tmpDown = new Blob([buf], { type: '' })
  download(URL.createObjectURL(tmpDown), data.excelName + '.xlsx')

  let body = auditLogBody || {
    OperateType: 'REPORT_MGR',
    OperateContent: `用户${currentUser.name}导出报表数据${data.excelName}.xlsx`,
  }

  // TODO
  // Http.post('/audit/log', body)
}

export async function exportImage(
  el: HTMLElement,
  name: string,
  auditLogBody?: any
) {
  const rect = el.getBoundingClientRect()
  const canvas = await html2canvas(el, {
    x: rect.left + window.pageXOffset,
    y: rect.top + window.pageYOffset,
  })
  download(canvas.toDataURL(), name)

  let body = auditLogBody || {
    OperateType: 'REPORT_MGR',
    OperateContent: `用户${currentUser.name}导出报表图片${name}.png`,
  }
  // TODO
  // Http.post('/audit/log', body)
}

function download(href: string, name: string) {
  const a = document.createElement('a')
  a.style.display = 'none'
  document.body.appendChild(a)
  a.href = href
  a.download = name
  a.click()
  a.remove()
}
