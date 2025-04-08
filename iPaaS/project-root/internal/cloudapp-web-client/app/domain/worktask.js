import xhr from '@/infra/xhr'
import log from '@/infra/log'

export default class WorkTask {
  constructor(id, userId, roomId) {
    this._id = id
    this._userId = userId
    this._roomId = roomId
  }

  queryInfo() {
    const detailRoomId = this._roomId.split('-')[0]
    return xhr.get(`/visual/worktask/detail/${detailRoomId}`).then(res => {
      if (res.code !== 0) {
        log.error('fail to get worktask details, reason: ', res.message)
        return null
      }
      return res.data
    })
  }

  close() {
    return xhr
      .post(`/visual/worktask/stop`)
      .send({ user_id: this._userId, work_task_id: this._id })
  }
}
