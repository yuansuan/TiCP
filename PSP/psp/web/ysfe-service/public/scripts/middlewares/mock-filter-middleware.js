// Copyright (C) 2018 LambdaCal Inc.

/* eslint no-eval: false */

/**
 * Express mock配置中间件
 * @module MockFilterMiddleware
 * @description 根据配置信息 mock 数据
 */
const path = require('path')
const fs = require('fs')
const debug = require('debug')('dev:server:mock')
const Mock = require('mockjs')

/**
 * js 文件动态读取函数
 * @method readJSFile
 * @param {String} filename
 */
function readJSFile(filename) {
  const content = fs.readFileSync(filename, {
    encoding: 'utf8',
  })
  return content ? eval(content) : null
}

module.exports = function MockFilterMiddleware(params) {
  return function mockFilterHandler(req, res, next) {
    const baseUrl = req.url.split('?')[0]
    const method = String(req.method).toLowerCase()

    /** 数据 mock 配置 */
    const rootConfig = params.config
    // 默认 mock 配置
    let mockConfig = {}

    /** 获取 mock 配置文件的信息 */
    const { configFile } = rootConfig
    if (configFile && fs.existsSync(configFile)) {
      // 读取文件内容
      const mockConfigContent = fs.readFileSync(configFile, {
        encoding: 'utf8',
      })
      // 解析配置文件
      mockConfig = mockConfigContent ? eval(mockConfigContent) : {}
    }

    // 配置了数据 mock，Middleware 可触发
    if (rootConfig) {
      // 是否命中标志
      let hit = false

      /** 遍历 mock 配置，匹配 url */
      Object.keys(mockConfig).forEach(key => {
        let data = mockConfig[key]

        /** data 支持对象形式和字符串形式 */
        // 如果 data 为字符串形式，则扩展成对象形式
        if (typeof data === 'string') {
          data = {
            filename: data,
          }
        }

        /** 解析 key（分离 method 和 url） */
        const matchUrl = key.replace(/^(.+)\s+/, '')
        const matchMethodRes = key.match(/^(.+)\s+.*/)
        const matchMethod = matchMethodRes
          ? String(matchMethodRes[1]).toLowerCase()
          : ''

        /**
         * 匹配 url；
         * 匹配 method（如果 matchMethod 为空，则忽略 method；否则需要匹配相应 method）；
         */
        if (
          baseUrl === matchUrl &&
          (!matchMethod || (matchMethod && matchMethod === method))
        ) {
          // 获取 mockData
          const mockRes = readJSFile(path.join(rootConfig.root, data.filename))

          // 设置为已命中
          hit = true

          debug(
            `mock ${
              matchMethod ? `${matchMethod}:` : ''
            }${matchUrl} with delay ${data.delay || 0}`
          )
          // mock 延时处理
          setTimeout(() => {
            res.send(Mock.mock(mockRes))
          }, data.delay || 0)
        }
      })

      // 如果未命中，则转到下一个中间件处理
      if (!hit) {
        next()
      }
    } else {
      // 未配置数据 mock，则转到下一个中间件处理
      next()
    }
  }
}
