function getDiscoviedList(userid) {
  return fetch(`/api/v3/nodediscovery/list`, {
    headers: { 'x-userid': userid }
  }).then(function (response) {
    return response.json()
  })
}

let intervalId = null
const _eventName = 'nodediscovery'

const ctx: Worker = self as any

ctx.addEventListener('message', event => {
  // 每 30 秒获取一次数据
  const { eventName, eventData } = event.data

  if (eventName === _eventName) {
    if (intervalId === null) {
      intervalId = setInterval(async () => {
        const res = await getDiscoviedList(eventData.userId)
        ctx.postMessage({
          eventName: _eventName,
          eventData: { res }
        })
      }, 1000 * 30)
    }
  }
})

export default null as any
