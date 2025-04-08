const fs = require('fs-extra')
const defaultConfig = require('../babel.config')
const getYsfeConfig = require('./getYsfeConfig')

module.exports = function () {
  const ysfeConfig = getYsfeConfig()
  if (ysfeConfig && ysfeConfig.babel) {
    return ysfeConfig.babel(defaultConfig)
  }

  return defaultConfig
}
