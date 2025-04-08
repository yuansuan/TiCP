const BaseHandler = require('./BaseHandler')

class SymbolHandler extends BaseHandler {
  /**
   * 单字体文件下载：
   * 串行下载多个文件：[url_01, url_02]
   * 并行下载多个文件：[[url_01, url_02]]
   */
  async handle(data) {
    const urls = (typeof this.url === 'string' ? [this.url] : this.url).filter(
      Boolean
    )
    if (urls.length === 0) {
      return data
    }

    for (let i in urls) {
      const iconUrls = (typeof urls[i] === 'string'
        ? [urls[i]]
        : urls[i]
      ).filter(Boolean)

      if (iconUrls.length > 0) {
        await Promise.all(
          iconUrls.map(async (url) => {
            // js content
            const content = await this.getRemoteContent(url)
            const name = this.getNameFromUrl(url)
            this.resources.push({ name: `${name}.js`, content })
            const iconJsLink = `<script src="${this.publicPath}iconfont/${name}.js"></script>`
            data.html = this.insertTag(data.html, iconJsLink)
          })
        )
      }
    }
    return data
  }
}

module.exports = SymbolHandler
