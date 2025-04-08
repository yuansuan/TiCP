import moment from 'moment'

function getDates() {
  return [moment().subtract(1, 'days'), moment()].map(m => m.valueOf())
}

function getUserList(userid) {
  return fetch(`/api/v1/auth/onlineList?_timestamp_=${Date.now()}`, {
    method: 'POST',
    headers: { 'x-userid': userid },
    body: JSON.stringify({
      index: 1,
      size: 1,
    }),
  }).then(function (response) {
    return response.json()
  }).then(function (res) {
    return {
      data: res.data?.page?.total || 0
    }
  })
}

function getDashboardInfo(type, dates, userid) {
  let url = null
  // ClUSTER_INFO, RESOURCE_INFO, JOB_INFO, SOFTWARE_INFO
  if (type === 'ClUSTER_INFO') {
    url = `/api/v1/dashboard/clusterInfo?start=${dates[0]}&end=${dates[1]}`
  } else if (type === 'RESOURCE_INFO') {
    url = `/api/v1/dashboard/resourceInfo?start=${dates[0]}&end=${dates[1]}`
  } else if (type === 'JOB_INFO') {
    url = `/api/v1/dashboard/jobInfo?start=${dates[0]}&end=${dates[1]}`
  } else if (type === 'SOFTWARE_INFO') {
    url = `/api/v1/job/appJobNum?start=${dates[0]}&end=${dates[1]}`
  } else if (type === 'USER_JOB_INFO') {
    url = `/api/v1/job/userJobNum?start=${dates[0]}&end=${dates[1]}`
  }

  url += `&_timestamp_=${Date.now()}`

  return fetch(
    url,
    {
      headers: { 'x-userid': userid },
    }
  ).then(function (response) {
    return response.json()
  })
}

let intervalId = null
const _eventName = 'dashboard'

const ctx: Worker = self as any

function getDashboardData(dates, userid) {
  // @ts-ignore
  return Promise.all([
    getDashboardInfo('ClUSTER_INFO', [], userid),
    getDashboardInfo('SOFTWARE_INFO', dates, userid),
    getDashboardInfo('RESOURCE_INFO', dates, userid),
    getUserList(userid),
    getDashboardInfo('USER_JOB_INFO', dates, userid),
  ])
}

ctx.addEventListener('message', event => {
  // 每 3 分钟获取 Dashboard 数据
  const { eventName, eventData } = event.data

  if (eventName === _eventName) {
    if (intervalId === null) {
      intervalId = setInterval(async () => {
        const res = await getDashboardData(getDates(), eventData.userId)
        ctx.postMessage({
          eventName: _eventName,
          eventData: { res },
        })
      }, 1000 * 60 * 3)
    }
  }
})

export default null as any
