import qs from 'qs'
import XLSX from 'xlsx-color'

const _eventName = 'jobReport'
const ctx: Worker = self as any

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


export function makeExcelBlob(data: excelData) {
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
  
  return tmpDown
}

async function getTotal(url, params) {
  let fetchUrl = `/api/v1${url}?${qs.stringify(params)}`
  return fetch(fetchUrl).then(function (response) {
    return response.json()
  })
}

// 然后，1页页取数据，并发为 2
async function getPageData(url, {start, size=1000, ...params}) {
  let fetchUrl = `/api/v1${url}?${qs.stringify({
    ...params,
    page_size: size,
    page_index: start,
  })}`
  return fetch(fetchUrl).then(function (response) {
    return response.json()
  })
}

class SupperTask {
  size = null
  queue = null
  runningCount = 0

  constructor(size) {
    this.size = size || 2
    this.queue = []
    this.runningCount = 0
  }

  add(fn) {
    return new Promise((resolve, reject) => {
      this.queue.push({fn, resolve, reject})
      this._run()
    })
  }

  _run() {
    while (this.runningCount < this.size && this.queue.length) {
      const {fn, resolve, reject} = this.queue.shift()
      
      this.runningCount++

      fn().then(resolve, reject).finally(() => {
        this.runningCount--
        this._run()
      })
    }
  }
}

export const supperTask = new SupperTask(1);

export async function startExport(url, params, callback, dataKey) {
  let start = 1, 
  totalNumbers = 1000,
  size = 1000,
  allData = []

  const N = 30
  let partIndex = 1

  const { data: { total }} = await getTotal(url, {...params})

  if (total === 0) {
    ctx.postMessage({
      eventName: _eventName,
      eventData: {error: 'no data to export'},
    })
    return
  }

  totalNumbers = total

  let successArray = new Array(totalNumbers%size === 0 ? totalNumbers/size :  Math.floor(totalNumbers/size) + 1).fill(false) 

  function addTask(params) {
    supperTask
    .add(() => getPageData(url, params))
    .then((res) => {
      const { data } = res as any
      allData[params.start-1] = data[dataKey] || []
      successArray[params.start-1] = true

      ctx.postMessage({
        eventName: _eventName,
        eventData: { 
          running: true,
          error: null,
          data: null
        },
      })

      if (params.start%N === 0) {
        // 部分导出
        callback && callback(allData, partIndex++)
        allData = []
      }

      if (successArray.every(succ => succ)) {
        // 导出成功
        callback && callback(allData, partIndex++, true)
      }

    }).catch((err) => {
      // 导出失败
      successArray[params.start-1] = false
      supperTask.queue = []

      ctx.postMessage({
        eventName: _eventName,
        eventData: { 
          error: err,
          running: false,
          data: null
        },
      })
    })
  }

  do {
    addTask({...params, start, size})
    start += 1
  } while (totalNumbers > (start - 1) * size)
}


ctx.addEventListener('message', event => {
  // 获取开始下载数据信号，开始请求数据
  const { eventName, eventData } = event.data
  if (eventName === _eventName) {
    const { params, type, computeTypesMap } = eventData

    const map = {
      overview : 'overviews',
      detail: 'job_details'
    }

    const columnName = params.query_type === 'app' ? '应用' : '用户'

    const execlInfo = {
      overview: {
        execlName: '统计总览',
        columnNames: [`${columnName}编号`, `${columnName}名称`, '计算类型', '使用时长(小时)'],
        columnKeys: ['id', 'name', 'computeType', 'cpu_time'],
        formatter: [val => val, val => val, val => computeTypesMap[val], val => val],
        url: '/job/statistics/overview',
      },
      detail: {
        execlName: '统计明细',
        columnNames: ['作业编号', '作业名称', '计算类型', '应用名称', '用户名称', '提交时间', '开始时间', '结束时间', '使用时长(小时)'],
        columnKeys: ['id', 'name', 'type', 'app_name', 'user_name', 'submit_time','start_time', 'end_time', 'cpu_time'],
        formatter: [val => val, val => val, val => computeTypesMap[val], val => val, val => val, val => val, val => val, val => val, val => val],
        url: '/job/statistics/detail',
      }
    }

    const dataKey = map[type]
  
    const { url, execlName, columnKeys, columnNames, formatter} = execlInfo[type]


    const callback = (allData, partIndex, isLast = false) => {
      const sheetNameMap = {
        name: execlName,
      }

      let sheetData = allData.flat().map(d => {
        let row = []
        columnKeys.forEach((key, index) => {
          row.push(formatter[index](d[key]))
        })

        return row
      })

      sheetData.unshift(columnNames)
  
      const sheets = []

      sheets.push({
        sheetName: sheetNameMap.name,
        data: sheetData,
      })

      const blob = makeExcelBlob({
        excelName: `${execlName}报表`,
        sheets,
      })

      ctx.postMessage({
        eventName: _eventName,
        eventData: { 
          data: {
            blob,
            execlName: `${execlName}报表`,
          },
          error: null,
          running: !isLast,
          partIndex,
        },
      })
    }

    startExport(url, params, callback, dataKey)
  }
})

export default null as any