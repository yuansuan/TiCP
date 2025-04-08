/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

'use strict'

// define development env
process.env.BABEL_ENV = 'development'
process.env.NODE_ENV = 'development'

// handle unhandledRejection
process.on('unhandledRejection', (err) => {
  throw err
})

// Ensure environment variables are read.
require('../config/env')

const paths = require('../config/paths')

// Warn and crash if required files are missing
const checkRequiredFiles = require('react-dev-utils/checkRequiredFiles')
if (!checkRequiredFiles([paths.appHtml, paths.appIndexJs])) {
  process.exit(1)
}

const http = require('http')
const webpackDevMiddleware = require('webpack-dev-middleware')
const webpackHotMiddleware = require('webpack-hot-middleware')
const historyApiFallback = require('connect-history-api-fallback')
const express = require('express')
const path = require('path')
const { computeConfig } = require('../config/webpack.config.js')
const webpackConfig = computeConfig('development')
const isHttps = webpackConfig.devServer.https
const resolveApp = require('../config/utils/resolveApp')

const https = require('https')
const fs = require('fs-extra')
const privateKey = fs.readFileSync(require.resolve('./ssl/server.key'), 'utf8')
const certificate = fs.readFileSync(require.resolve('./ssl/server.crt'), 'utf8')
const credentials = { key: privateKey, cert: certificate }

function createApp({ port, host }) {
  const webpack = require('webpack')
  const compiler = webpack(webpackConfig)

  const openBrowser = require('react-dev-utils/openBrowser')

  const app = express()

  // dev-middleware
  const devMiddleware = webpackDevMiddleware(compiler, {
    publicPath: webpackConfig.output.publicPath,
    stats: 'errors-only',
  })
  app.use(devMiddleware)

  // dev-middleware onComplete: open browser
  devMiddleware.waitUntilValid(() => {
    const pageUrl = `${
      webpackConfig.devServer.https ? 'https' : 'http'
    }://${host}:${port}`
    openBrowser(pageUrl)
  })
  // hot-middleware
  app.use(webpackHotMiddleware(compiler))

  // mock's priority is higher than proxy
  app.use(
    require('./middlewares/mock-filter-middleware')({
      config: {
        root: resolveApp('config/mock/'),
        configFile: resolveApp('config/mock/config.js'),
      },
    })
  )

  // proxy
  const { createProxyMiddleware } = require('http-proxy-middleware')
  const proxyConfig = webpackConfig.devServer.proxy || {}
  Object.keys(proxyConfig).forEach((key) => {
    const config = proxyConfig[key]
    if (isHttps) {
      config.headers = config.headers || {}
      config.headers['X-Forwarded-Proto'] = 'https'
    }
    app.use(key, createProxyMiddleware(key, config))
  })

  webpackConfig.devServer.historyApiFallback && app.use(historyApiFallback())
  app.use(devMiddleware)

  return app
}

function runInUnusedPort({ port, host }, fn) {
  const net = require('net')
  const tester = net
    .createServer()
    .once('error', function (err) {
      if (err.code !== 'EADDRINUSE') {
        throw err
      }

      // try next port
      console.warn(`port:${port} is in use, check port:${port + 1}`)
      runInUnusedPort({ port: ++port, host }, fn)
    })
    .once('listening', function () {
      tester
        .once('close', function () {
          fn({ port, host })
        })
        .close()
    })
    .listen(port, host)
}

runInUnusedPort(
  {
    port: webpackConfig.devServer.port || 8080,
    host: webpackConfig.devServer.host || 'localhost',
  },
  ({ port, host }) => {
    const app = createApp({ port, host })
    let server = null

    if (isHttps) {
      server = https.createServer(credentials, app)
    } else {
      server = http.createServer(app)
    }

    server.listen(port, host)
    server.on('listening', () => {
      console.log('Listening on port %s !', port)
    })
  }
)
