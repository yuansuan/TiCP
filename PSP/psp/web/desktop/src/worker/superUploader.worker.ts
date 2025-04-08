function getUploaderFileList(taskKey) {
  return fetch(
    `/api/v1/storage/hpcUpload/fileTaskList?taskKey=${taskKey}`
  ).then(function (response) {
    return response.json()
  })
}

let intervalId = null
const _eventName = 'super_uploader'

const ctx: Worker = self as any

ctx.addEventListener('message', event => {
  // 每 1 秒获取一次数据
  const { eventName, eventData } = event.data

  if (eventName === _eventName) {
    if (intervalId === null) {
      intervalId = setInterval(async () => {
        const res = await getUploaderFileList(eventData.taskKey)
        ctx.postMessage({
          eventName: _eventName,
          eventData: { res }
        })
      }, 1000)
    }
  }
})

export default null as any
