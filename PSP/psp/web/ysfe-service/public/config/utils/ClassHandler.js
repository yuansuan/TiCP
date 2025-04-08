const BaseHandler = require('./BaseHandler')

class ClassHandler extends BaseHandler {
  async get(iconUrl) {
    // 获取远程css内容
    let cssContent = await this.getRemoteContent(iconUrl)
    const reg = /url\('(.*?)'\)/g
    const fonts = Array.from(Array(5)).map((_) => {
      if (reg.test(cssContent)) {
        return reg.exec(cssContent)[1]
      }
      return ''
    })

    // 下载字体文件
    for (let index in fonts) {
      const url = fonts[index]
      // 如果不是带文件名的url
      if (!/^.*\/(.*?)\?/.test(url)) {
        continue
      }
      // 从远程的url获取对应的内容
      let font = { name: /^.*\/(.*?)\?/.exec(url)[1], url }
      const content = await this.getRemoteContent(font.url)
      this.resources.push({ name: font.name, content })
      const regStr = font.url.replace('?', '\\?')
      cssContent = cssContent.replace(new RegExp(regStr, 'g'), `./${font.name}`)
    }
    const name = this.getNameFromUrl(iconUrl)
    this.resources.push({ name: `${name}.css`, content: cssContent })
  }

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
        console.log('正在下载资源：', iconUrls.join(', '))
        await Promise.all(
          iconUrls.map(async (url) => {
            await this.get(url)
            const name = this.getNameFromUrl(url)
            const iconCssLink = `<link rel="stylesheet" href="${this.publicPath}/iconfont/${name}.css">`
            data.html = this.insertTag(data.html, iconCssLink)
          })
        )
        console.log('下载完成')
      }
    }

    return data
  }
}

module.exports = ClassHandler
