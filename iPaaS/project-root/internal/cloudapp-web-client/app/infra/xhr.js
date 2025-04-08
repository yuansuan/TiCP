// Copyright (C) 2018 LambdaCal Inc.

import req from 'superagent'

function formatBodyResponse(resp) {
  if (!resp.ok) {
    throw resp.error
  }
  return resp.body
}

class Xhr {
  get(url) {
    return req.get(url).use(this._formatRequest.bind(this))
  }

  post(url) {
    return req.post(url).use(this._formatRequest.bind(this))
  }

  put(url) {
    return req.put(url).use(this._formatRequest.bind(this))
  }

  delete(url) {
    return req.del(url).use(this._formatRequest.bind(this))
  }

  _formatRequest(sReq) {
    const rawResponseThen = sReq.then.bind(sReq)
    sReq.then = (onFulfilled, onRejected) =>
      rawResponseThen(formatBodyResponse.bind(this)).then(
        onFulfilled,
        onRejected
      )
  }
}

export default new Xhr()
