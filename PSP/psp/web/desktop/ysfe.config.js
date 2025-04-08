/* Copyright (C) 2016-present, Yuansuan.cn */
module.exports = {
  iconfontUrl: [
    '//at.alicdn.com/t/font_1727885_epy6khpqw6.js',
    '//at.alicdn.com/t/font_3429707_ji71q6it3uf.js',
    '//at.alicdn.com/t/font_1553138_jvtenh7aix.js'
  ],
  webpack(config) {
    if (config.mode === 'development') {
      config.devServer = {
        hot: true,
        overlay: true,
        proxy: {
          '/api/v1': {
            target: 'http://10.0.4.48:32432',
            ws: true,
            secure: false
          }
        },
        host: '0.0.0.0',
        port: process.env.PORT || 8082,
        historyApiFallback: true
      }
    }
    config.module.rules.unshift({
      test: /\.worker\.ts$/,
      loader: 'worker-loader'
    })
    return config
  }
}
