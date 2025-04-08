const http = require('http')

class BaseHandler {
  constructor(options) {
    this.resources = []
    this.url = options.url
    this.publicPath = options.publicPath || './'

    // format publicPath
    if (this.publicPath[this.publicPath.length - 1] !== '/') {
      this.publicPath += '/'
    }
  }

  /**
   * 获取远程的内容
   *
   * @memberof AliIconPlugin
   */
  getRemoteContent(url) {
    url = 'http:' + url
    return new Promise((resolve) => {
      http.get(url, (res) => {
        let content = ''
        res.on('data', (chunk) => {
          content += chunk
        })
        res.on('end', () => {
          resolve(content)
        })
      })
    })
  }

  /**
   * 插入标签
   *
   * @param {*} html
   * @param {*} tag
   * @returns
   * @memberof BaseHandler
   */
  insertTag(html, tag) {
    const index = html.indexOf('</head>')
    return html.substring(0, index) + tag + html.substring(index)
  }

  /**
   * 处理
   *
   * @abstract
   * @memberof BaseHandler
   */
  handle() {
    throw new Error('You should not call a abstract function!')
  }

  getNameFromUrl(url) {
    return /(.*\/)*([^.]+).*/.exec(url)[2]
  }
}

module.exports = BaseHandler
