import Http from '@/services/AxiosHttp'

export default {
  closeWorkTask(userId, workTaskId) {
    const data = {
      user_id: userId,
      work_task_id: workTaskId
    }
    return Http.post('/visual/worktask/stop', data)
  },
  getWorkTaskDetail(workTaskId) {
    const url = `/visual/worktask/detail/${workTaskId}`
    return Http.get(url)
  },
  reportAppMetric(data) {
    return Http.post('/visual/metrics/worktask', data)
  },
  activeApp(workTaskId) {
    return Http.put(`/visual/worktask/${workTaskId}`, {
      status: 'active'
    })
  }
}
