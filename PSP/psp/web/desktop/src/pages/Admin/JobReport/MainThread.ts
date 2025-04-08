import Worker from '@/worker/jobReport.worker'
import { eventEmitter, IEventData } from '@/utils'
import { message } from 'antd'
import { download } from './exportData'

function workerFactory(type) {
  const worker = new Worker()
  worker.addEventListener('message', event => {
    const { eventName, eventData } = event.data
    const { data, error, running, partIndex } = eventData

    if (eventName === _eventName) {

      console.debug('receive message (jobReport worker)', type, running)

      if (running) {
        eventEmitter.emit(type + '_export_execl_running', {
          message: { error: null },
        } as IEventData)

        if (data) {
          console.debug('receive message (jobReport worker) data part success', type, running, partIndex)
          download(URL.createObjectURL(data.blob), data.execlName + `_s${partIndex}.xlsx`)
        }

        return 
      }
      
      if (error) {
        if (error === 'no data to export') {
          message.error('没有数据，无法导出')
          eventEmitter.emit(type + '_export_execl_error', {
            message: { error: 'no data to export' },
          } as IEventData) 
        } else {
          message.error('导出数据出错了')
          eventEmitter.emit(type + '_export_execl_error', {
            message: { error },
          } as IEventData) 
        }

      } else {

        console.debug('receive message (jobReport worker) success', type, running)

        download(URL.createObjectURL(data.blob), data.execlName + (partIndex === 1 ? '.xlsx' : `_s${partIndex}.xlsx`))
        
        eventEmitter.emit(type + '_export_execl_success', {
          message: { error: null },
        } as IEventData)
      }
    }
  })

  worker.addEventListener("error", (event) => {
    console.error('worker error (jobReport worker)', type, event)

    message.error('导出数据出错了')
    
    eventEmitter.emit(type + '_export_execl_error', {
      message: { error: 'worker error' },
    } as IEventData) 
  });

  worker.addEventListener("messageerror", (event) => {
    console.error('worker messageerror (jobReport worker)', type, event)

    message.error('导出数据出错了')
    
    eventEmitter.emit(type + '_export_execl_error', {
      message: { error: 'worker messageerror' },
    } as IEventData) 
  });

  return worker
}

export const appWorker = workerFactory('app')
export const userWorker = workerFactory('user')

export const _eventName = 'jobReport'
