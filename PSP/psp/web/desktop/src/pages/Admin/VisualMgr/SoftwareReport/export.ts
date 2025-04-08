import { Http } from '@/utils'
import { eventEmitter, IEventData } from '@/utils'
import { message } from 'antd'

import { DATE_FORMAT } from './data'

async function getTotal(url, {time, appIds}) {
  return Http.post(url, {
    start_time: time[0]?.format(DATE_FORMAT) || null,
    end_time: time[1]?.format(DATE_FORMAT) || null,
    app_ids: appIds.length === 0 ? null : appIds,
    page_size: 1,
    page_num: 1,
  }, { baseURL: '' })
}

// 然后，1页页取数据，并发为 3
async function getPageData(url, {time, appIds, start, size=1000}) {
  return Http.post(url, {
    start_time: time[0]?.format(DATE_FORMAT) || null,
    end_time: time[1]?.format(DATE_FORMAT)|| null,
    app_ids: appIds.length === 0 ? null : appIds,
    page_size: size,
    page_num: start,
  }, { baseURL: '' })
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

const supperTask = new SupperTask(3);

export async function startExport(url, time, appIds, callback) {
  let start = 1, 
  totalNumbers = 1000,
  size = 1000,
  allData = []

  const { data: { total }} = await getTotal(url, {time, appIds})

  if (total === 0) {
    eventEmitter.emit('export_execl_error', {
      message: { error: 'no data to export' },
    } as IEventData)
    message.error('没有数据，无法导出')
    return
  }

  totalNumbers = total

  let successArray = new Array(totalNumbers%size === 0 ? totalNumbers/size :  Math.floor(totalNumbers/size) + 1).fill(false) 

  function addTask(params) {
    supperTask
    .add(() => getPageData(url, params))
    .then((res) => {
      const { data } = res as any
      allData = allData[params.start-1] = data?.list
      successArray[params.start-1] = true

      if (successArray.every(succ => succ)) {
        // 导出成功
        eventEmitter.emit('export_execl_success', {
          message: { allData, },
        } as IEventData)

        callback && callback(allData)
      }

    }).catch((err) => {
      // 导出失败
      successArray[params.start-1] = false
      supperTask.queue = []
      eventEmitter.emit('export_execl_error', {
        message: { error: err },
      } as IEventData)

      message.error('导出失败')
    })
  }

  do {
    addTask({time, appIds, start, size})
    start += 1
  } while (totalNumbers > (start - 1) * size)
}
