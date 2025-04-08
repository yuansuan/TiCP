const resolveApp = require('./resolveApp')
const appYsfeConfig = resolveApp('ysfe.config.js')
const fs = require('fs-extra')

module.exports = function () {
  if (fs.existsSync(appYsfeConfig)) {
    return require(appYsfeConfig)
  }

  return undefined
}
