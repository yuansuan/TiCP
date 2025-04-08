const HtmlWebpackPlugin = require('html-webpack-plugin')

const ClassHandler = require('./utils/ClassHandler')
const SymbolHandler = require('./utils/SymbolHandler')

/**
 * iconfont自动下载替换插件
 *
 * @description 配置好iconfont的url后，会自动下载css和font文件，并插入到html中
 * @class AliIconPlugin
 */
class AliIconPlugin {
  constructor(options) {
    this.url = options.url
    this.type = options.type || 'class'
    this.handler =
      this.type === 'class'
        ? new ClassHandler(options)
        : new SymbolHandler(options)
  }

  async apply(compiler) {
    if (!this.url) {
      throw new Error('AliIconPlugin require a url parameter')
    }

    // 更新html，插入对应资源
    compiler.hooks.compilation.tap('AliIconPlugin', (compilation) => {
      HtmlWebpackPlugin.getHooks(compilation).beforeEmit.tapAsync(
        'AliIconPlugin',
        async (data, cb) => {
          // 获取远程的所有文件
          const newData = await this.handler.handle(data)
          cb(null, newData)
        }
      )
    })

    // 生成文件
    compiler.hooks.emit.tap('AliIconPlugin', (compilation) => {
      this.handler.resources.forEach((item) => {
        compilation.assets[`iconfont/${item.name}`] = {
          source() {
            return item.content
          },
          size() {
            return item.content.length
          },
        }
      })
    })
  }
}

module.exports = AliIconPlugin
